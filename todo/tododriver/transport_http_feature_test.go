package tododriver_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-bdd/gobdd"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1/client/rest"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
	"github.com/sagikazarmark/todobackend-go-kit/todo/tododriver"
)

func TestRest(t *testing.T) {
	suite := gobdd.NewSuite(
		t,
		gobdd.WithFeaturesPath("../features/*.feature"),
		gobdd.WithAfterScenario(func(ctx gobdd.Context) {
			server, _ := ctx.Get("server", nil)
			if server == nil {
				return
			}

			server.(*httptest.Server).Close()
		}),
	)

	suite.AddStep(`an empty todo list`, givenAnEmptyTodoListRest)
	suite.AddStep(`(?:(?:I|the user)(?: also)? adds? )?(?:a new|an) item for {text}`, addAnItemRest)
	suite.AddStep(`it should be (?:the only item )?on the list`, shouldBeOnTheRest)
	suite.AddStep(`both items should be on the list`, allShouldBeOnTheListRest)
	suite.AddStep(`the list should be empty`, theListShouldBeEmptyRest)
	suite.AddStep(`it is marked as complete`, itemMarkedAsCompleteRest)
	suite.AddStep(`it should be complete`, itemShouldBeCompleteRest)
	suite.AddStep(`it is deleted`, deleteAnItemRest)
	suite.AddStep(`all items are deleted`, clearListRest)

	suite.Run()
}

func getRestClient(t gobdd.StepTest, ctx gobdd.Context) *todov1.APIClient {
	v, err := ctx.Get("client")
	require.NoError(t, err)

	return v.(*todov1.APIClient)
}

func givenAnEmptyTodoListRest(_ gobdd.StepTest, ctx gobdd.Context) {
	store := todo.NewInMemoryStore()
	service := todo.NewService(ulidgen.NewGenerator(), store)
	endpoints := tododriver.MakeEndpoints(service)

	router := mux.NewRouter()

	tododriver.RegisterHTTPHandlers(endpoints, router.PathPrefix("/todos").Subrouter())

	server := httptest.NewServer(router)

	ctx.Set("server", server)

	config := todov1.NewConfiguration()
	config.BasePath = server.URL
	config.HTTPClient = server.Client()

	ctx.Set("client", todov1.NewAPIClient(config))
}

func addAnItemRest(t gobdd.StepTest, ctx gobdd.Context, title string) {
	client := getRestClient(t, ctx)

	item, _, err := client.TodoListApi.AddItem(context.Background(), todov1.AddTodoItemRequest{Title: title})
	require.NoError(t, err)

	ctx.Set("id", item.Id)
	ctx.Set("title", title)

	ids, _ := ctx.Get("ids", []string{})
	titles, _ := ctx.Get("titles", []string{})

	ctx.Set("ids", append(ids.([]string), item.Id))
	ctx.Set("titles", append(titles.([]string), title))
}

func shouldBeOnTheRest(t gobdd.StepTest, ctx gobdd.Context) {
	if err, _ := ctx.GetError("error", nil); err != nil {
		t.Fatal(err)
	}

	client := getRestClient(t, ctx)

	items, _, err := client.TodoListApi.ListItems(context.Background())
	require.NoError(t, err)

	title, _ := ctx.GetString("title", "")

	assert.Len(t, items, 1, "there should be one item on the list")
	assert.Equal(t, items[0].Title, title, "the item on the list should match the added item")
}

func allShouldBeOnTheListRest(t gobdd.StepTest, ctx gobdd.Context) {
	if err, _ := ctx.GetError("error", nil); err != nil {
		t.Fatal(err)
	}

	client := getRestClient(t, ctx)

	items, _, err := client.TodoListApi.ListItems(context.Background())
	require.NoError(t, err)

	ids, _ := ctx.Get("ids", []string{})
	titles, _ := ctx.Get("titles", []string{})

	idMap := make(map[string]string)

	for i, id := range ids.([]string) {
		idMap[id] = titles.([]string)[i]
	}

	assert.Len(t, items, len(idMap))

	for _, item := range items {
		assert.Equal(t, idMap[item.Id], item.Title, "the item on the list should match the added item")
	}
}

func theListShouldBeEmptyRest(t gobdd.StepTest, ctx gobdd.Context) {
	client := getRestClient(t, ctx)

	items, _, err := client.TodoListApi.ListItems(context.Background())
	require.NoError(t, err)

	assert.Len(t, items, 0, "the list should be empty")
}

func itemMarkedAsCompleteRest(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getRestClient(t, ctx)

	completed := true

	_, _, err := client.TodoListApi.UpdateItem(context.Background(), id, todov1.UpdateTodoItemRequest{Completed: &completed}) // nolint: lll
	require.NoError(t, err)
}

func itemShouldBeCompleteRest(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getRestClient(t, ctx)

	item, _, err := client.TodoListApi.GetItem(context.Background(), id)
	require.NoError(t, err)

	assert.True(t, item.Completed, "item should be complete")
}

func deleteAnItemRest(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getRestClient(t, ctx)

	_, err := client.TodoListApi.DeleteItem(context.Background(), id)
	require.NoError(t, err)
}

func clearListRest(t gobdd.StepTest, ctx gobdd.Context) {
	client := getRestClient(t, ctx)

	_, err := client.TodoListApi.DeleteItems(context.Background())
	require.NoError(t, err)
}
