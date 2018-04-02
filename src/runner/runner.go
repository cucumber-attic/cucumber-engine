package runner

import (
	"fmt"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
	gherkin "github.com/cucumber/gherkin-go"
	uuid "github.com/satori/go.uuid"
)

// Runner executes a run of cucumber
type Runner struct {
	incomingCommands chan *dto.Command
	outgoingCommands chan *dto.Command
	runtimeConfig    *dto.RuntimeConfig
	pickleEvents     []*gherkin.PickleEvent
	testCaseRunner   *TestCaseRunner
	responseChannels map[string]chan *dto.Command
}

// NewRunner creates a runner
func NewRunner() *Runner {
	r := &Runner{
		incomingCommands: make(chan *dto.Command),
		outgoingCommands: make(chan *dto.Command),
		responseChannels: map[string]chan *dto.Command{},
	}
	go func() {
		for command := range r.incomingCommands {
			fmt.Printf("source %+v\n", command)
			go r.receiveCommand(command)
		}
	}()
	return r
}

// GetCommandChannels returns the command channels
func (r *Runner) GetCommandChannels() (chan *dto.Command, chan *dto.Command) {
	return r.incomingCommands, r.outgoingCommands
}

func (r *Runner) receiveCommand(command *dto.Command) {
	if command.Type == "start" {
		r.runtimeConfig = command.RuntimeConfig
		r.start(command.FeaturesConfig)
		return
	}
	if responseChannel, ok := r.responseChannels[command.ResponseTo]; ok {
		responseChannel <- command
	}
}

func (r *Runner) sendCommand(command *dto.Command) {
	r.outgoingCommands <- command
}

func (r *Runner) start(featuresConfig *dto.FeaturesConfig) {
	events, err := gherkin.GherkinEvents(featuresConfig.AbsolutePaths...)
	if err != nil {
		panic(err)
	}
	for _, event := range events {
		r.sendCommand(&dto.Command{
			Type:  "event",
			Event: event,
		})
		if pickleEvent, ok := event.(*gherkin.PickleEvent); ok {
			// TODO filter
			r.pickleEvents = append(r.pickleEvents, pickleEvent)
		}
	}
	r.sendCommand(&dto.Command{
		Type:  "event",
		Event: event.NewTestRunStarted(),
	})
	success := true
	_ = r.sendCommandAndAwaitResponse(&dto.Command{Type: dto.CommandTypeRunBeforeTestRunHooks})
	for _, pickleEvent := range r.pickleEvents {
		r.testCaseRunner = NewTestCaseRunner(&NewTestCaseRunnerOptions{
			AfterTestCaseHookDefinitions:  r.runtimeConfig.AfterTestCaseHookDefinitions,
			BeforeTestCaseHookDefinitions: r.runtimeConfig.BeforeTestCaseHookDefinitions,
			Pickle:                      pickleEvent.Pickle,
			SendCommand:                 r.sendCommand,
			SendCommandAndAwaitResponse: r.sendCommandAndAwaitResponse,
			URI: pickleEvent.URI,
		})
		result := r.testCaseRunner.Run()
		if r.shouldCauseFailure(result.Status) {
			success = false
		}
	}
	_ = r.sendCommandAndAwaitResponse(&dto.Command{Type: dto.CommandTypeRunAfterTestRunHooks})
	r.sendCommand(&dto.Command{
		Type:  "event",
		Event: event.NewTestRunFinished(success),
	})
	close(r.outgoingCommands)
}

func (r *Runner) sendCommandAndAwaitResponse(command *dto.Command) *dto.Command {
	id := uuid.NewV4().String()
	command.ID = id
	responseChannel := make(chan *dto.Command)
	r.responseChannels[id] = responseChannel
	go r.sendCommand(command)
	return <-responseChannel
}

func (r *Runner) shouldCauseFailure(status dto.Status) bool {
	return status == dto.StatusAmbiguous ||
		status == dto.StatusFailed ||
		status == dto.StatusUndefined ||
		(status == dto.StatusPending && r.runtimeConfig.IsStrict)
}
