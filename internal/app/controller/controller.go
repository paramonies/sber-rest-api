package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/paramonies/sber-rest-api/internal/app/service"
)

type Controller struct {
	services *service.Service
}

func NewController(services *service.Service) *Controller {
	return &Controller{services: services}
}

func (c *Controller) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/create", c.CreateUser)
	router.GET("/update", c.UpdateUser)

	return router
}

func (c *Controller) CreateUser(ctx *gin.Context) {
	SendErrorResponse(ctx, 200, "Create User")
	return
}

func (c *Controller) UpdateUser(ctx *gin.Context) {
	SendErrorResponse(ctx, 200, "Update User")
	return
}
