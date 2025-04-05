package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status":"Up"})
}

func ErrRouter(c *gin.Context) {
	c.String(http.StatusBadRequest, "url err")
}

func Cors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Next()
}