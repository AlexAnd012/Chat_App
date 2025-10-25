package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *Handlers) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger(), gin.Recovery())

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
