package main

import (
	"log"

	"github.com/devloeprsabbir/go-elasticsearch/config"
	"github.com/devloeprsabbir/go-elasticsearch/elastic"
	handler "github.com/devloeprsabbir/go-elasticsearch/handler"
	"github.com/devloeprsabbir/go-elasticsearch/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadEnv()

	config.ConnectDatabase(cfg)

	r := gin.New()
	r.Use(gin.Logger())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	todoRepo := repository.NewTodoRepository(config.DB)
	todoHandler := handler.NewTodoHandler(todoRepo, elastic.EsClient)
	v1 := r.Group("/api/v1")
	{
		todos := v1.Group("/todos")
		{
			todos.GET("", todoHandler.GetTodo)
			todos.POST("", todoHandler.CreateTodo)
		}
	}

	// v1 := r.Group("/api/v1")
	// {
	// 	todos := v1.Group("/todos")
	// 	{
	// 		todos.POST("", todoHandler.CreateTodo)
	// 		todos.GET("", todoHandler.ListTodos)
	// 		todos.GET("/search", todoHandler.SearchTodos) // ES-powered search
	// 		todos.POST("/seed", todoHandler.SeedTodos)    // bulk seed for testing
	// 		todos.GET("/:id", todoHandler.GetTodo)
	// 		todos.PUT("/:id", todoHandler.UpdateTodo)
	// 		todos.DELETE("/:id", todoHandler.DeleteTodo)
	// 	}
	// }

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
