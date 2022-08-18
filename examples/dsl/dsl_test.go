package main

import (
	"testing"

	"github.com/cucumber/godog"
	"github.com/obfu5c8/cucumboa"
	"github.com/obfu5c8/cucumboa/examples/mockserver"
)

// Load the OpenApi schema
var schema = cucumboa.MustLoadOpenApiSchemaFromUrl("https://petstore3.swagger.io/api/v3/openapi.json")

// Define our custom DSL to improve readability
func AddDSL(sc *godog.ScenarioContext, c *cucumboa.Context) {

	sc.Step("^the '([^']+)' operation is called for pet '(\\d+)'$", func(opId string, petId string) {
		c.SetOperation(opId)
		c.SetPathParams(map[string]string{"petId": petId})
	})

	sc.Step("^the pet will be called '([^']+)'$", func(expectedName string) {
		cucumboa.AssertResponseContentContainsValues(c, map[string]string{
			"name": expectedName,
		})
	})

}

// Create the scenario initializer
// This is where we configure our server and mocks
func InitializeScenario(ctx *godog.ScenarioContext) {
	// Create the API handler from our server implementation
	mockserver := mockserver.New()

	// Create a dispatcher to allow cucumboa to send requests to our API
	dispatcher := cucumboa.CreateHandlerDispatcher(mockserver.Handler())

	// Initialise cucumboa against the scenario
	cCtx, _ := cucumboa.InitializeScenario(ctx, cucumboa.Options{
		Schema:     schema,
		Dispatcher: dispatcher,
	})

	// Initialise our custom DSL steps
	AddDSL(ctx, cCtx)
}

// Wrap in a `go test` compatible function and execute
func TestApiSpec(t *testing.T) {
	cucumboa.RunSimpleTestSuite(t, InitializeScenario, []string{
		"./example.feature",
	})
}
