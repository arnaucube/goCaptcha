package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

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

	imgFakePath := ImgFakePath{}
	err := imgFakePathCollection.Find(bson.M{"fake": imageName}).One(&imgFakePath)
	check(err)

	fmt.Println(serverConfig.ImgsFolder + "/" + imgFakePath.Real)

	file, err := ioutil.ReadFile(serverConfig.ImgsFolder + "/" + imgFakePath.Real)
	if err != nil {
		fmt.Fprintln(w, err)
	}

	pathSplited := strings.Split(imgFakePath.Real, ".")
	imageExtension := pathSplited[len(pathSplited)-1]
	img, err := dataToImage(file, imageExtension)

	if err != nil {
		fmt.Fprintln(w, "image "+imageName+" does not exist in server")
	} else {
		jpeg.Encode(w, img, nil) // Write to the ResponseWriter
	}
}
func GetCaptcha(w http.ResponseWriter, r *http.Request) {

	resp := generateCaptcha(6)
	jsonResp, err := json.Marshal(resp)
	check(err)
	fmt.Fprintln(w, string(jsonResp))
}
