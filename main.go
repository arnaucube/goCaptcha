package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {
	savelog()
	log.Println("goCaptcha started")

	//start the server
	//http server start
	readServerConfig("./serverConfig.json")

	log.Println("server running")
	log.Print("port: ")
	log.Println(serverConfig.ServerPort)
	router := NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+serverConfig.ServerPort, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
