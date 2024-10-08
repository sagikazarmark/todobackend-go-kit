package tododriver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	appkithttp "github.com/sagikazarmark/appkit/transport/http"
	kitxhttp "github.com/sagikazarmark/kitx/transport/http"

	"github.com/sagikazarmark/todobackend-go-kit/internal/generated/api/todo/v1/rest/api"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
)

type contextKey int

const (
	ContextKeyBaseURL contextKey = iota
)

// MakeHTTPHandler mounts all of the service endpoints into an {http.Handler}.
func MakeHTTPHandler(endpoints Endpoints, options ...kithttp.ServerOption) http.Handler {
	router := chi.NewRouter()
	errorEncoder := kitxhttp.NewJSONProblemErrorResponseEncoder(appkithttp.NewDefaultProblemConverter())

	router.Method(http.MethodPost, "/", kithttp.NewServer(
		endpoints.AddItem,
		decodeAddItemHTTPRequest,
		kitxhttp.ErrorResponseEncoder(encodeAddItemHTTPResponse, errorEncoder),
		options...,
	))

	router.Method(http.MethodGet, "/", kithttp.NewServer(
		endpoints.ListItems,
		kithttp.NopRequestDecoder,
		kitxhttp.ErrorResponseEncoder(encodeListItemsHTTPResponse, errorEncoder),
		options...,
	))

	router.Method(http.MethodDelete, "/", kithttp.NewServer(
		endpoints.DeleteItems,
		kithttp.NopRequestDecoder,
		kitxhttp.ErrorResponseEncoder(kitxhttp.StatusCodeResponseEncoder(http.StatusNoContent), errorEncoder),
		options...,
	))

	router.Method(http.MethodGet, "/{id}", kithttp.NewServer(
		endpoints.GetItem,
		decodeGetItemHTTPRequest,
		kitxhttp.ErrorResponseEncoder(encodeGetItemHTTPResponse, errorEncoder),
		options...,
	))

	router.Method(http.MethodPatch, "/{id}", kithttp.NewServer(
		endpoints.UpdateItem,
		decodeUpdateItemHTTPRequest,
		kitxhttp.ErrorResponseEncoder(encodeUpdateItemHTTPResponse, errorEncoder),
		options...,
	))

	router.Method(http.MethodDelete, "/{id}", kithttp.NewServer(
		endpoints.DeleteItem,
		decodeDeleteItemHTTPRequest,
		kitxhttp.ErrorResponseEncoder(kitxhttp.StatusCodeResponseEncoder(http.StatusNoContent), errorEncoder),
		options...,
	))

	return router
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
		Order:     int32(item.Order), //nolint
		Url:       fmt.Sprintf("%s/%s", baseURL, item.ID),
	}
}

func getIDParamFromRequest(r *http.Request) (string, error) {
	id := chi.URLParam(r, "id")
	if id == "" {
		return "", errors.NewWithDetails("missing parameter from the URL", "param", "id")
	}

	return id, nil
}
