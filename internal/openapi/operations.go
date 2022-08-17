package openapi

import (
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

func findOperation(schema Schema, operationId string) (method string, pattern string, operation Operation, err error) {

	for pattern, path := range schema.Paths {

		methods := map[string]*openapi3.Operation{
			"GET":    path.Get,
			"POST":   path.Post,
			"PUT":    path.Put,
			"PATCH":  path.Patch,
			"DELETE": path.Delete,
			"HEAD":   path.Head,
		}

		for method, op := range methods {
			if op != nil && op.OperationID == operationId {
				return method, pattern, op, nil
			}
		}
	}
	return "", "", nil, errors.New(fmt.Sprintf("No operation called '%s' found", operationId))
}
