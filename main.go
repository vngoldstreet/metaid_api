package main

import (
	"time"
	"vietvd/mql-api/controllers"
	"vietvd/mql-api/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	repository.ConnectToMongoDB()
	repository.SetupLogger()
}

func main() {
	r := gin.Default()
	r.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"POST", "GET", "PUT"},
			AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, ResponseType, accept, origin, Cache-Control, X-Requested-With"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		repository.Logger(),
	)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	api := r.Group("/api/v1")
	{
		api.GET("/get-new-mtid", controllers.HandleGetNewMetaID)
		api.GET("/get-mtid", controllers.HandleGetMetaID)
		api.PUT("/delete-mtid", controllers.HandleDeleteMetaID)
	}
	r.Run(":50100")
}
