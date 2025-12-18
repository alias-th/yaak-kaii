package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"yaak-kaii/shared/contracts"

	"github.com/gin-gonic/gin"
)

func RunGinServer(addr string) {
	router := setupRouter()
	srv := &http.Server{
		Addr:    addr,
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		log.Printf("Server listening on %s", addr)
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

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/users", HandleCreateUser)
	// router.POST("/users/login", server.loginUser)
	// router.POST("/tokens/renew_access", server.renewAccessToken)
	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// authRoutes.POST("/accounts", server.createAccount)
	// authRoutes.GET("/accounts/:id", server.getAccount)
	// authRoutes.GET("/accounts", server.listAccounts)
	// authRoutes.POST("/transfers", server.createTransfer)

	return router
}

func errorResponse(err error, code int) gin.H {
	return gin.H{
		"error": contracts.APIError{
			Code:    code,
			Message: err.Error(),
		},
	}
}
