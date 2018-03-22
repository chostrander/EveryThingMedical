package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"gopkg.in/mgo.v2"
)

/** Reading - Structure for reading results **/
type Reading struct {
	ReadingTime string
	Systolic    int8
	Diastolic   int8
	Notes       string
}

func main() {

	tlsConfig := &tls.Config{}

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{"hoss00-shard-00-00-qyw8j.mongodb.net:27017"},
		Database: "admin", // this is really the authentication DB
		Username: "hossville",
		Password: "USMC##mongo18",
	}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)

	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection := session.DB("BloodPressureTrakker").C("Readings")
	cntToDelete, err1 := collection.Count()
	if err1 != nil {
		panic(err1)
	}
	if cntToDelete > 10 {
		collection.RemoveAll(nil)
	}

	currentTime := time.Now()
	timeStampString := currentTime.Format("2006-01-02 15:04:05")

	currentReading := Reading{}
	currentReading.ReadingTime = timeStampString
	currentReading.Diastolic = 80
	currentReading.Systolic = 111
	currentReading.Notes = "Happy Happy Joy Joy"

	collection.Insert(&currentReading)

	var allReadings []Reading

	collection.Find(nil).All(&allReadings)

	var counter int8
	counter = 0

	for _, v := range allReadings {
		fmt.Println(counter, "--->", v)
		counter++
	}
}
