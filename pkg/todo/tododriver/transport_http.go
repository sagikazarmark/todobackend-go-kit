package tododriver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"emperror.dev/errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	appkithttp "github.com/sagikazarmark/appkit/transport/http"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"

	api "github.com/sagikazarmark/todobackend-go-kit/internal/.generated/api/v1/rest/go"
	"github.com/sagikazarmark/todobackend-go-kit/pkg/todo"
)

type contextKey int

const (
	ContextKeyBaseURL contextKey = iota
)

// RegisterHTTPHandlers mounts all of the service endpoints into a router.
func RegisterHTTPHandlers(endpoints Endpoints, router *mux.Router, options ...kithttp.ServerOption) {
	errorEncoder := kitxhttp.NewJSONProblemErrorResponseEncoder(appkithttp.NewDefaultProblemConverter())

	router.Methods(http.MethodPost).Path("").Handler(kithttp.NewServer(
		endpoints.AddItem,
		decodeAddItemHTTPRequest,
		kitxhttp.ErrorResponseEncoder(encodeAddItemHTTPResponse, errorEncoder),
		options...,
	))

	router.Methods(http.MethodGet).Path("").Handler(kithttp.NewServer(
		endpoints.ListItems,
		kithttp.NopRequestDecoder,
		kitxhttp.ErrorResponseEncoder(encodeListItemsHTTPResponse, errorEncoder),
		options...,
	))

	router.Methods(http.MethodDelete).Path("").Handler(kithttp.NewServer(
		endpoints.DeleteItems,
		kithttp.NopRequestDecoder,
		kitxhttp.ErrorResponseEncoder(kitxhttp.StatusCodeResponseEncoder(http.StatusNoContent), errorEncoder),
		options...,
	))

	router.Methods(http.MethodGet).Path("/{id}").Handler(kithttp.NewServer(
		endpoints.GetItem,
		decodeGetItemHTTPRequest,
		kitxhttp.ErrorResponseEncoder(encodeGetItemHTTPResponse, errorEncoder),
		options...,
	))

	router.Methods(http.MethodPatch).Path("/{id}").Handler(kithttp.NewServer(
		endpoints.UpdateItem,
		decodeUpdateItemHTTPRequest,
		kitxhttp.ErrorResponseEncoder(encodeUpdateItemHTTPResponse, errorEncoder),
		options...,
	))

	router.Methods(http.MethodDelete).Path("/{id}").Handler(kithttp.NewServer(
		endpoints.DeleteItem,
		decodeDeleteItemHTTPRequest,
		kitxhttp.ErrorResponseEncoder(kitxhttp.StatusCodeResponseEncoder(http.StatusNoContent), errorEncoder),
		options...,
	))
}

func decodeAddItemHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var apiRequest api.AddTodoItemRequest

	err := json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		return nil, errors.Wrap(err, "decode request")
	}

	return AddItemRequest{
		NewItem: todo.NewItem{
			Title: apiRequest.Title,
			Order: int(apiRequest.Order),
		},
	}, nil
}

func encodeAddItemHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(AddItemResponse)

	apiResponse := marshalItemHTTP(ctx, resp.Item)

	return kitxhttp.JSONResponseEncoder(ctx, w, kitxhttp.WithStatusCode(apiResponse, http.StatusCreated))
}

func encodeListItemsHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(ListItemsResponse)

	items := make([]api.TodoItem, 0, len(resp.Items))

	for _, item := range resp.Items {
		items = append(items, marshalItemHTTP(ctx, item))
	}

	return kitxhttp.JSONResponseEncoder(ctx, w, items)
}

func decodeGetItemHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := getIDParamFromRequest(r)
	if err != nil {
		return nil, err
	}

	return GetItemRequest{
		Id: id,
	}, nil
}

func encodeGetItemHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(GetItemResponse)

	apiResponse := marshalItemHTTP(ctx, resp.Item)

	return kitxhttp.JSONResponseEncoder(ctx, w, apiResponse)
}

func decodeUpdateItemHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := getIDParamFromRequest(r)
	if err != nil {
		return nil, err
	}

	var apiRequest api.UpdateTodoItemRequest

	err = json.NewDecoder(r.Body).Decode(&apiRequest)
	if err != nil {
		return nil, errors.Wrap(err, "decode request")
	}

	var order *int

	if apiRequest.Order != nil {
		o := int(*apiRequest.Order)
		order = &o
	}

	return UpdateItemRequest{
		Id: id,
		ItemUpdate: todo.ItemUpdate{
			Title:     apiRequest.Title,
			Completed: apiRequest.Completed,
			Order:     order,
		},
	}, nil
}

func encodeUpdateItemHTTPResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(UpdateItemResponse)

	apiResponse := marshalItemHTTP(ctx, resp.Item)

	return kitxhttp.JSONResponseEncoder(ctx, w, apiResponse)
}

func decodeDeleteItemHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := getIDParamFromRequest(r)
	if err != nil {
		return nil, err
	}

	return DeleteItemRequest{
		Id: id,
	}, nil
}

func marshalItemHTTP(ctx context.Context, item todo.Item) api.TodoItem {
	baseURL, _ := ctx.Value(ContextKeyBaseURL).(string)

	return api.TodoItem{
		Id:        item.ID,
		Title:     item.Title,
		Completed: item.Completed,
		Order:     int32(item.Order),
		Url:       fmt.Sprintf("%s/%s", baseURL, item.ID),
	}
}

func getIDParamFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok || id == "" {
		return "", errors.NewWithDetails("missing parameter from the URL", "param", "id")
	}

	return id, nil
}
