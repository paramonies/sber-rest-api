package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/paramonies/sber-rest-api/internal/app/model"
)

func (c *Controller) CreateItem(ctx *gin.Context) {
	var input model.Item

	if err := ctx.BindJSON(&input); err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, "invalid input item body")
		return
	}

	if input.UserId < 1 {
		SendErrorResponse(ctx, http.StatusBadRequest, "invalid user-id")
		return
	}
	newItem, err := c.service.CreateItem(input)
	if err != nil {
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, newItem)
}

func (c *Controller) UpdateItem(ctx *gin.Context) {
	id := ctx.Param("id")
	itemId, err := strconv.Atoi(id)
	if err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, "item id must be integer")
		return
	}
	var input model.UpdateItem

	if err := ctx.BindJSON(&input); err != nil {
		SendErrorResponse(ctx, http.StatusBadRequest, "invalid input item body")
		return
	}

	if *input.UserId < 1 {
		SendErrorResponse(ctx, http.StatusBadRequest, "invalid user-id")
		return
	}
	newItem, err := c.service.UpdateItem(itemId, input)
	if err != nil {
		SendErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, newItem)
}
