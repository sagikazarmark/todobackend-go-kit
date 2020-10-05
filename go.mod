module github.com/sagikazarmark/todobackend-go-kit

go 1.14

require (
	emperror.dev/errors v0.7.0
	github.com/99designs/gqlgen v0.11.3
	github.com/go-bdd/gobdd v1.1.0
	github.com/go-kit/kit v0.10.0
	github.com/golang/protobuf v1.4.2
	github.com/goph/idgen v0.4.0
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/markbates/pkger v0.17.0
	github.com/oklog/run v1.1.0
	github.com/sagikazarmark/appkit v0.10.0
	github.com/sagikazarmark/kitx v0.14.0
	github.com/sagikazarmark/todobackend-go-kit/api v0.3.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.1
	github.com/vektah/gqlparser/v2 v2.0.1
	google.golang.org/grpc v1.32.0
)

replace github.com/sagikazarmark/todobackend-go-kit/api => ./api
