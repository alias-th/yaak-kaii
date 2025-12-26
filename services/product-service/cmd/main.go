package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"
	"yaak-kaii/services/product-service/internal/db/mongodb"
	"yaak-kaii/services/product-service/internal/grpc"
	"yaak-kaii/services/product-service/internal/repository"
	"yaak-kaii/services/product-service/internal/service"

	grpcserver "google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// MongoDB connection URI
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Database name
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "yaak_kaii_products"
	}

	// Connect to MongoDB
	mongoClient, err := mongodb.Connect(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	defer func() {
		disconnectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(disconnectCtx); err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		}
	}()
	mongodb := mongoClient.GetDatabase()

	// Initialize repositories
	mongoDBRepo := repository.NewMongoRepository(mongodb)
	svc := service.NewService(mongoDBRepo)

	// Initialize Grpc server
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, svc)

	log.Printf("Starting gRPC server Auth service on port %s", lis.Addr().String())
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
