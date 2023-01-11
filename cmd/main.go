package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/enzofaliMELI/web-server/cmd/handlers"
	"github.com/enzofaliMELI/web-server/cmd/routes"
	"github.com/enzofaliMELI/web-server/internal/domain"
	"github.com/enzofaliMELI/web-server/internal/product"
)

func main() {

	// Get .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error in loading .env")
	}

	// Read all files
	db := []domain.Product{}
	product.OpenProducts(&db)

	// Server
	server := gin.Default()

	// Router
	server.GET("/ping", handlers.Pong)
	routes := routes.NewRouter(server, &db)
	routes.SetRoutes()

	// Start Server
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
