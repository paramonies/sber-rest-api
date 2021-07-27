package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/paramonies/sber-rest-api/internal/app/service"
)

type Controller struct {
	service *service.Service
}

func NewController(services *service.Service) *Controller {
	return &Controller{service: services}
}

func (c *Controller) InitRoutes() *gin.Engine {
	router := gin.Default()

	user := router.Group("/user")
	{
		user.POST("/create", c.CreateUser)
		user.PUT("/update", c.UpdateUser)
		user.GET("/get/:id", c.GetUser)
		user.DELETE("/delete/:id", c.DeleteUser)
		user.GET("/list", c.GetListUsers)
	}

	item := router.Group("/item")
	{
		item.POST("/create", c.CreateItem)
		item.PUT("/update/:id", c.UpdateItem)
	}

	return router
}
