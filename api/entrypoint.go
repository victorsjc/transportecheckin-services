package api

import (
 "fmt"
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

//entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
    app = gin.New()
    route := app.Group("/")

    // register route
    registerRouter(route)
    app.ServeHTTP(w,r)
}