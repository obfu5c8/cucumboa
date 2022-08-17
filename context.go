package cucumboa

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/obfu5c8/cucumboa/internal/openapi"
)

type Context struct {
	schema  openapi.Schema
	openapi *openapi.Utils

	dispatcher Dispatcher

	operation openapi.Operation

	requestPath       string
	requestPathParams map[string]string
	requestMethod     string
	requestBody       string

	requestSent bool
	response    *http.Response
}

func (c *Context) GetResponse() *http.Response {
	if !c.requestSent {
		c.sendRequest()
	}
	return c.response
}

func (c *Context) GetResponseBody() []byte {
	if !c.requestSent {
		c.sendRequest()
	}
	body, _ := ioutil.ReadAll(c.response.Body)

	// Put the content back in the body to read next time
	c.response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}

// Specifies the operation we want to call
func (c *Context) SetOperation(operationId string) error {
	// Find the operation in the schema
	method, pattern, operation, err := c.openapi.GetOperation(operationId)
	if err != nil {
		return errors.New(fmt.Sprintf("Unknown operation '%s'", operationId))
	}

	// // Fail if operation is expecting path params, since we didn't pass any
	// expectedPathParams := openapi.ExtractOperationPathParamNames(operation)
	// if len(expectedPathParams) > 0 {
	// 	return errors.New(fmt.Sprintf("Operation '%s' expects %d path params, but none were specified.", operationId, len(expectedPathParams)))
	// }

	c.operation = operation
	c.requestPath = pattern
	c.requestMethod = method
	c.requestSent = false

	return nil
}

func (c *Context) SetPathParams(params map[string]string) error {
	c.requestPathParams = params
	return nil
}

func (c *Context) sendRequest() {
	request := c.buildRequest()
	c.response, _ = c.dispatcher.Dispatch(request)
	c.requestSent = true
}

func (c *Context) buildRequest() *http.Request {
	url := c.buildRequestUrl()
	req, _ := http.NewRequest(c.requestMethod, url, strings.NewReader(c.requestBody))

	return req
}

func (c *Context) buildRequestUrl() string {

	var url string = c.requestPath

	pathParams := openapi.ExtractOperationPathParamNames(c.operation)

	for _, param := range pathParams {
		slug := fmt.Sprintf("{%s}", param)
		url = strings.ReplaceAll(url, slug, c.requestPathParams[param])
	}

	return url
}
