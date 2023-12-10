package pkg

import (
	"context"
	"fmt"
	"glide/pkg/api"
	"go.uber.org/multierr"
	"os"
	"os/signal"
	"syscall"
)

// Gateway represents an instance of running Glide gateway.
// It loads configs, start API server(s), and listen to termination signals to shut down
type Gateway struct {
	// serverManager controls API over different protocols
	serverManager *api.ServerManager
	// signalChannel is used to receive termination signals from the OS.
	signalC chan os.Signal
	// shutdownC is used to terminate the gateway
	shutdownC chan struct{}
}

func NewGateway() (*Gateway, error) {
	serverManager, err := api.NewServerManager()

	if err != nil {
		return nil, err
	}

	return &Gateway{
		serverManager: serverManager,
		signalC:       make(chan os.Signal, 2), // equal to number of signal types we expect to receive
		shutdownC:     make(chan struct{}),
	}, nil
}

// Run starts and runs the gateway according to given configuration
func (gw *Gateway) Run(ctx context.Context) error {
	// TODO: init server manager
	// TODO: init configs
	signal.Notify(gw.signalC, os.Interrupt, syscall.SIGTERM)

LOOP:
	for {
		select {
		// TODO: Watch for config updates
		case <-gw.signalC:
			// TODO: log this occurrence
			break LOOP
		case <-gw.shutdownC:
			// TODO: log this occurrence
			break LOOP
		case <-ctx.Done():
			// TODO: log this occurrence
			// Call shutdown with background context as the passed in context has been canceled
			return gw.shutdown(context.Background())
		}
	}

	return gw.shutdown(ctx)
}

func (gw *Gateway) Shutdown() {
	close(gw.shutdownC)
}

func (gw *Gateway) shutdown(ctx context.Context) error {
	var errs error

	if err := gw.serverManager.Shutdown(ctx); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to shutdown servers: %w", err))
	}

	return errs
}
