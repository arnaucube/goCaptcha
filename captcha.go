package main

import (
	"math/rand"
	"os/exec"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type Captcha struct {
	Id       string   `json:"id"`
	Imgs     []string `json:"imgs"`
	Question string   `json:"question"`
	Date     string   `json:"date"`
}
type CaptchaSol struct {
	Id           string   `json:"id"`
	Imgs         []string `json:"imgs"`
	ImgsSolution []string `json:"imgssolution"`
	Question     string   `json:"question"` //select all X
	Date         string   `json:"date"`
}
type CaptchaAnswer struct {
	CaptchaId string `json:"captchaid"`
	Selection []int  `json:"selection"`
}
type ImgFakePath struct {
	CaptchaId string `json:"captchaid"`
	Real      string `json:"real"`
	Fake      string `json:"fake"`
}

func generateUUID() string {
	out, err := exec.Command("uuidgen").Output()
	check(err)
	uuid := string(out)
	uuid = strings.Replace(uuid, "\n", "", -1)
	return uuid
}
func generateRandInt(min int, max int) int {
	//rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}
func generateCaptcha(count int) Captcha {
	var captcha Captcha
	var captchaSol CaptchaSol

	captcha.Id = generateUUID()
	captchaSol.Id = captcha.Id

	for i := 0; i < count; i++ {
		nCateg := generateRandInt(0, len(categDataset))
		nImg := generateRandInt(0, len(dataset[categDataset[nCateg]]))
		//imgFakePath
		var imgFakePath ImgFakePath
		imgFakePath.CaptchaId = captcha.Id
		imgFakePath.Real = categDataset[nCateg] + "/" + dataset[categDataset[nCateg]][nImg]
		imgFakePath.Fake = generateUUID() + ".png"
		err := imgFakePathCollection.Insert(imgFakePath)
		check(err)
		captcha.Imgs = append(captcha.Imgs, imgFakePath.Fake)
		captchaSol.Imgs = append(captchaSol.Imgs, dataset[categDataset[nCateg]][nImg])
		captchaSol.ImgsSolution = append(captchaSol.ImgsSolution, categDataset[nCateg])
	}
	captcha.Question = "leopard"
	captchaSol.Question = "leopard"
	err := captchaCollection.Insert(captcha)
	check(err)
	err = captchaSolCollection.Insert(captchaSol)
	check(err)
	return captcha
}
func validateCaptcha(captchaAnswer CaptchaAnswer) bool {
	var solved bool
	solved = true
	captchaSol := CaptchaSol{}
	err := captchaSolCollection.Find(bson.M{"id": captchaAnswer.CaptchaId}).One(&captchaSol)
	check(err)
	for k, imgSol := range captchaSol.ImgsSolution {
		if imgSol == captchaSol.Question {
			if captchaAnswer.Selection[k] == 1 {
				//correct
			} else {
				solved = false
			}
		}
		if imgSol != captchaSol.Question {
			if captchaAnswer.Selection[k] == 0 {
				//correct
			} else {
				solved = false
			}
		}
	}

	return solved
}
