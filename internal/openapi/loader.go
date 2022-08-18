package openapi

import (
	"context"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadOpenApiSchemaFromFile(filePath string) (Schema, error) {
	ctx := context.Background()

	doc, err := openapi3.NewLoader().LoadFromFile(filePath)
	if err != nil {
		return nil, err
	}

	err = doc.Validate(ctx)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func LoadOpenApiSchemaFromUrl(url *url.URL) (Schema, error) {
	ctx := context.Background()

	doc, err := openapi3.NewLoader().LoadFromURI(url)
	if err != nil {
		return nil, err
	}

	err = doc.Validate(ctx)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
