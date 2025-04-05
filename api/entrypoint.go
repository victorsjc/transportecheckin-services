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
    r.POST("/api/register", handler.Register)
}

func Init(){
    app = gin.New()
    r := app.Group("/")
    gin.Use(CORSMiddleware())
    // register route
    registerRouter(r)
}

//entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
    Init()
    app.ServeHTTP(w,r)
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "https://ui-transportecheckin-app.vercel.app")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusNoContent)
            return
        }
        c.Next()
    }
}