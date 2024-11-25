module github.com/sagikazarmark/todobackend-go-kit

go 1.23.1

require (
	emperror.dev/errors v0.8.1
	github.com/99designs/gqlgen v0.17.32
	github.com/go-bdd/gobdd v1.1.3
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-chi/cors v1.2.1
	github.com/go-kit/kit v0.13.0
	github.com/golang/protobuf v1.5.4
	github.com/goph/idgen v0.4.0
	github.com/oklog/run v1.1.0
	github.com/sagikazarmark/appkit v0.15.0
	github.com/sagikazarmark/kitx v0.19.0
	github.com/sagikazarmark/todobackend-go-kit/api v0.7.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.10.0
	github.com/vektah/gqlparser/v2 v2.5.2-0.20230422221642-25e09f9d292d
	google.golang.org/grpc v1.67.1
)

require (
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/cucumber/gherkin-go/v13 v13.0.0 // indirect
	github.com/cucumber/messages-go/v12 v12.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.2 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/moogar0880/problems v0.1.1 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.25.1 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/oauth2 v0.23.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240930140551-af27646dc61f // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/sagikazarmark/todobackend-go-kit/api => ./api

// required due to recent submodule changes
replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240930140551-af27646dc61f
