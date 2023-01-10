package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/enzofaliMELI/web-server/cmd/handlers"
	"github.com/enzofaliMELI/web-server/internal/product"
)

const filename = "../products.json"

func main() {
	// Read all files
	product.OpenProducts(filename)

	// Server
	server := gin.Default()

	// Router
	server.GET("/ping", handlers.Pong)
	products := server.Group("/products")
	products.POST("/", handlers.StoreProduct)
	products.GET("/", handlers.GetAllProducts)
	products.GET("/:id", handlers.GetProductsId)
	products.GET("/search", handlers.GetProductsSearch)
	//products.PUT("/:id", handlers.UpdateProduct)
	//products.PATCH("/:id", handlers.UpdateProductName)
	//products.DELETE("/:id", handlers.DeleteProduct)

	// Start Server
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
