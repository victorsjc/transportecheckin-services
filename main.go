package main

import "net/http"
import "transportecheckin/api"

func main() {
    http.HandleFunc("/", api.Handler)
    http.ListenAndServe(":5000", nil)
}
