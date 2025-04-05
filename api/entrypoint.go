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

func Init(){
    app = gin.New()
    r := app.Group("/")

    // register route
    registerRouter(r)
    fmt.Println("API inicializada com sucesso")
}

//entrypoint
func Handler(w http.ResponseWriter, r *http.Request) {
    Init()
    app.ServeHTTP(w,r)
}