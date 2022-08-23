package tododriver_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/go-bdd/gobdd"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/goph/idgen/ulidgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
	"github.com/sagikazarmark/todobackend-go-kit/todo/tododriver"
)

func TestGRPC(t *testing.T) {
	suite := gobdd.NewSuite(
		t,
		gobdd.WithFeaturesPath("../features/*.feature"),
		gobdd.WithAfterScenario(func(ctx gobdd.Context) {
			server, _ := ctx.Get("server", nil)
			if server != nil {
				server.(*grpc.Server).Stop()
			}

			conn, _ := ctx.Get("conn", nil)
			if conn != nil {
				_ = conn.(*grpc.ClientConn).Close()
			}
		}),
	)

	suite.AddStep(`an empty todo list`, givenAnEmptyTodoListGRPC)
	suite.AddStep(`(?:(?:I|the user)(?: also)? adds? )?(?:a new|an) item for {text}`, addAnItemGRPC)
	suite.AddStep(`it should be (?:the only item )?on the list`, shouldBeOnTheGRPC)
	suite.AddStep(`both items should be on the list`, allShouldBeOnTheListGRPC)
	suite.AddStep(`the list should be empty`, theListShouldBeEmptyGRPC)
	suite.AddStep(`it is marked as complete`, itemMarkedAsCompleteGRPC)
	suite.AddStep(`it should be complete`, itemShouldBeCompleteGRPC)
	suite.AddStep(`it is deleted`, deleteAnItemGRPC)
	suite.AddStep(`all items are deleted`, clearListGRPC)

	suite.Run()
}

func getGRPCClient(t gobdd.StepTest, ctx gobdd.Context) todov1.TodoListServiceClient {
	v, err := ctx.Get("client")
	if err != nil {
		t.Fatal(err)
	}

	return v.(todov1.TodoListServiceClient)
}

func newLocalListener() (net.Listener, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			return nil, fmt.Errorf("grpctest: failed to listen on a port: %w", err)
		}
	}

	return l, nil
}

func givenAnEmptyTodoListGRPC(t gobdd.StepTest, ctx gobdd.Context) {
	l, err := newLocalListener()
	require.NoError(t, err)

	store := todo.NewInMemoryStore()
	service := todo.NewService(ulidgen.NewGenerator(), store)
	endpoints := tododriver.MakeEndpoints(service)

	server := grpc.NewServer()

	todov1.RegisterTodoListServiceServer(server, tododriver.MakeGRPCServer(endpoints))

	go func() {
		_ = server.Serve(l)
	}()

	ctx.Set("server", server)

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		server.Stop()

		t.Fatal(err)
	}

	ctx.Set("conn", conn)

	client := todov1.NewTodoListServiceClient(conn)

	ctx.Set("client", client)
}

func addAnItemGRPC(t gobdd.StepTest, ctx gobdd.Context, title string) {
	client := getGRPCClient(t, ctx)

	resp, err := client.AddItem(context.Background(), &todov1.AddItemRequest{Title: title})
	require.NoError(t, err)

	ctx.Set("id", resp.GetItem().GetId())
	ctx.Set("title", title)

	ids, _ := ctx.Get("ids", []string{})
	titles, _ := ctx.Get("titles", []string{})

	ctx.Set("ids", append(ids.([]string), resp.GetItem().GetId()))
	ctx.Set("titles", append(titles.([]string), title))
}

func shouldBeOnTheGRPC(t gobdd.StepTest, ctx gobdd.Context) {
	if err, _ := ctx.GetError("error", nil); err != nil {
		t.Fatal(err)
	}

	client := getGRPCClient(t, ctx)

	resp, err := client.ListItems(context.Background(), new(todov1.ListItemsRequest))
	require.NoError(t, err)

	title, _ := ctx.GetString("title", "")

	assert.Len(t, resp.GetItems(), 1, "there should be one item on the list")
	assert.Equal(t, resp.GetItems()[0].GetTitle(), title, "the item on the list should match the added item")
}

func allShouldBeOnTheListGRPC(t gobdd.StepTest, ctx gobdd.Context) {
	if err, _ := ctx.GetError("error", nil); err != nil {
		t.Fatal(err)
	}

	client := getGRPCClient(t, ctx)

	resp, err := client.ListItems(context.Background(), new(todov1.ListItemsRequest))
	require.NoError(t, err)

	ids, _ := ctx.Get("ids", []string{})
	titles, _ := ctx.Get("titles", []string{})

	idMap := make(map[string]string)

	for i, id := range ids.([]string) {
		idMap[id] = titles.([]string)[i]
	}

	assert.Len(t, resp.GetItems(), len(idMap))

	for _, item := range resp.GetItems() {
		assert.Equal(t, idMap[item.GetId()], item.GetTitle(), "the item on the list should match the added item")
	}
}

func theListShouldBeEmptyGRPC(t gobdd.StepTest, ctx gobdd.Context) {
	client := getGRPCClient(t, ctx)

	resp, err := client.ListItems(context.Background(), new(todov1.ListItemsRequest))
	require.NoError(t, err)

	assert.Len(t, resp.GetItems(), 0, "the list should be empty")
}

func itemMarkedAsCompleteGRPC(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getGRPCClient(t, ctx)

	_, err := client.UpdateItem(context.Background(), &todov1.UpdateItemRequest{Id: id, Completed: &wrappers.BoolValue{Value: true}}) //nolint: lll
	require.NoError(t, err)
}

func itemShouldBeCompleteGRPC(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getGRPCClient(t, ctx)

	resp, err := client.GetItem(context.Background(), &todov1.GetItemRequest{Id: id})
	require.NoError(t, err)

	assert.True(t, resp.GetItem().GetCompleted(), "item should be complete")
}

func deleteAnItemGRPC(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getGRPCClient(t, ctx)

	_, err := client.DeleteItem(context.Background(), &todov1.DeleteItemRequest{Id: id})
	require.NoError(t, err)
}

func clearListGRPC(t gobdd.StepTest, ctx gobdd.Context) {
	client := getGRPCClient(t, ctx)

	_, err := client.DeleteItems(context.Background(), new(todov1.DeleteItemsRequest))
	require.NoError(t, err)
}
