package cucumboa

import (
	"net/http"
	"net/http/httptest"
)

type Dispatcher interface {
	Dispatch(request *http.Request) (*http.Response, error)
}

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
