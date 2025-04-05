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
    r.GET("/api/ping", handler.Ping)
}

func Init(){
    app = gin.New()
    r := app.Group("/")

    // register route
    registerRouter(r)
}

//entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
    app.ServeHTTP(w,r)
}