package graph

import "backend-service/internal/application"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
//
// The Resolver receives the Application aggregate so GraphQL resolvers never
// depend directly on domain services or infrastructure — only on the application layer.

type Resolver struct {
	App *application.Application
}
