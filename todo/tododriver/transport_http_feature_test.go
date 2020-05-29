package tododriver_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-bdd/gobdd"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/v1/client"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
	"github.com/sagikazarmark/todobackend-go-kit/todo/tododriver"
)

func TestHTTP(t *testing.T) {
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

	suite.AddStep(`an empty todo list`, givenAnEmptyTodoList)
	suite.AddStep(`(?:(?:I|the user)(?: also)? adds? )?(?:a new|an) item for "(.*)"`, addAnItem)
	suite.AddStep(`it should be (?:the only item )?on the list`, shouldBeOnTheList)
	suite.AddStep(`both items should be on the list`, allShouldBeOnTheList)
	suite.AddStep(`the list should be empty`, theListShouldBeEmpty)
	suite.AddStep(`it is marked as complete`, itemMarkedAsComplete)
	suite.AddStep(`it should be complete`, itemShouldBeComplete)
	suite.AddStep(`it is deleted`, deleteAnItem)
	suite.AddStep(`all items are deleted`, clearList)

	suite.Run()
}

func getClient(t gobdd.StepTest, ctx gobdd.Context) *todov1.APIClient {
	v, err := ctx.Get("client")
	if err != nil {
		t.Fatal(err)
	}

	return v.(*todov1.APIClient)
}

func givenAnEmptyTodoList(_ gobdd.StepTest, ctx gobdd.Context) {
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

func addAnItem(t gobdd.StepTest, ctx gobdd.Context, title string) {
	client := getClient(t, ctx)

	item, _, err := client.TodoListApi.AddItem(context.Background(), todov1.AddTodoItemRequest{Title: title})
	if err != nil {
		t.Fatal(err)
	}

	ctx.Set("id", item.Id)
	ctx.Set("title", title)

	ids, _ := ctx.Get("ids", []string{})
	titles, _ := ctx.Get("titles", []string{})

	ctx.Set("ids", append(ids.([]string), item.Id))
	ctx.Set("titles", append(titles.([]string), title))
}

func shouldBeOnTheList(t gobdd.StepTest, ctx gobdd.Context) {
	if err, _ := ctx.GetError("error", nil); err != nil {
		t.Fatal(err)
	}

	client := getClient(t, ctx)

	items, _, err := client.TodoListApi.ListItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	title, _ := ctx.GetString("title", "")

	assert.Len(t, items, 1, "there should be one item on the list")
	assert.Equal(t, items[0].Title, title, "the item on the list should match the added item")
}

func allShouldBeOnTheList(t gobdd.StepTest, ctx gobdd.Context) {
	if err, _ := ctx.GetError("error", nil); err != nil {
		t.Fatal(err)
	}

	client := getClient(t, ctx)

	items, _, err := client.TodoListApi.ListItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}

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

func theListShouldBeEmpty(t gobdd.StepTest, ctx gobdd.Context) {
	client := getClient(t, ctx)

	items, _, err := client.TodoListApi.ListItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(t, items, 0, "the list should be empty")
}

func itemMarkedAsComplete(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getClient(t, ctx)

	completed := true

	_, _, err := client.TodoListApi.UpdateItem(context.Background(), id, todov1.UpdateTodoItemRequest{Completed: &completed})
	if err != nil {
		t.Fatal(err)
	}
}

func itemShouldBeComplete(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getClient(t, ctx)

	item, _, err := client.TodoListApi.GetItem(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, item.Completed, "item should be complete")
}

func deleteAnItem(t gobdd.StepTest, ctx gobdd.Context) {
	id, _ := ctx.GetString("id")

	client := getClient(t, ctx)

	_, err := client.TodoListApi.DeleteItem(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
}

func clearList(t gobdd.StepTest, ctx gobdd.Context) {
	client := getClient(t, ctx)

	_, err := client.TodoListApi.DeleteItems(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
