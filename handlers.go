package main

import (
	"time"
	"gopkg.in/mgo.v2/bson"
	"github.com/mitchellh/mapstructure"
)

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Letter struct {
	MongoId             bson.ObjectId `bson:"_id" json:"_id"`
	LetterId            string `bson:"letterId" json:"letterId"`
	MessageClass        string `bson:"messageClass" json:"messageClass"`
	MessageId           string `bson:"messageId" json:"messageId"`
	ReportName          string `bson:"reportName" json:"reportName"`
	ExternalStudentId   string `bson:"externalStudentId" json:"externalStudentId"`
	CustomFields        map[string]string `bson:"customFields" json:"customFields"`
	AwardYear           int `bson:"awardYear" json:"awardYear"`
	LetterCode          string `bson:"letterCode" json:"letterCode"`
	Status              string `bson:"status" json:"status"`
	LetterEventDateTime time.Time `bson:"letterEventDateTime" json:"letterEventDateTime"`
}

type QueryMessage struct {
	Collection string `json:"collection"`
	Query      string `json:"query"`
}

func queryExecute(client *Client, data interface{}) {

	var queryMessage QueryMessage

	err := mapstructure.Decode(data, &queryMessage)
	if err != nil {
		client.sendChannel <- Message{"error", err.Error()}
	}
	session := client.session.Copy()
	defer session.Close()

	c := session.DB("vm").C("letters")

	var letters []Letter

	err = c.Find(bson.M{}).Limit(5).All(&letters)
	if err != nil {
		client.sendChannel <- Message{"error", err.Error()}
		return
	}

	for _, letter := range letters {
		client.sendChannel <- Message{"entry add", letter}
	}
}
