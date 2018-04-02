package main

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/runner"
)

func main() {
	r := runner.NewRunner()
	incoming, outgoing := r.GetCommandChannels()
	done := make(chan bool)
	go func() {
		for command := range outgoing {
			data, err := json.Marshal(command)
			if err != nil {
				panic(err)
			}
			os.Stdout.Write(data)
		}
		done <- true
	}()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := dto.Command{}
		data := scanner.Bytes()
		if err := json.Unmarshal(data, command); err != nil {
			panic(err)
		}
		incoming <- &command
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	<-done
}
