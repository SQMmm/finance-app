package server

import (
	"context"
	"sync"
)

type Server interface {
	Serve(ctx context.Context, wg *sync.WaitGroup) error
}
