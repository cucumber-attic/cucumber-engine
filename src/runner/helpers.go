package runner

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/cucumber/cucumber-engine/src/dto"
	"github.com/olekukonko/tablewriter"
)

func getAmbiguousStepDefinitionsMessage(stepDefinitions []*dto.StepDefinition, baseDirectory string) (string, error) {
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
			uri := stepDefinition.URI
			if baseDirectory != "" {
				var err error
				uri, err = filepath.Rel(baseDirectory, uri)
				if err != nil {
					return "", err
				}
			}
			location = fmt.Sprintf("%s:%d", uri, stepDefinition.Line)
		}
		data = append(data, []string{"'" + stepDefinition.Expression.Source() + "'", location})
	}
	table.AppendBulk(data)
	table.Render()
	return fmt.Sprintf("Multiple step definitions match:\n%v", buf.String()), nil
}
