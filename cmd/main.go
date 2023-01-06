package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/enzofaliMELI/web-server/cmd/handlers"
	"github.com/enzofaliMELI/web-server/internal/product/repository"
)

const filename = "../products.json"

func main() {
	// Read all files
	repository.OpenProducts(filename)

	// server
	server := gin.Default()

	// router
	server.GET("/ping", handlers.Pong)
	products := server.Group("/products")
	products.POST("/", handlers.POSTProducts)
	products.GET("/", handlers.GETProducts)
	products.GET("/:id", handlers.ProductsId)
	products.GET("/search", handlers.ProductsSearch)

	// start
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
