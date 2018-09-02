package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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
	Route{
		"AnswerCaptcha",
		"POST",
		"/answer",
		AnswerCaptcha,
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

	ip := strings.Split(r.RemoteAddr, ":")[0]
	resp := generateCaptcha(serverConfig.NumImgsCaptcha, ip)
	jsonResp, err := json.Marshal(resp)
	check(err)
	fmt.Fprintln(w, string(jsonResp))
}

func AnswerCaptcha(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var captchaAnswer CaptchaAnswer
	err := decoder.Decode(&captchaAnswer)
	check(err)
	defer r.Body.Close()

	ip := strings.Split(r.RemoteAddr, ":")[0]
	resp := validateCaptcha(captchaAnswer, ip)
	//if resp==false add ip to blacklist
	if resp == false {
		//get SuspiciousIP
		suspiciousIP := SuspiciousIP{}
		err := suspiciousIPCollection.Find(bson.M{"ip": ip}).One(&suspiciousIP)
		if err != nil {
			//if not exist, add
			suspiciousIP.Date = time.Now().Unix()
			suspiciousIP.IP = ip
			suspiciousIP.Count = 0
			//store suspiciousIP in MongoDB
			err := suspiciousIPCollection.Insert(suspiciousIP)
			check(err)
		} else {
			//if exist
			suspiciousIP.Date = time.Now().Unix()
			suspiciousIP.Count++
			err := suspiciousIPCollection.Update(bson.M{"ip": ip}, suspiciousIP)
			check(err)
		}
	}
	//if count > limit, resp=false
	suspiciousIP := SuspiciousIP{}
	err = suspiciousIPCollection.Find(bson.M{"ip": ip}).One(&suspiciousIP)
	if err == nil {
		//if exist, and time.Since(suspiciousIP.Date) < serverConfig.TimeBan, increase counter
		if time.Since(time.Unix(suspiciousIP.Date, 0)).Seconds() < serverConfig.TimeBan {
			if suspiciousIP.Count > serverConfig.SuspiciousIPCountLimit {
				if resp == false {
					log.Println("IP: " + ip + ", has reached limit count of SuspiciousIP")
				}
				resp = false
			}
		} else {
			//timeBan is completed, delete counter
			err := suspiciousIPCollection.Remove(bson.M{"ip": ip})
			check(err)
		}
		if resp == true {
			//answered correct, delete counter
			err := suspiciousIPCollection.Remove(bson.M{"ip": ip})
			check(err)
		}
	}

	// delete the captchaSolution from MongoDB
	captchaSolCollection.RemoveAll(bson.M{"id": captchaAnswer.CaptchaId})
	check(err)
	// delete the fakepaths from MongoDB
	imgFakePathCollection.RemoveAll(bson.M{"captchaid": captchaAnswer.CaptchaId})
	check(err)

	jsonResp, err := json.Marshal(resp)
	check(err)
	fmt.Fprintln(w, string(jsonResp))
}
