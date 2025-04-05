package handler

import (
	"encoding/json"	
	"net/http"
	"github.com/gin-gonic/gin"
)

type RequestCheckin struct {
	id string `json:"id"`
	date   string `json:"date"`
	direction string `json:"direction"`
	returnTime   string `json:"returnTime"`
	status string `json:"status"`
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status":"Up"})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","name": "John", "email": "johndoe@gmail.com","userType":"mensalista"})
}

func Checkin(c *gin.Context) {
	var checkin RequestCheckin
	//c.JSON(http.StatusOK, gin.H{"id":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","date":"","direction":"","status":"REGISTERED"})
    if err := c.BindJSON(&ncheckin); err != nil {
        return
    }
    checkin.id = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    checkin.status = "REGISTERED"
	c.IndentedJSON(http.StatusCreated, checkin)
}

func ErrRouter(c *gin.Context) {
	c.String(http.StatusBadRequest, "url err")
}