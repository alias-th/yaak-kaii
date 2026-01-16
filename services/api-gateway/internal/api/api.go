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
			auth.POST("/login", app.login)
			auth.POST("/rotate-token", app.rotateToken)
		}

		admin := v1.Group("/admin").Use(app.adminMiddleware())
		{
			admin.GET("/products", app.listSellerProducts)
			admin.GET("/:shopId/products/:slug", app.getProductDetails)
			admin.POST("/products", app.createProduct)
			admin.POST("/products/:id/images", app.uploadProductImages)
			admin.POST("/products/:id/variants", app.createProductVariants)

			admin.GET("/categories", app.getCategories)
			admin.GET("/categories/:id", app.getCategory)
			admin.POST("/categories", app.createCategories)
		}

		buyer := v1.Group("/buyer")
		{
			buyer.GET("/products", app.listProducts)
			buyer.GET("/:shopId/products/:slug", app.getProductDetails)

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
