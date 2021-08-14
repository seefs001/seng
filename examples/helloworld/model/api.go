/*
 * Swagger Petstore
 *
 * A sample API that uses a petstore as an example to demonstrate features in the OpenAPI 3.0 specification
 *
 * API version: 1.0.0
 * Contact: apiteam@swagger.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package model

import (
	"context"
	"net/http"
)

// DefaultApiRouter defines the required methods for binding the api requests to a responses for the DefaultApi
// The DefaultApiRouter implementation should parse necessary information from the http request,
// pass the data to a DefaultApiServicer to perform the required actions, then write the service results to the http response.
type DefaultApiRouter interface {
	AddPet(http.ResponseWriter, *http.Request)
	DeletePet(http.ResponseWriter, *http.Request)
	FindPetByID(http.ResponseWriter, *http.Request)
	FindPets(http.ResponseWriter, *http.Request)
}

// DefaultApiServicer defines the api actions for the DefaultApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type DefaultApiServicer interface {
	AddPet(context.Context, NewPet) (ImplResponse, error)
	DeletePet(context.Context, int64) (ImplResponse, error)
	FindPetByID(context.Context, int64) (ImplResponse, error)
	FindPets(context.Context, []string, int32) (ImplResponse, error)
}