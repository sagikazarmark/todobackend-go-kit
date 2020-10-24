bin/gqlgen: go.mod
	@mkdir -p bin
	go build -o bin/gqlgen github.com/99designs/gqlgen

.PHONY: graphql
graphql: bin/gqlgen ## Generate GraphQL code
	bin/gqlgen
