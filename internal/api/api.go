package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitApi() *api {
	return &api{}
}

type api struct {
}

func (a *api) InitRouter() *gin.Engine {
	router := gin.Default()
	a.initHandlers(router)
	return router
}

func (a *api) initHandlers(r *gin.Engine) {
	r.LoadHTMLFiles("index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	api := r.Group("/api")
	{
		api.POST("/get-id")
	}
}
