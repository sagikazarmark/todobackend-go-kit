package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
	"github.com/oklog/run"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
	"github.com/sagikazarmark/todobackend-go-kit/internal/.generated/api/v1/graphql"
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
	grpcAddr := flags.String("grpc-addr", ":8001", "gRPC Server address")
	publicURL := flags.String("public-url", "http://localhost:8000", "Publicly available base URL")

	_ = flags.Parse(os.Args[1:])

	log.Println("starting application version", version, fmt.Sprintf("(%s)", commitHash), "built on", buildDate)

	todoURL := *publicURL + "/todos"

	router := mux.NewRouter()

	grpcServer := grpc.NewServer()

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
		todov1.RegisterTodoListServiceServer(grpcServer, tododriver.MakeGRPCServer(endpoints))
		router.PathPrefix("/graphql").Handler(handler.NewDefaultServer(
			graphql.NewExecutableSchema(graphql.Config{
				Resolvers: tododriver.MakeGraphQLResolver(endpoints),
			}),
		))
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete}),
		handlers.AllowedHeaders([]string{"content-type"}),
	)

	httpServer := &http.Server{
		Addr:    *httpAddr,
		Handler: cors(router),
	}

	log.Println("listening on", *httpAddr)

	httpLn, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on", *grpcAddr)

	grpcLn, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	var group run.Group

	group.Add(
		func() error { return httpServer.Serve(httpLn) },
		func(err error) { _ = httpServer.Shutdown(context.Background()) },
	)
	defer httpServer.Close()

	group.Add(
		func() error { return grpcServer.Serve(grpcLn) },
		func(err error) { grpcServer.GracefulStop() },
	)
	defer grpcServer.Stop()

	// Setup signal handler
	group.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	err = group.Run()
	if err != nil {
		if _, ok := err.(run.SignalError); ok {
			log.Println(err)

			return
		}

		// Fatal error
		// We don't use fatal, so deferred functions can do their jobs.
		log.Println(err)
	}
}
