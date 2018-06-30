package runner

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sync"

	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/cucumber/cucumber-engine/src/dto/event"
	gherkin "github.com/cucumber/gherkin-go"
	uuid "github.com/satori/go.uuid"
)

// Runner executes a run of cucumber
type Runner struct {
	incomingCommands     chan *dto.Command
	outgoingCommands     chan *dto.Command
	isStrict             bool
	responseChannelMutex sync.RWMutex
	responseChannels     map[string]chan *dto.Command
	result               *dto.TestRunResult
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
	r.responseChannelMutex.RLock()
	if responseChannel, ok := r.responseChannels[command.ResponseTo]; ok {
		responseChannel <- command
	}
	r.responseChannelMutex.RUnlock()
}

func (r *Runner) sendCommand(command *dto.Command) {
	r.outgoingCommands <- command
}

func (r *Runner) start(command *dto.Command) {
	acceptedPickleEvents, err := r.getAcceptedPickleEvents(command.BaseDirectory, command.FeaturesConfig)
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
	var runTestCasesFunc func(*runTestCasesOptions) (*dto.TestRunResult, error)
	if command.RuntimeConfig.MaxParallel == -1 || command.RuntimeConfig.MaxParallel > 1 {
		runTestCasesFunc = RunTestCasesInParallel
	} else {
		runTestCasesFunc = RunTestCasesSequentially
	}
	testRunResult, err := runTestCasesFunc(&runTestCasesOptions{
		baseDirectory:               command.BaseDirectory,
		pickleEvents:                acceptedPickleEvents,
		runtimeConfig:               command.RuntimeConfig,
		sendCommand:                 r.sendCommand,
		sendCommandAndAwaitResponse: r.sendCommandAndAwaitResponse,
		supportCodeLibrary:          supportCodeLibrary,
	})
	if err != nil {
		r.sendCommand(&dto.Command{
			Type:  dto.CommandTypeError,
			Error: err.Error(),
		})
		return
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

func (r *Runner) getAcceptedPickleEvents(baseDirectory string, featuresConfig *dto.FeaturesConfig) ([]*gherkin.PickleEvent, error) {
	pickleFilter, err := NewPickleFilter(featuresConfig.Filters)
	if err != nil {
		return nil, err
	}
	gherkinEvents, err := gherkin.GherkinEventsForLanguage(featuresConfig.AbsolutePaths, featuresConfig.Language)
	if err != nil {
		return nil, err
	}
	acceptedPickleEvents := []*gherkin.PickleEvent{}
	for _, gherkinEvent := range gherkinEvents {
		r.sendCommand(&dto.Command{
			Type:  "event",
			Event: gherkinEvent,
		})
		if attachmentEvent, ok := gherkinEvent.(*gherkin.AttachmentEvent); ok && attachmentEvent.Media.Type == "text/x.cucumber.stacktrace+plain" {
			uri, err := filepath.Rel(baseDirectory, attachmentEvent.Source.URI)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("Parse error in '%s': %s", uri, attachmentEvent.Data)
		}
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
	r.responseChannelMutex.Lock()
	r.responseChannels[id] = responseChannel
	r.responseChannelMutex.Unlock()
	go r.sendCommand(command)
	result := <-responseChannel
	r.responseChannelMutex.Lock()
	delete(r.responseChannels, id)
	r.responseChannelMutex.Unlock()
	return result
}

func reorderPickleEvents(pickleEvents []*gherkin.PickleEvent, seed int64) {
	seededRand := rand.New(rand.NewSource(seed))
	N := len(pickleEvents)
	for i := 0; i < N; i++ {
		j := i + seededRand.Intn(N-i)
		pickleEvents[j], pickleEvents[i] = pickleEvents[i], pickleEvents[j]
	}
}
