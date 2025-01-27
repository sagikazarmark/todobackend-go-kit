module github.com/sagikazarmark/todobackend-go-kit

go 1.23.1

require (
	emperror.dev/errors v0.8.1
	github.com/99designs/gqlgen v0.17.63
	github.com/go-bdd/gobdd v1.1.3
	github.com/go-chi/chi/v5 v5.2.0
	github.com/go-chi/cors v1.2.1
	github.com/go-kit/kit v0.13.0
	github.com/golang/protobuf v1.5.4
	github.com/goph/idgen v0.4.0
	github.com/oklog/run v1.1.0
	github.com/sagikazarmark/appkit v0.16.0
	github.com/sagikazarmark/kitx v0.20.0
	github.com/sagikazarmark/todobackend-go-kit/api v0.7.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.10.0
	github.com/vektah/gqlparser/v2 v2.5.21
	google.golang.org/grpc v1.70.0
)

require (
	github.com/agnivade/levenshtein v1.2.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/cucumber/gherkin-go/v13 v13.0.0 // indirect
	github.com/cucumber/messages-go/v12 v12.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/moogar0880/problems v0.1.1 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/urfave/cli/v2 v2.27.5 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241230172942-26aa7a208def // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/sagikazarmark/todobackend-go-kit/api => ./api

// required due to recent submodule changes
replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240930140551-af27646dc61f
