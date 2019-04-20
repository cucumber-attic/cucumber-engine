package cli

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/cucumber/cucumber-engine/src/runner"
	messages "github.com/cucumber/cucumber-messages-go/v2"
	protobufio "github.com/gogo/protobuf/io"
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
		writer := protobufio.NewDelimitedWriter(os.Stdout)
		for command := range outgoing {
			if *debugFlag {
				fmt.Fprintf(os.Stderr, "cucumber-engine OUT: %+v\n", command)
			}
			err := writer.WriteMsg(command)
			if err != nil {
				panic(err)
			}
		}
		done <- true
	}()
	reader := protobufio.NewDelimitedReader(os.Stdin, math.MaxInt32)
	for {
		command := &messages.Wrapper{}
		err := reader.ReadMsg(command)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if *debugFlag {
			fmt.Fprintf(os.Stderr, "cucumber-engine IN: %+v\n", command)
		}
		incoming <- command
	}
	<-done
}
