package router

import (
	e "errors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"goportscan/core/models"
	v1 "goportscan/core/router/v1"
	_ "goportscan/core/swagger"
	"net/http"
)

func NewHttpRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.NoRoute(func(c *gin.Context) {
		resp := models.Gin{Ctx: c}
		resp.Fail(http.StatusNotFound, e.New("not found route"))
	})
	router.NoMethod(func(c *gin.Context) {
		resp := models.Gin{Ctx: c}
		resp.Fail(http.StatusNotFound, e.New("not found method"))
	})
	router.GET("/ping", func(c *gin.Context) {
		resp := models.Gin{Ctx: c}
		resp.Success("pong")
	})
	v1.Register(router.Group("/api/v1")) //register api/v1 for producer
	return router
}
