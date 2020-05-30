module github.com/sagikazarmark/todobackend-go-kit

go 1.14

require (
	emperror.dev/errors v0.7.0
	github.com/go-bdd/gobdd v1.0.1
	github.com/go-kit/kit v0.10.0
	github.com/goph/idgen v0.4.0
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/sagikazarmark/appkit v0.9.0
	github.com/sagikazarmark/kitx v0.12.0
	github.com/sagikazarmark/todobackend-go-kit/api v0.0.0-00010101000000-000000000000
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.5.1
)

replace github.com/sagikazarmark/todobackend-go-kit/api => ./api
