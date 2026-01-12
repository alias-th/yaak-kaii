package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"yaak-kaii/services/product-service/internal/database"
	grpcserver "yaak-kaii/services/product-service/internal/grpc_server"
	"yaak-kaii/services/product-service/internal/repositories"
	"yaak-kaii/services/product-service/internal/services"

	pb "yaak-kaii/shared/proto/product"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	GRPCServer *grpc.Server
	DB         *gorm.DB
}

func NewApp() *App {

	db := database.InitDB()

	// Run migrations for models
	database.Migrate(db)

	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)

	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	server := grpc.NewServer()
	productGrpc := grpcserver.NewProductGRPCServer(productService, categoryService)

	pb.RegisterProductServiceServer(server, productGrpc)

	return &App{
		GRPCServer: server,
		DB:         db,
	}

}

func (a *App) Run() {
	port := ":9091"
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
