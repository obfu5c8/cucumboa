package cucumboa

import (
	"net/url"
	"testing"

	"github.com/cucumber/godog"
	"github.com/obfu5c8/cucumboa/internal/openapi"
)

type Options struct {
	// Path or URL to OpenApi schema for this API
	Schema openapi.Schema
	// Http dispatcher to handle sending of HTTP requests to your API
	Dispatcher Dispatcher
}

// Initialize a godog scenario with a default context.
// Returns the context object so you can use it for extended DSLs
func InitializeScenario(scenario *godog.ScenarioContext, opts Options) (*Context, error) {
	ctx := newContext(opts)
	return InitializeScenarioWithContext(scenario, &ctx)
}

// Initialzie a godog scenario by passing in your own context
// Returns the context object so you can use it for extended DSLs
func InitializeScenarioWithContext(scenario *godog.ScenarioContext, ctx *Context) (*Context, error) {
	bindSteps(scenario, ctx)
	return ctx, nil
}

// Load an OpenApi schema from a file
func LoadOpenApiSchemaFromFile(filePath string) (openapi.Schema, error) {
	return openapi.LoadOpenApiSchemaFromFile(filePath)
}

// Load an OpenApi schema from a file, or panic if something goes wrong
func MustLoadOpenApiSchemaFromFile(filePath string) openapi.Schema {
	schema, err := LoadOpenApiSchemaFromFile(filePath)
	if err != nil {
		panic(err)
	}
	return schema
}

// Load an OpenApi schema from a URL
func LoadOpenApiSchemaFromUrl(uri string) (openapi.Schema, error) {
	url, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	return openapi.LoadOpenApiSchemaFromUrl(url)
}

// Load an OpenApi schema from a URL, or panic if something goes wrong
func MustLoadOpenApiSchemaFromUrl(uri string) openapi.Schema {
	schema, err := LoadOpenApiSchemaFromUrl(uri)
	if err != nil {
		panic(err)
	}
	return schema
}

func RunSimpleTestSuite(t *testing.T, initScenario func(ctx *godog.ScenarioContext), features []string) {
	suite := godog.TestSuite{
		ScenarioInitializer: initScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    features,
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
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
