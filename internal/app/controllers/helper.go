package controllers

import (
	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context, errStr string) {
	c.AbortWithStatusJSON(400, gin.H{"code": "BadRequest", "message": errStr})
}

func Unauthorized(c *gin.Context, errStr string) {
	c.AbortWithStatusJSON(401, gin.H{"code": "Unauthorized", "message": errStr})
}

func NotFound(c *gin.Context, errStr string) {
	c.AbortWithStatusJSON(404, gin.H{"code": "NotFound", "message": errStr})
}

func InternalServerError(c *gin.Context, errStr string) {
	c.AbortWithStatusJSON(500, gin.H{"code": "InternalServerError", "message": errStr})
}
func SuccessWithData(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"message": "success", "data": data})
}

func Success(c *gin.Context) {
	c.JSON(200, gin.H{"message": "success"})
}
