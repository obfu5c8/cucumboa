package cucumboa

import (
	"errors"
	"fmt"

	"github.com/obfu5c8/cucumboa/internal/openapi"
	"github.com/tidwall/gjson"
)

func ValidateResponseBody(ctx *Context) error {

	response := ctx.GetResponse()
	request := response.Request

	err := ctx.openapi.ValidateResponseAgainstRoute(request, response, openapi.Route{
		Path:       ctx.requestPath,
		PathParams: ctx.requestPathParams,
		Method:     ctx.requestMethod,
		Operation:  ctx.operation,
	})
	return err
}

func AssertResponseContentContainsValues(ctx *Context, values map[string]string) error {
	body := ctx.GetResponseBody()

	for keyPath, expected := range values {
		result := gjson.GetBytes(body, keyPath)
		got := result.String()

		if got != expected {
			return errors.New(fmt.Sprintf("Expected content property '%s' to be '%s'. Got '%s'.\n%s", keyPath, expected, got, body))
		}
	}
	return nil
}
