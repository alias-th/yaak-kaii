package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yaak-kaii/services/api-gateway/internal/auth"
	grpcclients "yaak-kaii/services/api-gateway/internal/grpc_clients"
	"yaak-kaii/services/api-gateway/internal/s3"

	"github.com/gin-gonic/gin"
)

type Application struct {
	Config      Config
	GrpcClients *grpcclients.GrpcClients
	JwtAuth     *auth.JWTAuthenticator
	S3Uploader  *s3.S3Uploader
}

type Config struct {
	Addr string
}

func (app *Application) Run() {
	router := gin.Default()

	{
		v1 := router.Group("/v1")
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		auth := v1.Group("/authentication")
		{
			auth.POST("/users", app.createUser)
			auth.POST("/guest", app.createGuest)
			auth.POST("/login", app.login)
			auth.POST("/rotate-token", app.rotateToken)
		}

		seller := v1.Group("/seller").Use(app.sellerMiddleware())
		{
			seller.GET("/:shopId/products/:slug", app.getProductDetails)
			seller.POST("/products", app.createProduct)
			seller.POST("/products/:id/images", app.uploadProductImages)
			seller.POST("/products/:id/variant", app.createProductVariant)
			seller.POST("/categories", app.createCategory)
		}

		buyer := v1.Group("/buyer")
		{
			buyer.GET("/products", app.listProducts)
		}

	}

	srv := &http.Server{
		Addr:    app.Config.Addr,
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		log.Printf("Server listening on %s", app.Config.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
