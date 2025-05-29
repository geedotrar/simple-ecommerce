package main

import (
	"log"
	"product-service/config"
	"product-service/internal/handlers"
	"product-service/internal/repository"
	"product-service/internal/routes"
	"product-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	server()
}

func server() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	redisClient := config.InitRedis()
	if redisClient == nil {
		log.Fatal("Failed to connect to Redis")
	}

	g := gin.Default()
	g.Use(gin.Recovery())

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	gormConfig := config.NewGormPostgres()
	if gormConfig == nil {
		log.Fatal("Failed to initialize database connection")
	}

	db := gormConfig.GetConnection()
	if db == nil {
		log.Fatal("Database connection is nil")
	}

	productGroup := g.Group("/products")
	productRepo := repository.NewProductRepository(gormConfig)
	productSvc := service.NewProductService(productRepo)
	uploadDir := "./uploads/products"
	productHdl := handlers.NewproductHandler(productSvc, uploadDir)
	productRouter := routes.NewProductRouter(productGroup, productHdl)
	productRouter.Mount()
	g.Static("/uploads", "./uploads")

	g.Run(":8080")
}
