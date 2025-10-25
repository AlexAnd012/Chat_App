package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "chatapp/docs" // путь к сгенерированным докам
)

func NewRouter(handler *Handlers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Swagger UI на /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/healthz", func(c *gin.Context) {
		if err := handler.repo.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"ok": false, "redis": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "ts": time.Now().UTC()})
	})

	api := r.Group("/api")
	{
		api.GET("/messages", handler.GetMessages)
		api.POST("/messages", handler.PostMessage)
	}
	return r
}
