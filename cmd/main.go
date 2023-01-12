package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/enzofaliMELI/web-server/cmd/docs"
	"github.com/enzofaliMELI/web-server/cmd/handlers"
	"github.com/enzofaliMELI/web-server/cmd/routes"
	"github.com/enzofaliMELI/web-server/internal/domain"
	"github.com/enzofaliMELI/web-server/internal/product"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title MELI Bootcamp API
// @version 1.0
// @description This API Handle MELI Products.
// @termsOfService https://developers.mercadolibre.com.ar/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url https://developers.mercadolibre.com.ar/support

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
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

	// Swagger
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	server.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Router
	server.GET("/ping", handlers.Pong)
	routes := routes.NewRouter(server, &db)
	routes.SetRoutes()

	// Start Server
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
