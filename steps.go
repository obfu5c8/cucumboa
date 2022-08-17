package cucumboa

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/tidwall/gjson"
)

func bindSteps(sc *godog.ScenarioContext, ctx Context) {
	steps := steps{ctx: ctx}

	// Request configuration
	sc.Step(`^the \'([^']+)\' operation is called$`, steps.callOperationStep)
	sc.Step(`^the \'([^']+)\' operation is called with path params:$`, steps.callOperationWithPathParamsStep)

	// // Response assertions
	sc.Step(`^the response status will be \'(\d+)\'$`, steps.assertResponseStatusStep)
	sc.Step(`^the content will have values:$`, steps.assertResponseContentContainsTableValues)
}

type steps struct {
	ctx Context
}

// Sets the operation to call in the request
func (s *steps) callOperationStep(operationId string) error {
	return s.ctx.SetOperation(operationId)
}

// Sets the operation to call along with any path params
func (s *steps) callOperationWithPathParamsStep(operationId string, pathParamsTable *godog.Table) error {
	err := s.ctx.SetOperation(operationId)
	if err != nil {
		return err
	}

	// Turn gherkin table into a map
	pathParams := make(map[string]string)
	for _, row := range pathParamsTable.Rows {
		pathParams[row.Cells[0].Value] = row.Cells[1].Value
	}

	return s.ctx.SetPathParams(pathParams)
}

// Asserts the response status is expected
func (s *steps) assertResponseStatusStep(status string) error {
	expected, _ := strconv.ParseInt(status, 10, 16)

	// Assert the expected response code was received
	res := s.ctx.GetResponse()
	got := res.StatusCode
	if int(expected) != got {
		return errors.New(fmt.Sprintf("Expected response status was '%d'. Got '%d'", expected, got))
	}

	// Validate the response body against the schema
	err := ValidateResponseBody(&s.ctx)
	if err != nil {
		ourMsg := "The expected response code was received, however"
		return errors.New(fmt.Sprintf("%s %s", ourMsg, err.Error()))
	}

	return nil
}

// Asserts that the response body contains certain values
// This is not an exclusive test, it will only check for properties
// passed to the step handler, and additional properties in the response
// body are ignored
func (s *steps) assertResponseContentContainsTableValues(table *godog.Table) error {
	// Convert table to map of values
	values := make(map[string]string)
	for _, row := range table.Rows {
		values[row.Cells[0].Value] = row.Cells[1].Value
	}

	return s.assertResponseContentContainsValues(values)
}

func (s *steps) assertResponseContentContainsValues(values map[string]string) error {
	body := s.ctx.GetResponseBody()

	for keyPath, expected := range values {
		result := gjson.GetBytes(body, keyPath)
		got := result.String()

		if got != expected {
			return errors.New(fmt.Sprintf("Expected content property '%s' to be '%s'. Got '%s'.\n%s", keyPath, expected, got, body))
		}
	}
	return nil
}
