package tododriver

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/sagikazarmark/todobackend-go-kit/internal/.generated/api/v1/graphql"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
)

// MakeGraphQLResolver mounts all of the service endpoints into a GraphQL resolver.
func MakeGraphQLResolver(endpoints Endpoints) graphql.ResolverRoot {
	return &resolver{
		endpoints: endpoints,
	}
}

type resolver struct {
	endpoints Endpoints
}

func (r *resolver) Mutation() graphql.MutationResolver {
	return &mutationResolver{r}
}
func (r *resolver) Query() graphql.QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *resolver }

func (r *mutationResolver) AddTodoItem(ctx context.Context, input todo.NewItem) (*todo.Item, error) {
	req := AddItemRequest{
		NewItem: input,
	}

	resp, err := r.endpoints.AddItem(ctx, req)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok && f.Failed() != nil {
		return nil, f.Failed()
	}

	item := resp.(AddItemResponse).Item

	return &item, nil
}

func (r *mutationResolver) UpdateTodoItem(ctx context.Context, input graphql.TodoItemUpdate) (*todo.Item, error) {
	req := UpdateItemRequest{
		Id: input.ID,
		ItemUpdate: todo.ItemUpdate{
			Title:     input.Title,
			Completed: input.Completed,
			Order:     input.Order,
		},
	}

	resp, err := r.endpoints.UpdateItem(ctx, req)
	if err != nil {
		return nil, errors.New("internal server error")
	}

	if f, ok := resp.(endpoint.Failer); ok && f.Failed() != nil {
		return nil, f.Failed()
	}

	item := resp.(UpdateItemResponse).Item

	return &item, nil
}

type queryResolver struct{ *resolver }

func (r *queryResolver) TodoItems(ctx context.Context) ([]todo.Item, error) {
	resp, err := r.endpoints.ListItems(ctx, ListItemsRequest{})
	if err != nil {
		return nil, errors.New("internal server error")
	}

	return resp.(ListItemsResponse).Items, nil
}
