package handler

import (
	"time"
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
	checkins := generateFakeCheckins()
    c.JSON(http.StatusOK, gin.H{"data": checkins})
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
    for i := 0; i < 3; i++ { // Gera 3 dias de checkins (hoje + 2 dias atrás)
        date := currentDate.AddDate(0, 0, -i).Format("2006-01-02")
        checkins = append(checkins, Checkin{uuid.New().String(), date, "ida", "", "REGISTERED"})
        checkins = append(checkins, Checkin{uuid.New().String(), date, "retorno", "17h10", "REGISTERED"})
    }
    return checkins
}