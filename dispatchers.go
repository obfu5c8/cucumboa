package cucumboa

import (
	"net/http"
	"net/http/httptest"
)

// The Dispatcher interface provides a way for cucumboa to call
// OpenApi operations on a target system
type Dispatcher interface {
	Dispatch(request *http.Request) (*http.Response, error)
}

// Craetes a cucumboa dispatcher that runs against a local http.Handler
// instance to allow testing in-memory without needing an http server
func CreateHandlerDispatcher(handler http.Handler) Dispatcher {
	return httpHandlerDispatcher{handler: handler}
}

type httpHandlerDispatcher struct {
	handler http.Handler
}

func (d httpHandlerDispatcher) Dispatch(request *http.Request) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	d.handler.ServeHTTP(recorder, request)

	res := recorder.Result()
	res.Request = request

	return res, nil
}
