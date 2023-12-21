package pkg

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"glide/pkg/telemetry"
	"go.uber.org/zap"

	"glide/pkg/api"
	"glide/pkg/api/http"

	"go.uber.org/multierr"
)

// Gateway represents an instance of running Glide gateway.
// It loads configs, start API server(s), and listen to termination signals to shut down
type Gateway struct {
	// telemetry holds logger, meter, and tracer
	telemetry *telemetry.Telemetry
	// serverManager controls API over different protocols
	serverManager *api.ServerManager
	// signalChannel is used to receive termination signals from the OS.
	signalC chan os.Signal
	// shutdownC is used to terminate the gateway
	shutdownC chan struct{}
}

func NewGateway() (*Gateway, error) {
	// TODO: gonna be read from a config file
	logConfig := telemetry.NewLogConfig()
	logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logConfig.Encoding = "console"

	tel, err := telemetry.NewTelemetry(&telemetry.Config{LogConfig: logConfig})
	if err != nil {
		return nil, err
	}

	serverManager, err := api.NewServerManager(&http.ServerConfig{})
	if err != nil {
		return nil, err
	}

	return &Gateway{
		telemetry:     tel,
		serverManager: serverManager,
		signalC:       make(chan os.Signal, 3), // equal to number of signal types we expect to receive
		shutdownC:     make(chan struct{}),
	}, nil
}

// Run starts and runs the gateway according to given configuration
func (gw *Gateway) Run(ctx context.Context) error {
	// TODO: init configs
	gw.serverManager.Start()

	signal.Notify(gw.signalC, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(gw.signalC)

LOOP:
	for {
		select {
		// TODO: Watch for config updates
		case sig := <-gw.signalC:
			gw.telemetry.Logger().Info("Received signal from OS", zap.String("signal", sig.String()))
			break LOOP
		case <-gw.shutdownC:
			gw.telemetry.Logger().Info("received shutdown request")
			break LOOP
		case <-ctx.Done():
			gw.telemetry.Logger().Info("context done, terminating process")
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
