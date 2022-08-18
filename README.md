# Cucumboa
## Helpers for using cucumber to describe OpenAPI schemas in go

---


## Why?

OpenApi schemas are great at describing the _structure_ of your API. You can list all the possible
response status codes for your endpoints and the data models they might return, but that's where
it stops. They don't describe the API's _behaviour_. In what situations would a 404 status be returned? Which model from an `anyOf` schema will I receive if I pass an invalid query param?

Behavioural tests help to solve this problem. By using cucumber to describe the behaviours of the API
service it allows us to document these behaviours in a way that is both easy to understand as a consumer
but also easy to generate tests for to automatically validate the implementation.


## Example

```gherkin
Scenario: Testing the GetPet endpoint

    Feature: Successfully retrieving a Pet
        Given pet '1234' exists
        When the 'Find pet by ID' operation is called with path params:
            | petId | 1234 |
        Then the response status will be '200'

    Feature: Invalid Id supplied
        When the 'Find pet by ID' operation is called with path params:
            | petId |  this_isnt_a_number |
        Then the response status will be '400'

    Feature: Requesting a nonexistant pet
        Given pet '9876' does not exist
        When the 'Find pet by ID' operation is called with path params:
            | petId |  9876 |
        Then the response status will be '404'

```

## Installation
`go get github.com/obfu5c8/cucumboa`

## Getting Started
```go
// spec_test.go


// Load our OpenApi schema
var schema = cucumboa.MustLoadOpenApiSchemaFromFile("../../../api/openapi.yml")


func initializeScenario(ctx *godog.ScenarioContext) {

	// Create the API http.Handler from our server implementation
	httpHandler := myApi.NewHttpHandler()

	// Create a dispatcher to allow cucumboa to send requests to our API
	dispatcher := cucumboa.CreateHandlerDispatcher(httpHandler)

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

```



## Built-in Steps

### Request configuration

#### _`... the {name} operation is called`_

#### _`... the {name} operation is called with path params:`_

### Response Assertions

#### _`... the response status will be '{code}'`_



## Extending with a DSL
While cucumboa provides some great built-in steps for testing any generic OpenApi service, sometimes they can seem a bit verbose and clutter the view of the specification.
To make the spec cleaner to read you can also extend the provided steps to create a DSL for your service.

### Example
```gherkin
Scenario: Comparing DSL & vanilla steps
    
    Feature: Using the vanilla steps
        Given pet '1234' does not exist
        When the 'Find pet by ID' operation is called with path params:
            | petId | 1234 |
        Then the response status will be '404'
        And the content will contain values:
            | error.code | 0x234 |

    Feature: Wrapping in a custom DSL
        Given pet '1234' does not exist
        When the 'Find pet by ID' operation is called with id '1234'
        Then the response status will be '404'
        And the error code will be '0x234'
```

```go

type DSL struct {
    c *cucumboa.Context
}

func (dsl *DSL) CallOperationWithIdPathParamStep(operation string, id string) error {
    dsl.c.SetOperation(operation)
    dsl.c.SetPathParams(map[string]string{
        petId: id
    })
    return nil
}

func (dsl *DSL) AssertResponseContentErrorCodeStep(code string) {
    dsl.c.AssertResponseContentContainsValues(map[string]string{
        code: code,
    })

}
```