package runner

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sync"

	"github.com/cucumber/cucumber-engine/src/dto"
	messages "github.com/cucumber/cucumber-messages-go/v2"
	gherkin "github.com/cucumber/gherkin-go"
	uuid "github.com/satori/go.uuid"
)

// Runner executes a run of cucumber
type Runner struct {
	incomingCommands     chan *messages.Wrapper
	outgoingCommands     chan *messages.Wrapper
	responseChannelMutex sync.RWMutex
	responseChannels     map[string]chan *messages.Wrapper
	result               *dto.TestRunResult
}

// NewRunner creates a runner
func NewRunner() *Runner {
	r := &Runner{
		incomingCommands: make(chan *messages.Wrapper),
		outgoingCommands: make(chan *messages.Wrapper),
		responseChannels: map[string]chan *messages.Wrapper{},
		result:           dto.NewTestRunResult(),
	}
	go func() {
		for command := range r.incomingCommands {
			go r.receiveCommand(command)
		}
	}()
	return r
}

// GetCommandChannels returns the command channels
func (r *Runner) GetCommandChannels() (chan *messages.Wrapper, chan *messages.Wrapper) {
	return r.incomingCommands, r.outgoingCommands
}

func (r *Runner) receiveCommand(command *messages.Wrapper) {
	switch x := command.Message.(type) {
	case *messages.Wrapper_CommandStart:
		r.start(x.CommandStart)
	case *messages.Wrapper_CommandActionComplete:
		r.responseChannelMutex.RLock()
		if responseChannel, ok := r.responseChannels[x.CommandActionComplete.GetCompletedId()]; ok {
			responseChannel <- command
		}
		r.responseChannelMutex.RUnlock()
	}
}

func (r *Runner) sendCommand(command *messages.Wrapper) {
	r.outgoingCommands <- command
}

func (r *Runner) sendError(err error) {
	r.sendCommand(&messages.Wrapper{
		Message: &messages.Wrapper_CommandError{
			CommandError: err.Error(),
		},
	})
}

func (r *Runner) start(command *messages.CommandStart) {
	acceptedPickles, err := r.getAcceptedPickles(command.GetBaseDirectory(), command.SourcesConfig)
	if err != nil {
		r.sendError(err)
		return
	}
	supportCodeLibrary, err := NewSupportCodeLibrary(command.SupportCodeConfig)
	if err != nil {
		r.sendError(err)
		return
	}
	r.sendCommand(&messages.Wrapper{
		Message: &messages.Wrapper_TestRunStarted{
			TestRunStarted: &messages.TestRunStarted{},
		},
	})
	if len(acceptedPickles) > 0 {
		_ = r.sendCommandAndAwaitResponse(&messages.Wrapper{
			Message: &messages.Wrapper_CommandRunBeforeTestRunHooks{
				CommandRunBeforeTestRunHooks: &messages.CommandRunBeforeTestRunHooks{},
			},
		})
	}
	var runTestCasesFunc func(*runTestCasesOptions) (bool, error)
	if command.RuntimeConfig.MaxParallel == 0 || command.RuntimeConfig.MaxParallel > 1 {
		runTestCasesFunc = RunTestCasesInParallel
	} else {
		runTestCasesFunc = RunTestCasesSequentially
	}
	testRunResult, err := runTestCasesFunc(&runTestCasesOptions{
		baseDirectory:               command.BaseDirectory,
		pickles:                     acceptedPickles,
		runtimeConfig:               command.RuntimeConfig,
		sendCommand:                 r.sendCommand,
		sendCommandAndAwaitResponse: r.sendCommandAndAwaitResponse,
		supportCodeLibrary:          supportCodeLibrary,
	})
	if err != nil {
		r.sendError(err)
		return
	}
	if len(acceptedPickles) > 0 {
		_ = r.sendCommandAndAwaitResponse(&messages.Wrapper{
			Message: &messages.Wrapper_CommandRunAfterTestRunHooks{
				CommandRunAfterTestRunHooks: &messages.CommandRunAfterTestRunHooks{},
			},
		})
	}
	r.sendCommand(&messages.Wrapper{
		Message: &messages.Wrapper_TestRunFinished{
			TestRunFinished: &messages.TestRunFinished{Success: testRunResult},
		},
	})
	close(r.outgoingCommands)
}

func (r *Runner) getAcceptedPickles(baseDirectory string, sourcesConfig *messages.SourcesConfig) ([]*messages.Pickle, error) {
	pickleFilter, err := NewPickleFilter(sourcesConfig.Filters)
	if err != nil {
		return nil, err
	}
	gherkinMessages, err := gherkin.Messages(sourcesConfig.AbsolutePaths, nil, sourcesConfig.Language, true, true, true, nil, false)
	if err != nil {
		return nil, err
	}
	acceptedPickles := []*messages.Pickle{}
	for i, gherkinMessage := range gherkinMessages {
		r.sendCommand(&gherkinMessages[i])
		switch x := gherkinMessage.Message.(type) {
		case *messages.Wrapper_Attachment:
			uri, err := filepath.Rel(baseDirectory, x.Attachment.Source.Uri)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("Parse error in '%s': %s", uri, x.Attachment.Data)
		case *messages.Wrapper_Pickle:
			pickle := x.Pickle
			if pickleFilter.Matches(pickle) {
				r.sendCommand(&messages.Wrapper{
					Message: &messages.Wrapper_PickleAccepted{
						PickleAccepted: &messages.PickleAccepted{PickleId: pickle.Id},
					},
				})
				acceptedPickles = append(acceptedPickles, pickle)
			} else {
				r.sendCommand(&messages.Wrapper{
					Message: &messages.Wrapper_PickleRejected{
						PickleRejected: &messages.PickleRejected{PickleId: pickle.Id},
					},
				})
			}
		}
	}
	if sourcesConfig.Order.Type == messages.SourcesOrderType_RANDOM {
		reorderPickles(acceptedPickles, sourcesConfig.Order.Seed)
	}
	return acceptedPickles, nil
}

func (r *Runner) sendCommandAndAwaitResponse(command *messages.Wrapper) *messages.Wrapper {
	id := uuid.NewV4().String()
	switch x := command.Message.(type) {
	case *messages.Wrapper_CommandRunBeforeTestRunHooks:
		x.CommandRunBeforeTestRunHooks.ActionId = id
	case *messages.Wrapper_CommandRunAfterTestRunHooks:
		x.CommandRunAfterTestRunHooks.ActionId = id
	case *messages.Wrapper_CommandInitializeTestCase:
		x.CommandInitializeTestCase.ActionId = id
	case *messages.Wrapper_CommandRunBeforeTestCaseHook:
		x.CommandRunBeforeTestCaseHook.ActionId = id
	case *messages.Wrapper_CommandRunAfterTestCaseHook:
		x.CommandRunAfterTestCaseHook.ActionId = id
	case *messages.Wrapper_CommandRunTestStep:
		x.CommandRunTestStep.ActionId = id
	case *messages.Wrapper_CommandGenerateSnippet:
		x.CommandGenerateSnippet.ActionId = id
	}
	responseChannel := make(chan *messages.Wrapper)
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

func reorderPickles(pickles []*messages.Pickle, seed uint64) {
	seededRand := rand.New(rand.NewSource(int64(seed)))
	N := len(pickles)
	for i := 0; i < N; i++ {
		j := i + seededRand.Intn(N-i)
		pickles[j], pickles[i] = pickles[i], pickles[j]
	}
}
