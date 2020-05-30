package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/spf13/pflag"

	"github.com/sagikazarmark/todobackend-go-kit/todo"
	"github.com/sagikazarmark/todobackend-go-kit/todo/tododriver"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version    string
	commitHash string
	buildDate  string
)

func main() {
	flags := pflag.NewFlagSet("Go kit TodoBackend", pflag.ExitOnError)

	httpAddr := flags.String("http-addr", ":8000", "HTTP Server address")
	publicURL := flags.String("public-url", "http://localhost:8000", "Publicly available base URL")

	_ = flags.Parse(os.Args[1:])

	log.Println("starting application version", version, fmt.Sprintf("(%s)", commitHash), "built on", buildDate)

	todoURL := *publicURL + "/todos"

	router := mux.NewRouter()

	{
		file, err := pkger.Open("/static/index.html")
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		r := strings.NewReplacer("PUBLIC_URL", todoURL, "VERSION", version)

		body = []byte(r.Replace(string(body)))

		router.Methods(http.MethodGet).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write(body)
		})
	}

	{
		store := todo.NewInMemoryStore()
		service := todo.NewService(ulidgen.NewGenerator(), store)
		endpoints := tododriver.MakeEndpoints(service)

		tododriver.RegisterHTTPHandlers(
			endpoints,
			router.PathPrefix("/todos").Subrouter(),
			kithttp.ServerBefore(func(ctx context.Context, request *http.Request) context.Context {
				return context.WithValue(ctx, tododriver.ContextKeyBaseURL, todoURL)
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

	log.Println("listening on", *httpAddr)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
