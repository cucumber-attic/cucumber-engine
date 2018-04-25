package formatter_test

import (
	"io"
	"io/ioutil"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/cucumber/cucumber-pickle-runner/src/dto/event"
	"github.com/cucumber/cucumber-pickle-runner/src/formatter"
	gherkin "github.com/cucumber/gherkin-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getRerunFormatForTestCase(status dto.Status) string {
	reader, writer := io.Pipe()
	eventChannel := make(chan gherkin.CucumberEvent)
	formatter.FormatAsRerun(&formatter.CommonOptions{
		BaseDirectory: "/path/to/project/",
		EventChannel:  eventChannel,
		Stream:        writer,
	}, "\n")
	eventChannel <- event.NewTestCaseFinished(event.NewTestCaseFinishedOptions{
		Pickle: &gherkin.Pickle{
			Locations: []gherkin.Location{{Line: 1}},
		},
		Result: &dto.TestResult{Status: status},
		URI:    "/path/to/project/path/to/featureA",
	})
	eventChannel <- event.NewTestRunFinished(&dto.TestRunResult{})
	close(eventChannel)
	output, err := ioutil.ReadAll(reader)
	Expect(err).NotTo(HaveOccurred())
	return string(output)
}

var _ = Describe("FormatAsRerun", func() {
	Describe("with no test cases", func() {
		It("outputs nothing", func() {
			reader, writer := io.Pipe()
			eventChannel := make(chan gherkin.CucumberEvent)
			formatter.FormatAsRerun(&formatter.CommonOptions{
				EventChannel: eventChannel,
				Stream:       writer,
			}, "\n")
			eventChannel <- event.NewTestRunFinished(&dto.TestRunResult{})
			close(eventChannel)
			output, err := ioutil.ReadAll(reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(output)).To(Equal(""))
		})
	})

	Describe("with one ambiguous test case", func() {
		It("outputs the reference needed to rerun the test case", func() {
			output := getRerunFormatForTestCase(dto.StatusAmbiguous)
			Expect(output).To(Equal("path/to/featureA:1"))
		})
	})

	Describe("with one failing test case", func() {
		It("outputs the reference needed to rerun the test case", func() {
			output := getRerunFormatForTestCase(dto.StatusFailed)
			Expect(output).To(Equal("path/to/featureA:1"))
		})
	})

	Describe("with one passing test case", func() {
		It("outputs nothing", func() {
			output := getRerunFormatForTestCase(dto.StatusPassed)
			Expect(output).To(Equal(""))
		})
	})

	Describe("with one pending test case", func() {
		It("outputs the reference needed to rerun the test case", func() {
			output := getRerunFormatForTestCase(dto.StatusPending)
			Expect(output).To(Equal("path/to/featureA:1"))
		})
	})

	Describe("with one undefined test case", func() {
		It("outputs the reference needed to rerun the test case", func() {
			output := getRerunFormatForTestCase(dto.StatusSkipped)
			Expect(output).To(Equal("path/to/featureA:1"))
		})
	})

	Describe("with one skipped test case", func() {
		It("outputs the reference needed to rerun the test case", func() {
			output := getRerunFormatForTestCase(dto.StatusUndefined)
			Expect(output).To(Equal("path/to/featureA:1"))
		})
	})

	Describe("with two failing test cases in the same file", func() {
		It("outputs the references needed to rerun the test cases as a single element", func() {
			reader, writer := io.Pipe()
			eventChannel := make(chan gherkin.CucumberEvent)
			formatter.FormatAsRerun(&formatter.CommonOptions{
				BaseDirectory: "/path/to/project/",
				EventChannel:  eventChannel,
				Stream:        writer,
			}, "\n")
			eventChannel <- event.NewTestCaseFinished(event.NewTestCaseFinishedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
				},
				Result: &dto.TestResult{Status: dto.StatusFailed},
				URI:    "/path/to/project/path/to/featureA",
			})
			eventChannel <- event.NewTestCaseFinished(event.NewTestCaseFinishedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 2}},
				},
				Result: &dto.TestResult{Status: dto.StatusFailed},
				URI:    "/path/to/project/path/to/featureA",
			})
			eventChannel <- event.NewTestRunFinished(&dto.TestRunResult{})
			close(eventChannel)
			output, err := ioutil.ReadAll(reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(output)).To(Equal("path/to/featureA:1:2"))
		})
	})

	Describe("with two failing test cases in different files", func() {
		It("outputs the references needed to rerun the test cases as multiple elements", func() {
			reader, writer := io.Pipe()
			eventChannel := make(chan gherkin.CucumberEvent)
			formatter.FormatAsRerun(&formatter.CommonOptions{
				BaseDirectory: "/path/to/project/",
				EventChannel:  eventChannel,
				Stream:        writer,
			}, "\n")
			eventChannel <- event.NewTestCaseFinished(event.NewTestCaseFinishedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 1}},
				},
				Result: &dto.TestResult{Status: dto.StatusFailed},
				URI:    "/path/to/project/path/to/featureA",
			})
			eventChannel <- event.NewTestCaseFinished(event.NewTestCaseFinishedOptions{
				Pickle: &gherkin.Pickle{
					Locations: []gherkin.Location{{Line: 2}},
				},
				Result: &dto.TestResult{Status: dto.StatusFailed},
				URI:    "/path/to/project/path/to/featureB",
			})
			eventChannel <- event.NewTestRunFinished(&dto.TestRunResult{})
			close(eventChannel)
			output, err := ioutil.ReadAll(reader)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(output)).To(Equal("path/to/featureA:1" + "\n" + "path/to/featureB:2"))
		})
	})
})
