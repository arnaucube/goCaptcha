package main

import (
	"math/rand"
	"os/exec"
	"strings"
)

type Captcha struct {
	Id       string   `json:"id"`
	Imgs     []string `json:"imgs"`
	Question string   `json:"question"`
	Date     string   `json:"date"`
}
type CaptchaSolution struct {
	Id           string   `json:"id"`
	Imgs         []string `json:"imgs"`
	ImgsSolution []string `json:"imgssolution"`
	Question     string   `json:"question"`
	Date         string   `json:"date"`
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
	var captchaSol CaptchaSolution

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
	captcha.Question = "Select all leopards"
	captchaSol.Question = "Select all leopards"
	err := captchaCollection.Insert(captcha)
	check(err)
	err = captchaSolutionCollection.Insert(captchaSol)
	check(err)
	return captcha
}
