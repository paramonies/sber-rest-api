package controller

import "github.com/gin-gonic/gin"

type errorMessage struct {
	Message string `json:"error"`
}

func SendErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, errorMessage{message})
}
