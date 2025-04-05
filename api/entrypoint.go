package api

import (
 "net/http"
 "github.com/gin-gonic/gin"
 "transportecheckin/handler"
)

var (
    app *gin.Engine
)

func registerRouter(r *gin.RouterGroup) {
    r.GET("/api/health", handler.Health)
    r.POST("/api/logout", handler.Logout)
    r.OPTIONS("/api/logout", func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "https://ui-transportecheckin-app.vercel.app")
        c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        c.AbortWithStatus(200)
    })    
    r.POST("/api/register", handler.Register)
    r.OPTIONS("/api/register", func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "https://ui-transportecheckin-app.vercel.app")
        c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        c.AbortWithStatus(200)
    })
    r.GET("/api/checkins", handler.GetAllCheckins)
    r.POST("/api/checkins", handler.RegisterCheckin)
    r.OPTIONS("/api/checkins", func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "https://ui-transportecheckin-app.vercel.app")
        c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
        c.AbortWithStatus(200)
    })
}

func Init(){
    app = gin.New()
    r := app.Group("/")    
    // Middleware para lidar com CORS manualmente
    r.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "https://ui-transportecheckin-app.vercel.app")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
        }
        c.Next()
    })    
    // register route
    registerRouter(r)
}

//entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
    Init()
    app.ServeHTTP(w,r)
}