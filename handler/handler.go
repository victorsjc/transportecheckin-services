package handler

import (
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

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
	checkins := generateFakeCheckins()
    c.JSON(http.StatusOK, checkins)
}

func Login(c *gin.Context) {
	var req Login
    if err := c.BindJSON(&req); err != nil {
        return
    }    
	if(req.Password != nil && req.Password === "itau1234"){
		return c.JSON(http.StatusOK, gin.H{})
	}
	c.JSON(http.StatusUnauthorized)
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func ErrRouter(c *gin.Context) {
	c.String(http.StatusBadRequest, "url err")
}

// Função para criar checkins fake
func generateFakeCheckins() []Checkin {
    var checkins []Checkin
    currentDate := time.Now()
    date := currentDate.Format("2006-01-02")
    checkins = append(checkins, Checkin{uuid.New().String(), date, "ida", "", "REGISTERED"})
    return checkins
}