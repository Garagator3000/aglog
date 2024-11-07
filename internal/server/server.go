package server

import "context"

type Server interface {
	GetAddr() string
	Start()
	Stop()
	Listen(ctx context.Context, buffer chan string)
}
