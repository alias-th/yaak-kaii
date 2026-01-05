package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"yaak-kaii/services/order-service/internal/database"
	grpcserver "yaak-kaii/services/order-service/internal/grpc_server"
	"yaak-kaii/services/order-service/internal/repositories"
	"yaak-kaii/services/order-service/internal/services"

	pb "yaak-kaii/shared/proto/order"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	GRPCServer *grpc.Server
	DB         *gorm.DB
}

func NewApp() *App {

	db := database.InitDB()

	orderRepo := repositories.NewOrderRepository(db)
	orderItemRepo := repositories.NewOrderItemRepository(db)

	orderService := services.NewOrderService(orderRepo)
	orderItemService := services.NewOrderItemService(orderItemRepo)

	server := grpc.NewServer()
	orderGrpc := grpcserver.NewOrderGRPCServer(orderService, orderItemService)

	pb.RegisterOrderServiceServer(server, orderGrpc)

	return &App{
		GRPCServer: server,
		DB:         db,
	}

}

func (a *App) Run() {
	port := ":9092"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("gRPC server is running on port %s", port)
	a.GRPCServer.Serve(lis)

}

func (a *App) Shutdown() {
	log.Println("Shutting down gracefully...")

	a.GRPCServer.GracefulStop()
	log.Println("gRPC server stopped.")

	sqlDB, err := a.DB.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}

func main() {
	app := NewApp()

	go func() {
		app.Run()
	}()

	// Graceful stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	app.Shutdown()
	log.Println("Application exited gracefully")

}
