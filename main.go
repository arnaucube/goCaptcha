package main

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
)

func main() {
	savelog()
	log.Println("goCaptcha started")
	rand.Seed(time.Now().UTC().UnixNano())

	//connect with mongodb
	readMongodbConfig("./mongodbConfig.json")
	session, err := getSession()
	check(err)
	//captchaCollection = getCollection(session, "captchas")
	captchaSolCollection = getCollection(session, "captchassolutions")
	imgFakePathCollection = getCollection(session, "imgfakepath")
	suspiciousIPCollection = getCollection(session, "suspiciousip")

	//start the server
	//http server start
	readServerConfig("./serverConfig.json")

	//read the filenames of the dataset
	readDataset(serverConfig.ImgsFolder)
	log.Println("dataset read")
	log.Println("num of dataset categories: " + strconv.Itoa(len(dataset)))

	log.Println("server running")
	log.Print("port: ")
	log.Println(serverConfig.ServerPort)
	router := NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":"+serverConfig.ServerPort, handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
