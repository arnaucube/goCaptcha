package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	mgo "gopkg.in/mgo.v2"
)

//MongoConfig stores the configuration of mongodb to connect
type MongoConfig struct {
	Ip       string `json:"ip"`
	Database string `json:"database"`
}

var mongoConfig MongoConfig

var captchaCollection *mgo.Collection
var captchaSolutionCollection *mgo.Collection
var imgFakePathCollection *mgo.Collection

func readMongodbConfig(path string) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Println("error:", e)
	}
	content := string(file)
	json.Unmarshal([]byte(content), &mongoConfig)
}

func getSession() (*mgo.Session, error) {
	session, err := mgo.Dial("mongodb://" + mongoConfig.Ip)
	if err != nil {
		panic(err)
	}
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	return session, err
}
func getCollection(session *mgo.Session, collection string) *mgo.Collection {

	c := session.DB(mongoConfig.Database).C(collection)
	return c
}
