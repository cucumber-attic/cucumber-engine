package runner

import (
	"math/rand"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
	gherkin "github.com/cucumber/gherkin-go"
	uuid "github.com/satori/go.uuid"
)

// Runner executes a run of cucumber
type Runner struct {
	incomingCommands chan *dto.Command
	outgoingCommands chan *dto.Command
	isStrict         bool
	responseChannels map[string]chan *dto.Command
	result           *dto.TestRunResult
}

// NewRunner creates a runner
func NewRunner() *Runner {
	r := &Runner{
		incomingCommands: make(chan *dto.Command),
		outgoingCommands: make(chan *dto.Command),
		responseChannels: map[string]chan *dto.Command{},
		result:           &dto.TestRunResult{Success: true, Duration: 0},
	}
	go func() {
		for command := range r.incomingCommands {
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
		r.start(command)
		return
	}
	if responseChannel, ok := r.responseChannels[command.ResponseTo]; ok {
		responseChannel <- command
	}
}

func (r *Runner) sendCommand(command *dto.Command) {
	r.outgoingCommands <- command
}

func (r *Runner) start(command *dto.Command) {
	acceptedPickleEvents, err := r.getAcceptedPickleEvents(command.FeaturesConfig)
	if err != nil {
		r.sendCommand(&dto.Command{
			Type:  dto.CommandTypeError,
			Error: err.Error(),
		})
		return
	}
	supportCodeLibrary, err := NewSupportCodeLibrary(command.SupportCodeConfig)
	if err != nil {
		r.sendCommand(&dto.Command{
			Type:  dto.CommandTypeError,
			Error: err.Error(),
		})
		return
	}
	r.sendCommand(&dto.Command{
		Type:  "event",
		Event: event.NewTestRunStarted(),
	})
	if len(acceptedPickleEvents) > 0 {
		_ = r.sendCommandAndAwaitResponse(&dto.Command{Type: dto.CommandTypeRunBeforeTestRunHooks})
	}
	testRunResult := &dto.TestRunResult{Success: true, Duration: 0}
	isSkipped := command.RuntimeConfig.IsDryRun
	for _, pickleEvent := range acceptedPickleEvents {
		testCaseRunner, err := NewTestCaseRunner(&NewTestCaseRunnerOptions{
			BaseDirectory:               command.BaseDirectory,
			ID:                          uuid.NewV4().String(),
			IsSkipped:                   isSkipped,
			Pickle:                      pickleEvent.Pickle,
			SendCommand:                 r.sendCommand,
			SendCommandAndAwaitResponse: r.sendCommandAndAwaitResponse,
			SupportCodeLibrary:          supportCodeLibrary,
			URI:                         pickleEvent.URI,
		})
		if err != nil {
			r.sendCommand(&dto.Command{
				Type:  dto.CommandTypeError,
				Error: err.Error(),
			})
			return
		}
		testCaseResult := testCaseRunner.Run()
		testRunResult.Duration += testCaseResult.Duration
		if r.shouldCauseFailure(testCaseResult.Status, command.RuntimeConfig.IsStrict) {
			testRunResult.Success = false
			if !isSkipped && command.RuntimeConfig.IsFailFast {
				isSkipped = true
			}
		}
	}
	if len(acceptedPickleEvents) > 0 {
		_ = r.sendCommandAndAwaitResponse(&dto.Command{Type: dto.CommandTypeRunAfterTestRunHooks})
	}
	r.sendCommand(&dto.Command{
		Type:  "event",
		Event: event.NewTestRunFinished(testRunResult),
	})
	close(r.outgoingCommands)
}

func (r *Runner) getAcceptedPickleEvents(featuresConfig *dto.FeaturesConfig) ([]*gherkin.PickleEvent, error) {
	pickleFilter, err := NewPickleFilter(featuresConfig.Filters)
	if err != nil {
		return nil, err
	}
	gherkinEvents, err := gherkin.GherkinEvents(featuresConfig.AbsolutePaths...)
	if err != nil {
		return nil, err
	}
	acceptedPickleEvents := []*gherkin.PickleEvent{}
	for _, gherkinEvent := range gherkinEvents {
		r.sendCommand(&dto.Command{
			Type:  "event",
			Event: gherkinEvent,
		})
		if pickleEvent, ok := gherkinEvent.(*gherkin.PickleEvent); ok {
			if pickleFilter.Matches(pickleEvent) {
				r.sendCommand(&dto.Command{
					Type:  "event",
					Event: event.NewPickleAccepted(pickleEvent),
				})
				acceptedPickleEvents = append(acceptedPickleEvents, pickleEvent)
			} else {
				r.sendCommand(&dto.Command{
					Type:  "event",
					Event: event.NewPickleRejected(pickleEvent),
				})
			}
		}
	}
	if featuresConfig.Order.Type == dto.FeaturesOrderTypeRandom {
		reorderPickleEvents(acceptedPickleEvents, featuresConfig.Order.Seed)
	}
	return acceptedPickleEvents, nil
}

func (r *Runner) sendCommandAndAwaitResponse(command *dto.Command) *dto.Command {
	id := uuid.NewV4().String()
	command.ID = id
	responseChannel := make(chan *dto.Command)
	r.responseChannels[id] = responseChannel
	go r.sendCommand(command)
	return <-responseChannel
}

func (r *Runner) shouldCauseFailure(status dto.Status, isStrict bool) bool {
	return status == dto.StatusAmbiguous ||
		status == dto.StatusFailed ||
		status == dto.StatusUndefined ||
		(status == dto.StatusPending && isStrict)
}

func reorderPickleEvents(pickleEvents []*gherkin.PickleEvent, seed int64) {
	seededRand := rand.New(rand.NewSource(seed))
	N := len(pickleEvents)
	for i := 0; i < N; i++ {
		j := i + seededRand.Intn(N-i)
		pickleEvents[j], pickleEvents[i] = pickleEvents[i], pickleEvents[j]
	}
}
