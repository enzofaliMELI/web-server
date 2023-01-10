package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/enzofaliMELI/web-server/cmd/handlers"
	"github.com/enzofaliMELI/web-server/cmd/routes"
	"github.com/enzofaliMELI/web-server/internal/domain"
	"github.com/enzofaliMELI/web-server/internal/product"
)

func main() {
	// Read all files
	db := []domain.Product{}
	product.OpenProducts(&db)

	//repository := product.NewRepository(&db)
	//service := product.NewService(repository)
	//handler := handler.NewProduct(service)

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
