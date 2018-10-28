package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

var coll = connectToDB("webhooks")

//WEBHOOKForm is used
type WEBHOOKForm struct {
	WEBHOOKURL      string `json:"WEBHOOKURL"`
	MINTRIGGERVALUE int    `json:"MINTRIGGERVALUE"`
	WEBHOOKid       string
}

func returnTracks(n int64, x int64) string {

	trackFileDB := trackFile{}

	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		//http.Error(w, "Bad request!", 400)

	}
	resp := ""
	var i int64 = 0

	for cur.Next(context.Background()) {
		cur.Decode(&trackFileDB)
		if i >= (n - x) {
			resp += trackFileDB.UniqueID
			resp += ","
		}
		i++
	}
	resp = strings.TrimRight(resp, ",")
	return resp
}
func insertWebHookToDB(collection *mongo.Collection, webhook WEBHOOKForm) {

	res, err := collection.InsertOne(context.Background(), webhook)

	if err != nil {
		log.Fatal(err)
	}
	id := res.InsertedID
	fmt.Print(id)
	if id == nil {
		fmt.Print("ID nil!")
	}

}

//WebHookHandler is used
func WebHookHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		webhookInfo := WEBHOOKForm{}

		err := json.NewDecoder(r.Body).Decode(&webhookInfo)
		if err != nil {
			http.Error(w, "Bad request!", 400)
			return
		}

		if webhookInfo.WEBHOOKURL == "" {
			http.Error(w, "Bad request!", 400)
			return
		}
		if err != nil {
			//http.Error(w, "Bad request!", 400)

		}
		if webhookInfo.MINTRIGGERVALUE == 0 {

			webhookInfo.MINTRIGGERVALUE = 1
		}

		if checkUrl(coll, webhookInfo.WEBHOOKURL, "webhookurl") == 0 {

			rand.Seed(time.Now().UnixNano())
			initialID = rand.Intn(100)
			webhookInfo.WEBHOOKid = fmt.Sprintf("%d", initialID)

			insertWebHookToDB(coll, webhookInfo)

		}

		filter := bson.NewDocument(bson.EC.String("webhookurl", ""+webhookInfo.WEBHOOKURL+""))
		error := coll.FindOne(context.Background(), filter).Decode(&webhookInfo)
		if error != nil {
			log.Fatal(err)
		}

		resp := "{\n\"id\": " + "\"" + webhookInfo.WEBHOOKid + "\"\n}" //formating the response in json format

		fmt.Fprint(w, resp)

	}

}
func WebHookHandlerID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-type", "application/json")

	URLt := mux.Vars(r)
	if len(URLt) != 1 {
		http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
		return
	}
	if URLt["webhookID"] == "" {
		http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
		return
	}

	webhookDB := WEBHOOKForm{}
	filter := bson.NewDocument(bson.EC.String("webhookid", URLt["webhookID"])) //where uniqueid=urlt["webhookID"] qikjo mas barazimit osht mux variabla te url aty ..../64 qikjo id
	// decodde osht perdor per me kthy rreshtin e dbs ne strukture trackFileDB
	err := coll.FindOne(context.Background(), filter).Decode(&webhookDB)
	if err != nil {
		http.Error(w, "", 501)
		return
	}

	if r.Method == http.MethodGet {

		fmt.Fprint(w, "{\n\"webhookURL\": \""+webhookDB.WEBHOOKURL+"\",\n\"minTriggerValue\": "+
			"\""+fmt.Sprint(webhookDB.MINTRIGGERVALUE)+"\"\n}")

		//fmt.Fprint(w, response)

	} else if r.Method == http.MethodDelete {
		http.Error(w, "", 501)
		del, err := coll.DeleteOne(context.Background(), filter)
		if err != nil {
			http.Error(w, "", 501)
			return
		}
		if del.DeletedCount == 0 { //nese delete fail
			http.Error(w, "", 501)
			return
		}
		// response := "{"
		// response = `"webhookURL":"` + webhookDB.WEBHOOKURL + `"\n`
		// response = `"minTriggerValue":"` + fmt.Sprint(webhookDB.MINTRIGGERVALUE) + `"\n`
		//response = "}"
		fmt.Fprint(w, "{\n\"webhookURL\": \""+webhookDB.WEBHOOKURL+"\",\n\"minTriggerValue\": "+
			"\""+fmt.Sprint(webhookDB.MINTRIGGERVALUE)+"\"\n}")

		//fmt.Fprint(w, response)

	} else {
		http.Error(w, "", 501)
	}

}

func triggerWebhook() error{
	webhookinfo := WEBHOOKForm{}

	trackCount, err := collection.Count(context.Background(), nil)
	if err != nil {
		// http.Error(w, "", 400)
		return err
	}
	cursor, err := coll.Find(context.Background(), nil)
	if err != nil {
		// http.Error(w, "", 400)
		return err
	}

	for cursor.Next(context.Background()) {

		cursor.Decode(&webhookinfo)

		if trackCount%int64(webhookinfo.MINTRIGGERVALUE) != 0 {
			continue
		}

		processStart := time.Now() // Track when the process started

		url := webhookinfo.WEBHOOKURL

		trackString := returnTracks(int64(trackCount), int64(webhookinfo.MINTRIGGERVALUE))

		latestTS := tLatest()
		jsonPayload := "{"
		jsonPayload += `"username": "Tracks added",`
		jsonPayload += `"content": "Latest added track at ` + latestTS + `\n`
		jsonPayload += `New tracks are ` + trackString + `\n`
		jsonPayload += `The request took ` + strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 2, 64) + `ms to process"`
		jsonPayload += "}"

		var jsonStr = []byte(jsonPayload)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
	return err
}

func triggerWebhookPeriod() error{
	webhookinfo := WEBHOOKForm{}

	trackCount, err := collection.Count(context.Background(), nil)
	if err != nil {
		// http.Error(w, "", 400)
		return err
	}
	cursor, err := coll.Find(context.Background(), nil)
	if err != nil {
		// http.Error(w, "", 400)
		return err
	}

	for cursor.Next(context.Background()) {

		cursor.Decode(&webhookinfo)

		processStart := time.Now() // Track when the process started

		url := webhookinfo.WEBHOOKURL

		trackString := returnTracks(int64(trackCount), int64(webhookinfo.MINTRIGGERVALUE))

		latestTS := tLatest()
		jsonPayload := "{"
		jsonPayload += `"username": "Tracks added",`
		jsonPayload += `"content": "Latest added track at ` + latestTS + `\n`
		jsonPayload += `New tracks are ` + trackString + `\n`
		jsonPayload += `The request took ` + strconv.FormatFloat(float64(time.Since(processStart))/float64(time.Millisecond), 'f', 2, 64) + `ms to process"`
		jsonPayload += "}"

		var jsonStr = []byte(jsonPayload)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
	return err
}
