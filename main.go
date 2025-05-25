package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nikkmidl/rig-api/adapters"
	"github.com/nikkmidl/rig-api/internal/app"
	"github.com/nikkmidl/rig-api/internal/handler"
	"github.com/nikkmidl/rig-api/pkg/config"
	"github.com/nikkmidl/rig-api/pkg/opa"
	proto "github.com/nikkmidl/rig-api/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	if config.Config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Initialize context with signal handling
	// This will allow graceful shutdown on interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	// Initialize github client and OPA evaluator
	// Ensure that the GitHub token is set in the environment
	gh := adapters.NewGHClient(config.Config.GithubToken)
	opaEval, err := opa.NewEvaluator(ctx)

	if err != nil {
		log.Fatal().Err(err).Msg("evaluator initialization failed")
	}

	svc := app.New(gh, opaEval)
	handler := handler.NewHandler(svc)

	grpcServer := initGRPCServer(handler, stop)
	initHTTPServer(handler, ctx, stop)

	// Wait for interrupt signal
	<-ctx.Done()
	log.Info().Msg("Shutting down servers")
	grpcServer.GracefulStop()
}

func initGRPCServer(handler *handler.Handler, stop context.CancelFunc) *grpc.Server {
	grpcServer := grpc.NewServer()
	proto.RegisterAccessServiceServer(grpcServer, handler)

	listener, err := net.Listen("tcp", config.Config.GRPCServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
		stop()
	}

	// Start gRPC server in a goroutine
	go func() {
		log.Info().Msgf("gRPC server running on %s", config.Config.GRPCServerAddress)
		if err := grpcServer.Serve(listener); err != nil {
			log.Error().Err(err).Msg("gRPC server failed to serve")
			stop()
		}
	}()
	return grpcServer
}

func initHTTPServer(handler *handler.Handler, ctx context.Context, stop context.CancelFunc) {
	// Start HTTP gateway server in a goroutine
	go func() {
		grpcMux := runtime.NewServeMux()

		if err := proto.RegisterAccessServiceHandlerServer(ctx, grpcMux, handler); err != nil {
			log.Error().Err(err).Msg("failed to register HTTP gateway handler")
			stop()
		}

		log.Info().Msg("HTTP Gateway running on " + config.Config.HTTPServerAddress)

		// CORS middleware can be added here if needed
		if err := http.ListenAndServe(config.Config.HTTPServerAddress, grpcMux); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP gateway server failed")
			stop()
		}
	}()
}
