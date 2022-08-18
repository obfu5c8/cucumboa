package main

import (
	"testing"

	"github.com/cucumber/godog"
	"github.com/obfu5c8/cucumboa"
	"github.com/obfu5c8/cucumboa/examples/mockserver"
)

// Load the OpenApi schema from a remote url
var schema = cucumboa.MustLoadOpenApiSchemaFromUrl("https://petstore3.swagger.io/api/v3/openapi.json")

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Create the API handler from our server implementation
	mockserver := mockserver.New()

	// Create a dispatcher to allow cucumboa to send requests to our API
	dispatcher := cucumboa.CreateHandlerDispatcher(mockserver.Handler())

	// Initialise cucumboa against the scenario
	cucumboa.InitializeScenario(ctx, cucumboa.Options{
		Schema:     schema,
		Dispatcher: dispatcher,
	})
}

func TestApiSpec(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"./example.feature"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}

}
