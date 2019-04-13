package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/cucumber/cucumber-engine/src/runner"
	messages "github.com/cucumber/cucumber-messages-go/v2"
	"github.com/gogo/protobuf/proto"
)

var version string

// Execute implements the command line interface
func Execute() {
	versionFlag := flag.Bool("version", false, "print version")
	debugFlag := flag.Bool("debug", false, "print debug information")
	flag.Parse()
	if *versionFlag {
		fmt.Printf("cucumber-engine %s\n", version)
		os.Exit(0)
	}
	r := runner.NewRunner()
	incoming, outgoing := r.GetCommandChannels()
	done := make(chan bool)
	go func() {
		for command := range outgoing {
			data, err := proto.Marshal(command)
			if err != nil {
				panic(err)
			}
			if *debugFlag {
				fmt.Fprintf(os.Stderr, "cucumber-engine OUT: %s\n", string(data))
			}
			os.Stdout.Write(append(data, []byte("\n")...))
		}
		done <- true
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := &messages.Wrapper{}
		data := scanner.Bytes()
		if *debugFlag {
			fmt.Fprintf(os.Stderr, "CPR IN: %s\n", string(data))
		}
		if err := proto.Unmarshal(data, command); err != nil {
			panic(err)
		}
		incoming <- command
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	<-done
}
