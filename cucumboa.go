package cucumboa

import (
	"github.com/cucumber/godog"
	"github.com/obfu5c8/go-api/pkg/cucumboa/internal/openapi"
)

type Options struct {
	// Path or URL to OpenApi schema for this API
	Schema openapi.Schema
	// Http dispatcher to handle sending of HTTP requests to your API
	Dispatcher Dispatcher
}

// Initialize a godog scenario with a default context.
// Returns the context object so you can use it for extended DSLs
func InitializeScenario(scenario *godog.ScenarioContext, opts Options) (Context, error) {
	ctx := newContext(opts)
	return initializeScenarioWithContext(scenario, ctx)
}

// Load an OpenApi schema from a file, or panic if something goes wrong
func MustLoadOpenApiSchemaFromFile(filePath string) openapi.Schema {
	schema, err := openapi.LoadOpenApiSchemaFromFile(filePath)
	if err != nil {
		panic(err)
	}
	return schema
}

// Create a new cucumboa context
func newContext(opts Options) Context {
	ctx := Context{
		schema:     opts.Schema,
		openapi:    openapi.NewUtils(opts.Schema),
		dispatcher: opts.Dispatcher,
	}
	return ctx
}

// Initialzie a godog scenario by passing in your own context
func initializeScenarioWithContext(scenario *godog.ScenarioContext, ctx Context) (Context, error) {
	bindSteps(scenario, ctx)
	return ctx, nil
}
