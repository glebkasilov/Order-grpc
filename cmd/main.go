package main

import (
	"context"
	"fmt"
	"lesson3/internal/config"
	"lesson3/internal/service"
	"net/http"

	"lesson3/pkg/api/test/api"
	"lesson3/pkg/database/postgres"
	"lesson3/pkg/logger"
	"log"
	"net"
	"os"
	"os/signal"

	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)

	server_cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server_cfg.Host, server_cfg.Port))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	storage, err := postgres.New()
	if err != nil {
		log.Fatalf("failed to create storage: %v", err)
	}
	srv := service.New(storage)

	recoveryOpts := []recovery.Option{recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		return err
	})}
	go runRest(ctx, server_cfg)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(recoveryOpts...), logging.UnaryServerInterceptor(logger.InterceptorLogger(logger.GetLoggerFromCtx(ctx)))))

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Server started", zap.String("port", fmt.Sprintf("%d", (server_cfg.Port))))

	go func() {
		select {
		case <-ctx.Done():
			server.GracefulStop()
			logger.GetLoggerFromCtx(ctx).Info(ctx, "Server stopped")
			<-ctx.Done()
		}
	}()

	api.RegisterOrderServiceServer(server, srv)

	if err := server.Serve(lis); err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "failed to serve: %w", zap.Error(err))
	}

}

func runRest(ctx context.Context, server_cfg *config.Config) {
	rt := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	grpcEndpoint := fmt.Sprintf("%s:%d", server_cfg.Host, server_cfg.Port)

	err := api.RegisterOrderServiceHandlerFromEndpoint(ctx, rt, grpcEndpoint, opts)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to register gateway", zap.Error(err))
		return
	}

	h2Server := &http2.Server{}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", server_cfg.HttpPort),
		Handler: h2c.NewHandler(rt, h2Server),
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Starting REST gateway", zap.String("port", server_cfg.HttpPort))
	if err := server.ListenAndServe(); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "REST server failed", zap.Error(err))
	}
}
