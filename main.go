package main

import (
    "fmt"
    "log"
    "net/http"
)


func handleFunc(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintf(w, "Hello!")

}

func main() {
    http.HandleFunc("/api", handleFunc)
    log.Println("Server is running:", 3001);
    log.Fatal(http.ListenAndServe(":3001", nil))
}