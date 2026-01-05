package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"yaak-kaii/services/product-service/internal/db/mongodb"
	grpc_server "yaak-kaii/services/product-service/internal/grpc"
	"yaak-kaii/services/product-service/internal/repository"
	"yaak-kaii/services/product-service/internal/service"

	"google.golang.org/grpc"
	grpcserver "google.golang.org/grpc"
)

type App struct {
	GRPCServer  *grpc.Server
	MongoClient *mongodb.MongoClient
}

func NewApp(ctx context.Context) *App {

	// Connect to MongoDB
	mongoUsername := os.Getenv("MONGODB_USERNAME")
	mongoPASSWORD := os.Getenv("MONGODB_PASSWORD")
	mongoURI := fmt.Sprintf("mongodb://%s:%s@mongo-headless:27017", mongoUsername, mongoPASSWORD)
	log.Println(mongoURI, mongoUsername)
	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "yaak_kaii_products"
	}
	mongoClient, err := mongodb.Connect(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}

	mongodb := mongoClient.GetDatabase()

	// Initialize repositories
	mongoDBRepo := repository.NewMongoRepository(mongodb)
	svc := service.NewService(mongoDBRepo)

	grpcServer := grpcserver.NewServer()
	grpc_server.NewGRPCHandler(grpcServer, svc)

	return &App{
		GRPCServer:  grpcServer,
		MongoClient: mongoClient,
	}

}

func (a *App) Run() {
	// Initialize Grpc server
	port := ":9091"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("gRPC server is running on port %s", port)
	a.GRPCServer.Serve(lis)
}

func (a *App) Shutdown(ctx context.Context) {
	log.Println("Shutting down gracefully...")

	a.GRPCServer.GracefulStop()
	log.Println("gRPC server stopped.")

	err := a.MongoClient.Disconnect(ctx)
	if err != nil {
		log.Println("failed to disconnect mongo.", err)

	}

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := NewApp(ctx)

	go func() {
		app.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	app.Shutdown(ctx)
	log.Println("Application exited gracefully")
}
