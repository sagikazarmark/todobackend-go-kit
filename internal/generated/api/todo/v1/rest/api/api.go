// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * Todo API
 *
 * The Todo API manages a list of todo items as described by the TodoMVC backend project: http://todobackend.com 
 *
 * API version: 1.0.0
 */

package api

import (
	"context"
	"net/http"
)



// TodoListAPIRouter defines the required methods for binding the api requests to a responses for the TodoListAPI
// The TodoListAPIRouter implementation should parse necessary information from the http request,
// pass the data to a TodoListAPIServicer to perform the required actions, then write the service results to the http response.
type TodoListAPIRouter interface { 
	AddItem(http.ResponseWriter, *http.Request)
	DeleteItem(http.ResponseWriter, *http.Request)
	DeleteItems(http.ResponseWriter, *http.Request)
	GetItem(http.ResponseWriter, *http.Request)
	ListItems(http.ResponseWriter, *http.Request)
	UpdateItem(http.ResponseWriter, *http.Request)
}


// TodoListAPIServicer defines the api actions for the TodoListAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type TodoListAPIServicer interface { 
	AddItem(context.Context, AddTodoItemRequest) (ImplResponse, error)
	DeleteItem(context.Context, string) (ImplResponse, error)
	DeleteItems(context.Context) (ImplResponse, error)
	GetItem(context.Context, string) (ImplResponse, error)
	ListItems(context.Context) (ImplResponse, error)
	UpdateItem(context.Context, string, UpdateTodoItemRequest) (ImplResponse, error)
}
