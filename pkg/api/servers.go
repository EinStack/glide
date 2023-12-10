package api

import (
	"context"
	"sync"
)

type ServerManager struct {
	shutdownWG *sync.WaitGroup
}

func NewServerManager() (*ServerManager, error) {
	return &ServerManager{
		shutdownWG: &sync.WaitGroup{},
	}, nil
}

func (mgr ServerManager) Shutdown(ctx context.Context) error {
	return nil
}
