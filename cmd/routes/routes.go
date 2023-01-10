package routes

import (
	"github.com/enzofaliMELI/web-server/cmd/handlers"
	"github.com/enzofaliMELI/web-server/internal/domain"
	"github.com/enzofaliMELI/web-server/internal/product"
	"github.com/gin-gonic/gin"
)

type Router struct {
	db *[]domain.Product
	en *gin.Engine
}

func NewRouter(en *gin.Engine, db *[]domain.Product) *Router {
	return &Router{en: en, db: db}
}

func (r *Router) SetRoutes() {
	r.SetProduct()
}

func (r *Router) SetProduct() {
	// instances
	repository := product.NewRepository(r.db)
	service := product.NewService(repository)
	handler := handlers.NewProduct(service)

	products := r.en.Group("/products")
	products.POST("/", handler.Store())
	products.GET("/", handler.GetAll())
	products.GET("/:id", handler.GetById())
	products.GET("/search", handler.GetPriceGt())
	products.PUT("/:id", handler.UpdateProduct())
	products.PATCH("/:id", handler.UpdateProductName())
	products.DELETE("/:id", handler.DeleteProduct())
}
