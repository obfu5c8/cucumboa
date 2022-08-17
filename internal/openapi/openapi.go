package openapi

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
)

type Schema *openapi3.T
type Operation *openapi3.Operation

func NewUtils(schema Schema) *Utils {
	return &Utils{
		schema: schema,
	}
}

type Utils struct {
	schema Schema
}

func (u *Utils) GetOperation(id string) (method string, pattern string, operation Operation, err error) {
	method, pattern, operation, err = findOperation(u.schema, id)
	return
}

type Route struct {
	Path       string
	PathParams map[string]string
	Method     string
	Operation  Operation
}

func (u *Utils) ValidateResponseAgainstRoute(request *http.Request, response *http.Response, routeInfo Route) error {

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    request,
		PathParams: routeInfo.PathParams,
		Route: &routers.Route{
			Spec:      u.schema,
			Path:      routeInfo.Path,
			PathItem:  u.schema.Paths[routeInfo.Path],
			Method:    routeInfo.Method,
			Operation: routeInfo.Operation,
		},
	}

	// Validate response
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 response.StatusCode,
		Header:                 response.Header,
	}

	body, _ := io.ReadAll(response.Body)
	response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	responseValidationInput.SetBodyBytes(body)

	err := openapi3filter.ValidateResponse(response.Request.Context(), responseValidationInput)

	return err
}

func ExtractOperationPathParamNames(op Operation) []string {
	var paramNames []string

	for _, param := range op.Parameters {
		if param.Value.In == "path" {
			paramNames = append(paramNames, param.Value.Name)
		}
	}

	return paramNames
}
