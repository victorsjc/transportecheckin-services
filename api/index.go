package api

import (
 "encoding/json"
 "net/http"

 "github.com/gin-contrib/cors"
 "github.com/gin-gonic/gin" 
)

const KN_SECURITY_HOLDER = "security_holder"
const KN_AUTHORIZATION = "Authorization"

func HealthIndicatorHandler(c *gin.Context) {
    //talvez o ideal seria verificar o status do banco de dados
    c.JSON(http.StatusOK, gin.H{
        "status": "Up",
    })
}

func Serve(){
    r := gin.Default()
    // Configuração do CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Middleware para verificar o token
    r.Use(func(c *gin.Context) {
        token := c.GetHeader(KN_AUTHORIZATION)
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token não fornecido"})
            c.Abort()
            return
        }
        c.Next()
    })


    r.GET("/health", HealthIndicatorHandler)

    ginLambda = ginadapter.New(r)
    r.Run(":8080")    
}