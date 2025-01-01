package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/goph/idgen/ulidgen"
	"github.com/oklog/run"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	todov1 "github.com/sagikazarmark/todobackend-go-kit/api/todo/v1"
	"github.com/sagikazarmark/todobackend-go-kit/internal/generated/api/todo/v1/graphql"
	"github.com/sagikazarmark/todobackend-go-kit/static"
	"github.com/sagikazarmark/todobackend-go-kit/todo"
	"github.com/sagikazarmark/todobackend-go-kit/todo/tododriver"
)

func main() {
	flags := pflag.NewFlagSet("Go kit TodoBackend", pflag.ExitOnError)

	httpAddr := flags.String("http-addr", ":8000", "HTTP Server address")
	grpcAddr := flags.String("grpc-addr", ":8001", "gRPC Server address")
	publicURL := flags.String("public-url", "http://localhost:8000", "Publicly available base URL")

	_ = flags.Parse(os.Args[1:])

	log.Println("starting application version", version, fmt.Sprintf("(%s)", revision[:8]), "built on", revisionDate)

	todoURL := *publicURL + "/todos"

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	grpcServer := grpc.NewServer()

	{
		index, err := static.Files().ReadFile("index.html")
		if err != nil {
			panic(err)
		}

		r := strings.NewReplacer("PUBLIC_URL", todoURL, "VERSION", version)

		body := []byte(r.Replace(string(index)))

		router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)

			_, _ = w.Write(body)
		})
	}

	{
		store := todo.NewInMemoryStore()
		service := todo.NewService(ulidgen.NewGenerator(), store)
		endpoints := tododriver.MakeEndpoints(service)

		router.Mount("/todos", tododriver.MakeHTTPHandler(
			endpoints,
			kithttp.ServerBefore(func(ctx context.Context, _ *http.Request) context.Context {
				return context.WithValue(ctx, tododriver.ContextKeyBaseURL, todoURL)
			}),
		))
		todov1.RegisterTodoListServiceServer(grpcServer, tododriver.MakeGRPCServer(endpoints))
		router.Handle("/graphql/playground", playground.Handler("GraphQL playground", "/graphql/query"))
		router.Handle("/graphql/query", handler.NewDefaultServer( //nolint
			graphql.NewExecutableSchema(graphql.Config{
				Resolvers: tododriver.MakeGraphQLResolver(endpoints),
			}),
		))
	}

	httpServer := &http.Server{
		Addr:              *httpAddr,
		ReadHeaderTimeout: 30 * time.Second,
		Handler:           router,
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
		func(_ error) { _ = httpServer.Shutdown(context.Background()) },
	)
	defer httpServer.Close()

	group.Add(
		func() error { return grpcServer.Serve(grpcLn) },
		func(_ error) { grpcServer.GracefulStop() },
	)
	defer grpcServer.Stop()

	// Setup signal handler
	group.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	err = group.Run()
	if err != nil {
		if e := (run.SignalError{}); errors.As(err, &e) {
			log.Println(err)

			return
		}

		// Fatal error
		// We don't use fatal, so deferred functions can do their jobs.
		log.Println(err)
	}
}
