package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/enzofaliMELI/web-server/cmd/handlers"
	"github.com/enzofaliMELI/web-server/pkg/response"
	"github.com/gin-gonic/gin"
)

func FirstMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("First Middleware")
		ctx.Next()
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("token")
		if token != os.Getenv("TOKEN") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.Err(handlers.ErrUnauthorized))
			return
		}
		ctx.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		url := ctx.Request.URL
		time := time.Now()
		size := ctx.Request.ContentLength
		fmt.Printf("Method: %v \nURL: %v \nTime: %v \nSize: %v\n", method, url, time, size)

		ctx.Next()
	}
}

func Middlewares(f gin.HandlerFunc) []gin.HandlerFunc {
	list := []gin.HandlerFunc{
		FirstMiddleware(),
		TokenAuthMiddleware(),
		Logger(),
	}
	list = append(list, f)
	return list
}
