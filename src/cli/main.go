package main

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/runner"
)

func main() {
	r := runner.NewRunner(func(command *dto.Command) {
		data, err := json.Marshal(command)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(data)
	})
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := dto.Command{}
		data := scanner.Bytes()
		if err := json.Unmarshal(data, command); err != nil {
			panic(err)
		}
		r.ReceiveCommand(&command)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
