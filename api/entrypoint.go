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

    // register route
    registerRouter(r)
}

//entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
    Init()
    app.ServeHTTP(w,r)
}