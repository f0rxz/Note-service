package main

import (
	"log"
	"net/http"
	"note-service/routers"
)

func main() {
	router := routers.SetupRouter()
	log.Println("Server is running on port 80")
	log.Fatal(http.ListenAndServe(":80", router))
}
