/*
 * Todo API
 *
 * The Todo API manages a list of todo items as described by the TodoMVC backend project: http://todobackend.com 
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"context"
	"net/http"
	"errors"
)

// TodoListApiService is a service that implents the logic for the TodoListApiServicer
// This service should implement the business logic for every endpoint for the TodoListApi API. 
// Include any external packages or services that will be required by this service.
type TodoListApiService struct {
}

// NewTodoListApiService creates a default api service
func NewTodoListApiService() TodoListApiServicer {
	return &TodoListApiService{}
}

// AddItem - Add a new item to the list
func (s *TodoListApiService) AddItem(ctx context.Context, addTodoItemRequest AddTodoItemRequest) (ImplResponse, error) {
	// TODO - update AddItem with the required logic for this service method.
	// Add api_todo_list_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, TodoItem{}) or use other options such as http.Ok ...
	//return Response(201, TodoItem{}), nil

	//TODO: Uncomment the next line to return response Response(400, Error{}) or use other options such as http.Ok ...
	//return Response(400, Error{}), nil

	//TODO: Uncomment the next line to return response Response(422, Error{}) or use other options such as http.Ok ...
	//return Response(422, Error{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("AddItem method not implemented")
}

// DeleteItem - Delete an item
func (s *TodoListApiService) DeleteItem(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update DeleteItem with the required logic for this service method.
	// Add api_todo_list_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	//return Response(204, nil),nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteItem method not implemented")
}

// DeleteItems - Delete all items
func (s *TodoListApiService) DeleteItems(ctx context.Context) (ImplResponse, error) {
	// TODO - update DeleteItems with the required logic for this service method.
	// Add api_todo_list_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	//return Response(204, nil),nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("DeleteItems method not implemented")
}

// GetItem - Get an item
func (s *TodoListApiService) GetItem(ctx context.Context, id string) (ImplResponse, error) {
	// TODO - update GetItem with the required logic for this service method.
	// Add api_todo_list_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, TodoItem{}) or use other options such as http.Ok ...
	//return Response(200, TodoItem{}), nil

	//TODO: Uncomment the next line to return response Response(404, Error{}) or use other options such as http.Ok ...
	//return Response(404, Error{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetItem method not implemented")
}

// ListItems - List items
func (s *TodoListApiService) ListItems(ctx context.Context) (ImplResponse, error) {
	// TODO - update ListItems with the required logic for this service method.
	// Add api_todo_list_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []TodoItem{}) or use other options such as http.Ok ...
	//return Response(200, []TodoItem{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("ListItems method not implemented")
}

// UpdateItem - Update an existing item
func (s *TodoListApiService) UpdateItem(ctx context.Context, id string, updateTodoItemRequest UpdateTodoItemRequest) (ImplResponse, error) {
	// TODO - update UpdateItem with the required logic for this service method.
	// Add api_todo_list_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, TodoItem{}) or use other options such as http.Ok ...
	//return Response(200, TodoItem{}), nil

	//TODO: Uncomment the next line to return response Response(404, Error{}) or use other options such as http.Ok ...
	//return Response(404, Error{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("UpdateItem method not implemented")
}
