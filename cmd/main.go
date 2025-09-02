package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"tournaments-core/internal/config"
	_grpc "tournaments-core/internal/delivery/grpc"
	"tournaments-core/internal/domain/ports/repository"
	"tournaments-core/internal/repository/postgresql"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
)

func main() {
	cfg := config.MustLoad()

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Name,
		cfg.DatabaseConfig.SslMode,
	)
	log.Printf("[POSTGRES]: Successful connection to: %s\n", dbUrl)

	gamesRepository, err := postgresql.NewGamesRepository(dbUrl)
	if err != nil {
		log.Fatalf("[POSTGRES]: Error while initializing repository: %w", err)
	}

	resultRepository, err := postgresql.NewResultsRepository(dbUrl)
	if err != nil {
		log.Fatalf("[POSTGRES]: Error while initializing repository: %w", err)
	}

	// TODO: logger

	go RunGrpcServer(cfg, &gamesRepository, &resultRepository)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		log.Fatal(fmt.Sprintf("signal.Notify: %v", v))
	case done := <-ctx.Done():
		log.Fatal(fmt.Sprintf("ctx.Done: %v", done))
	}
}

func RunGrpcServer(config *config.Config, games_rep *repository.GamesRepository, res_rep *repository.ResultsRepository) {
	grpcServer := grpc.NewServer()
	_grpc.NewGamesGrpcServer(grpcServer, games_rep)
	_grpc.NewResultsGrpcServer(grpcServer, res_rep)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", config.GrpcConfig.Port)
	if err != nil {
		log.Fatalf("[gRPC]: Failed to listen: %v", err)
	}

	go func() {
		log.Println(fmt.Sprintf("[gRPC]: Grpc server listens to: %s", config.GrpcConfig.Port))
		log.Fatal(grpcServer.Serve(lis))
	}()
}
