package presentation

import "context"

// Server is an interface for HTTP server.
type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
