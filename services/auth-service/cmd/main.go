package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"yaak-kaii/services/auth-service/internal/auth"
	"yaak-kaii/services/auth-service/internal/config"
	db "yaak-kaii/services/auth-service/internal/db/sqlc"
	"yaak-kaii/services/auth-service/internal/grpc"
	"yaak-kaii/services/auth-service/internal/repository"
	"yaak-kaii/services/auth-service/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	grpcserver "google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load Config
	appConfig := config.LoadConfig()

	// Initialize PG connection pool
	connPool, err := pgxpool.New(ctx, appConfig.ConnString)
	if err != nil {
		log.Fatal("cannot connect to db")
	}
	store := db.NewStore(connPool)

	// Initialize repositories
	userRepo := repository.NewUserRepository(store)
	guestRepo := repository.NewGuestRepository(store)
	tokenRepo := repository.NewRefreshTokenRepository(store)
	roleRepo := repository.NewRoleRepository(store)

	// Initialize jwt
	jwtAuth := auth.NewJWTAuthenticator(appConfig.JwtSecret, appConfig.JwtISS, appConfig.JwtISS)

	// Initialize service
	svc := service.NewService(userRepo, guestRepo, tokenRepo, roleRepo, jwtAuth)
	log.Println("connected to database", appConfig.ConnString)

	// Graceful shutdown on interrupt signals
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	lis, err := net.Listen("tcp", appConfig.GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpcserver.NewServer()

	grpc.NewGRPCHandler(grpcServer, svc)

	log.Printf("Starting gRPC server Trip service on port %s", lis.Addr().String())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
