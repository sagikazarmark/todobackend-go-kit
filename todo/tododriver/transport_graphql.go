package tododriver

import (
	"context"

	graphql2 "github.com/99designs/gqlgen/graphql"
	kitxgraphql "github.com/sagikazarmark/kitx/transport/graphql"

	"github.com/sagikazarmark/todobackend-go-kit/internal/.generated/api/v1/graphql"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
)

// MakeGraphQLSchema mounts all of the service endpoints into a GraphQL executable schema.
func MakeGraphQLSchema(endpoints Endpoints, options ...kitxgraphql.ServerOption) graphql2.ExecutableSchema {
	return graphql.NewExecutableSchema(graphql.Config{
		Resolvers: MakeGraphQLResolver(endpoints, options...),
	})
}

// MakeGraphQLResolver mounts all of the service endpoints into a GraphQL resolver.
func MakeGraphQLResolver(endpoints Endpoints, options ...kitxgraphql.ServerOption) graphql.ResolverRoot {
	errorEncoder := func(_ context.Context, err error) error {
		return err
	}

	return &resolver{
		AddTodoItemHandler: kitxgraphql.NewErrorEncoderHandler(kitxgraphql.NewServer(
			endpoints.AddItem,
			decodeAddItemGraphQLRequest,
			kitxgraphql.ErrorResponseEncoder(encodeAddItemGraphQLResponse, errorEncoder),
			options...,
		), errorEncoder),
		UpdateTodoItemHandler: kitxgraphql.NewErrorEncoderHandler(kitxgraphql.NewServer(
			endpoints.UpdateItem,
			decodeUpdateItemGraphQLRequest,
			kitxgraphql.ErrorResponseEncoder(encodeUpdateItemGraphQLResponse, errorEncoder),
			options...,
		), errorEncoder),
		ListTodoItemsHandler: kitxgraphql.NewErrorEncoderHandler(kitxgraphql.NewServer(
			endpoints.ListItems,
			decodeListItemsGraphQLRequest,
			kitxgraphql.ErrorResponseEncoder(encodeListItemsGraphQLResponse, errorEncoder),
			options...,
		), errorEncoder),
	}
}

func decodeAddItemGraphQLRequest(_ context.Context, request interface{}) (interface{}, error) {
	return AddItemRequest{
		NewItem: request.(todo.NewItem),
	}, nil
}

func encodeAddItemGraphQLResponse(_ context.Context, response interface{}) (interface{}, error) {
	item := response.(AddItemResponse).Item

	return &item, nil
}

func decodeUpdateItemGraphQLRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(graphql.TodoItemUpdate)

	return UpdateItemRequest{
		Id: req.ID,
		ItemUpdate: todo.ItemUpdate{
			Title:     req.Title,
			Completed: req.Completed,
			Order:     req.Order,
		},
	}, nil
}

func encodeUpdateItemGraphQLResponse(_ context.Context, response interface{}) (interface{}, error) {
	item := response.(UpdateItemResponse).Item

	return &item, nil
}

func decodeListItemsGraphQLRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return ListItemsRequest{}, nil
}

func encodeListItemsGraphQLResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response.(ListItemsResponse).Items, nil
}

type resolver struct {
	AddTodoItemHandler    kitxgraphql.Handler
	UpdateTodoItemHandler kitxgraphql.Handler
	ListTodoItemsHandler  kitxgraphql.Handler
}

func (r *resolver) Mutation() graphql.MutationResolver {
	return &mutationResolver{r}
}

func (r *resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *resolver }

func (r *mutationResolver) AddTodoItem(ctx context.Context, input todo.NewItem) (*todo.Item, error) {
	_, resp, err := r.AddTodoItemHandler.ServeGraphQL(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.(*todo.Item), nil
}

func (r *mutationResolver) UpdateTodoItem(ctx context.Context, input graphql.TodoItemUpdate) (*todo.Item, error) {
	_, resp, err := r.UpdateTodoItemHandler.ServeGraphQL(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.(*todo.Item), nil
}

type queryResolver struct{ *resolver }

func (r *queryResolver) TodoItems(ctx context.Context) ([]todo.Item, error) {
	_, resp, err := r.ListTodoItemsHandler.ServeGraphQL(ctx, nil)
	if err != nil {
		return nil, err
	}

	return resp.([]todo.Item), nil
}
