package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	viper "github.com/spf13/viper"
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

	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath("$HOME/go") // path to look for the config file in
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("read in viper")
	fmt.Printf("--> %s \n", viper.GetString("mongo.username"))
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{"hoss00-shard-00-00-qyw8j.mongodb.net:27017"},
		Database: "admin", // this is really the authentication DB
		Username: viper.GetString("mongo.username"),
		Password: viper.GetString("mongo.password"),
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
