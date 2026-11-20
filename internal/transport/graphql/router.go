package graph

import (
	"github.com/binariado/go-service-template/internal/application"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes mounts the GraphQL endpoint at /query and, when
// playgroundEnabled is true, the playground at /.
//
// GraphQL is an optional adapter: to drop it, delete internal/transport/graphql,
// gqlgen.yml, the RegisterRoutes call in cmd/api/main.go and the gqlgen
// dependency (go mod tidy).
func RegisterRoutes(r chi.Router, app *application.Application, playgroundEnabled bool) {
	srv := gqlhandler.NewDefaultServer(NewExecutableSchema(Config{
		Resolvers: &Resolver{App: app},
	}))

	if playgroundEnabled {
		r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}
	r.Handle("/query", srv)
}
