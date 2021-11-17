package tododriver

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	appkitgrpc "github.com/sagikazarmark/appkit/transport/grpc"
	kitxgrpc "github.com/sagikazarmark/kitx/transport/grpc"

	api "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, options ...kitgrpc.ServerOption) api.TodoListServiceServer {
	errorEncoder := kitxgrpc.NewStatusErrorResponseEncoder(appkitgrpc.NewDefaultStatusConverter())

	return api.TodoListServiceKitServer{
		AddItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.AddItem,
			decodeAddItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeAddItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		ListItemsHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.ListItems,
			decodeListItemsGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeListItemsGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		DeleteItemsHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.DeleteItems,
			decodeDeleteItemsGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeDeleteItemsGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		GetItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.GetItem,
			decodeGetItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeGetItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		UpdateItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.UpdateItem,
			decodeUpdateItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeUpdateItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
		DeleteItemHandler: kitxgrpc.NewErrorEncoderHandler(kitgrpc.NewServer(
			endpoints.DeleteItem,
			decodeDeleteItemGRPCRequest,
			kitxgrpc.ErrorResponseEncoder(encodeDeleteItemGRPCResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeAddItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*api.AddItemRequest)

	return AddItemRequest{
		NewItem: todo.NewItem{
			Title: req.GetTitle(),
			Order: int(req.GetOrder()),
		},
	}, nil
}

func encodeAddItemGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(AddItemResponse)

	return &api.AddItemResponse{
		Item: marshalItemGRPC(resp.Item),
	}, nil
}

func decodeListItemsGRPCRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return ListItemsRequest{}, nil
}

func encodeListItemsGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(ListItemsResponse)

	items := make([]*api.TodoItem, 0, len(resp.Items))

	for _, item := range resp.Items {
		items = append(items, marshalItemGRPC(item))
	}

	return &api.ListItemsResponse{
		Items: items,
	}, nil
}

func decodeDeleteItemsGRPCRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return DeleteItemsRequest{}, nil
}

func encodeDeleteItemsGRPCResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &api.DeleteItemsResponse{}, nil
}

func decodeGetItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*api.GetItemRequest)

	return GetItemRequest{
		Id: req.GetId(),
	}, nil
}

func encodeGetItemGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetItemResponse)

	return &api.GetItemResponse{
		Item: marshalItemGRPC(resp.Item),
	}, nil
}

func decodeUpdateItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*api.UpdateItemRequest)

	var (
		title     *string
		completed *bool
		order     *int
	)

	if req.Title != nil {
		title = &req.Title.Value
	}

	if req.Completed != nil {
		completed = &req.Completed.Value
	}

	if req.Order != nil {
		o := int(req.Order.Value)
		order = &o
	}

	return UpdateItemRequest{
		Id: req.GetId(),
		ItemUpdate: todo.ItemUpdate{
			Title:     title,
			Completed: completed,
			Order:     order,
		},
	}, nil
}

func encodeUpdateItemGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(UpdateItemResponse)

	return &api.UpdateItemResponse{
		Item: marshalItemGRPC(resp.Item),
	}, nil
}

func decodeDeleteItemGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*api.DeleteItemRequest)

	return DeleteItemRequest{
		Id: req.GetId(),
	}, nil
}

func encodeDeleteItemGRPCResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &api.DeleteItemResponse{}, nil
}

func marshalItemGRPC(item todo.Item) *api.TodoItem {
	return &api.TodoItem{
		Id:        item.ID,
		Title:     item.Title,
		Completed: item.Completed,
		Order:     int32(item.Order),
	}
}
