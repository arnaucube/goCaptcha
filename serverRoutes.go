package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"GetImage",
		"GET",
		"/image/{imageName}",
		GetImage,
	},
	Route{
		"GetCaptcha",
		"GET",
		"/captcha",
		GetCaptcha,
	},
}

//ROUTES

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ask for images in /r")
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageName := vars["imageName"]

	file, err := ioutil.ReadFile(serverConfig.ImgsFolder + "/" + imageName)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	img, err := dataToPNG(file, imageName)

	if err != nil {
		fmt.Fprintln(w, "image "+imageName+" does not exist in server")
	} else {
		jpeg.Encode(w, img, nil) // Write to the ResponseWriter
	}
}
func GetCaptcha(w http.ResponseWriter, r *http.Request) {

	resp := ""
	jsonResp, err := json.Marshal(resp)
	check(err)
	fmt.Fprintln(w, string(jsonResp))
}
