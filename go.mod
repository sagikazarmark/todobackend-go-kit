module github.com/sagikazarmark/todobackend-go-kit

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/99designs/gqlgen v0.13.0
	github.com/go-bdd/gobdd v1.1.3
	github.com/go-kit/kit v0.12.0
	github.com/golang/protobuf v1.5.2
	github.com/goph/idgen v0.4.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/oklog/run v1.1.0
	github.com/sagikazarmark/appkit v0.12.0
	github.com/sagikazarmark/kitx v0.15.0
	github.com/sagikazarmark/todobackend-go-kit/api v0.4.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/vektah/gqlparser/v2 v2.1.0
	google.golang.org/grpc v1.42.0
)

replace github.com/sagikazarmark/todobackend-go-kit/api => ./api
