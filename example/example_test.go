package main

import (
	"testing"

	"github.com/cucumber/godog"
	"github.com/obfu5c8/go-api/internal/api/rest"
	"github.com/obfu5c8/go-api/pkg/cucumboa"
)

var schema = cucumboa.MustLoadOpenApiSchemaFromFile("../../../api/openapi.yml")

func InitializeScenario(ctx *godog.ScenarioContext) {
	// ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	// 	return ctx, nil
	// })

	// Create the API handler from our server implementation
	apiHandler := rest.DefaultRestApiHandler()

	// Create a dispatcher to allow cucumboa to send requests to our API
	dispatcher := cucumboa.CreateHandlerDispatcher(apiHandler)

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
