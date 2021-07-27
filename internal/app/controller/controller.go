package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paramonies/sber-rest-api/internal/app/model"
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

	router.POST("/create", c.CreateUser)
	router.PUT("/update", c.UpdateUser)
	router.GET("/get/:id", c.GetUser)
	router.DELETE("/delete/:id", c.DeleteUser)
	router.GET("/list", c.GetListUsers)

	return router
}

func (c *Controller) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, "user id must be integer")
		return
	}
	user, err := c.service.GetUserById(userId)
	if err != nil {
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) CreateUser(ctx *gin.Context) {
	var input model.User

	if err := ctx.BindJSON(&input); err != nil {
		fmt.Println(input)
		SendErrorResponse(ctx, http.StatusBadRequest, "invalid input body")
		return
	}

	newUser, err := c.service.CreateUser(input)
	if err != nil {
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, newUser)
}

func (c *Controller) UpdateUser(ctx *gin.Context) {
	var input model.UpdateUser

	if err := ctx.BindJSON(&input); err != nil {
		fmt.Println(input)
		SendErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	updatedUser, err := c.service.UpdateUser(input)
	if err != nil {
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(updatedUser)
	ctx.JSON(http.StatusOK, updatedUser)
}

func (c *Controller) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, "user id must be integer")
		return
	}
	err = c.service.DeleteUser(userId)
	if err != nil {
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, statusResponse{"OK"})
}

func (c *Controller) GetListUsers(ctx *gin.Context) {
	//..../list?page=2&limit=5
	pageStr := ctx.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limitStr := ctx.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 3
	}

	users, err := c.service.GetListUsers(page, limit)
	if err != nil {
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, users)
}
