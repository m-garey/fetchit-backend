package application

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/m-garey/fetchit-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/m-garey/fetchit-backend/internal/handler"
	"github.com/m-garey/fetchit-backend/internal/repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const port string = "8080"

func Run() {
	db := setupDB()
	defer db.Close(context.Background())

	repo := repository.New(db)
	h := handler.New(repo)
	router := setupRouter()
	setupHandler(router, h)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in goroutine for graceful shutdown
	go func() {
		log.Printf("HTTP server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	// Gracefully shutdown with 10s timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("server stopped gracefully")
}

func setupDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	return conn
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	return r
}

func setupHandler(r *gin.Engine, h handler.API) {
	api := r.Group("/api")
	{
		api.POST("/users", h.CreateUser)
		api.POST("/stores", h.CreateStore)
		api.POST("/purchase", h.RecordPurchase)
		api.GET("/stickers/:user_id", h.GetStickersByUser)
		api.GET("/stickers/:user_id/:store_id", h.GetSticker)
	}
}
