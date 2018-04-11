package runner

import (
	"bytes"
	"fmt"

	"github.com/cucumber/cucumber-pickle-runner/src/dto"
	"github.com/olekukonko/tablewriter"
)

func getAmbiguousStepDefinitionsMessage(stepDefinitions []*dto.StepDefinition) string {
	buf := bytes.Buffer{}
	table := tablewriter.NewWriter(&buf)
	table.SetBorder(false)
	table.SetRowSeparator(" ")
	table.SetColumnSeparator("-")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	data := [][]string{}
	for _, stepDefinition := range stepDefinitions {
		location := ""
		if (stepDefinition.Line != 0) && stepDefinition.URI != "" {
			location = fmt.Sprintf("%s:%d", stepDefinition.URI, stepDefinition.Line)
		}
		data = append(data, []string{"'" + stepDefinition.Expression.Source() + "'", location})
	}
	table.AppendBulk(data)
	table.Render()
	return fmt.Sprintf("Multiple step definitions match:\n%v", buf.String())
}
