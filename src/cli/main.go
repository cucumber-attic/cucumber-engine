package cli

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/runner"
)

func Execute() {
	// id := uuid.NewV4().String()
	r := runner.NewRunner()
	incoming, outgoing := r.GetCommandChannels()
	done := make(chan bool)
	go func() {
		for command := range outgoing {
			data, err := json.Marshal(command)
			if err != nil {
				panic(err)
			}
			// fmt.Fprintln(os.Stderr, id+" shared Out: "+string(data))
			os.Stdout.Write(append(data, []byte("\n")...))
		}
		done <- true
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := &dto.Command{}
		data := scanner.Bytes()
		// fmt.Fprintln(os.Stderr, id+" shared In: "+string(data))
		if err := json.Unmarshal(data, command); err != nil {
			panic(err)
		}
		incoming <- command
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	<-done
}
