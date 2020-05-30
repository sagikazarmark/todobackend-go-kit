package main

import (
	"context"
	"log"
	"net/http"
	"os"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"

	"github.com/sagikazarmark/todobackend-go-kit/todo"
	"github.com/sagikazarmark/todobackend-go-kit/todo/tododriver"
)

func main() {
	flags := pflag.NewFlagSet("Go kit TodoBackend", pflag.ExitOnError)

	httpAddr := flags.String("http-addr", ":8000", "HTTP Server address")
	publicURL := flags.String("public-url", "http://localhost:8000", "Publicly available base URL")

	_ = flags.Parse(os.Args[1:])

	router := mux.NewRouter()

	{
		store := todo.NewInMemoryStore()
		service := todo.NewService(ulidgen.NewGenerator(), store)
		endpoints := tododriver.MakeEndpoints(service)

		tododriver.RegisterHTTPHandlers(
			endpoints,
			router.PathPrefix("/todos").Subrouter(),
			kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
				return context.WithValue(ctx, tododriver.ContextKeyBaseURL, *publicURL+"/todos")
			}),
		)
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete}),
		handlers.AllowedHeaders([]string{"content-type"}),
	)

	server := &http.Server{
		Addr:    *httpAddr,
		Handler: cors(router),
	}

	log.Println("starting application at", *httpAddr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
