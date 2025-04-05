package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Checkin struct {
	Id string `json:"id"`
	Date   string `json:"date"`
	Direction string `json:"direction"`
	ReturnTime   string `json:"returnTime"`
	Status string `json:"status"`
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status":"Up"})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","name": "John", "email": "johndoe@gmail.com","userType":"mensalista"})
}

func RegisterCheckin(c *gin.Context) {
	var checkin Checkin
    if err := c.BindJSON(&checkin); err != nil {
        return
    }
    checkin.Id = uuid.New().String()
    checkin.Status = "REGISTERED"
	c.IndentedJSON(http.StatusOK, checkin)
}

func GetAllCheckins(c *gin.Context) {
 	result := []Checkin{
 		{ uuid.New().String(), "2025-04-01", "ida", "", "REGISTERED"},
 		{ uuid.New().String(), "2025-04-01", "retorno", "17h10", "REGISTERED"},
 	}
	c.JSON(http.StatusOK, gin.H{result})
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"id":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855","name": "John", "email": "johndoe@gmail.com","userType":"mensalista"})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func ErrRouter(c *gin.Context) {
	c.String(http.StatusBadRequest, "url err")
}