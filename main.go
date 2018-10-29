package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	//"github.com/mongodb/mongo-go-driver/mongo"
)

//time since the server starts
var startTime = time.Now()
var initialID int
var lengthTrig int64
var lengthTrigAfter int64
var collection = connectToDB("track")

type url struct {
	URL string `json:"url"`
}

//trackFile struct is used to get the data we need from an igc file
type trackFile struct {
	Pilot       string
	HDate      string
	Glider      string
	GliderID    string
	TrackLength string
	URL 		string
	UniqueID    string
	TimeStamp   time.Time
}

//Track is a struct that saves the ID and igcTrack data
type Track struct {
	ID       string    `json:"ID"`
	IGCTrack igc.Track `json:"igcTrack"`
}

//MetaInfo is a struct that saves meta information about the server
type MetaInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

//checkURL  func checks if the posted url is already in the database
func checkURL(collection *mongo.Collection, url string, urlDB string) int64 {
	//select * from collection where url(e postit)=urlDB
	//url is the url posted and the urlDB are the urls already in db
	//check if any of the urlDB is the url posted(filter the db documents so that the url posted is equal to one of the urls in DB)
	filter := bson.NewDocument(bson.EC.String(""+urlDB+"", ""+url+""))
	//length is 0 if the url is not in the database
	length, err := collection.Count(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	//return the length if the length is 0 the url will be inserted in the db (where the function is called)
	return length
}

//connectToDB is a function to connect the server to database
func connectToDB(col string) *mongo.Collection {
	client, err := mongo.Connect(context.Background(),"mongodb://gresad:prishtina123@ds113443.mlab.com:13443/paraglidingdb",nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("paraglidingdb").Collection(col)

	return collection
}

//GetAddr function gets the port assigned by heroku
func GetAddr() string {
	var port = os.Getenv("PORT")

	if port == "" {
		port = "8080"
		fmt.Println("No port  variable detected, defaulting to " + port)
	}
	return ":" + port
}
//IGCinfo is a function which returns error if the method used is not get and redirects from /paragliding/ to /paragliding/api
func IGCinfo(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "501 - Method not implemented", http.StatusNotImplemented)
		return
	}
	// Redirect to /paragliding/api
	http.Redirect(w, r, "/paragliding/api", 302)
	return
}

func getAPI(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("content-type", "application/json")

	URLs := mux.Vars(request)
	if len(URLs) != 0 {
		http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
		return
	}

	metaInfo := &MetaInfo{}
	metaInfo.Uptime = FormatSince(startTime)
	metaInfo.Info = "Service for IGC tracks"
	metaInfo.Version = "version 1.0"

	json.NewEncoder(w).Encode(metaInfo)
}
func postAPIigc(w http.ResponseWriter,request *http.Request){



	// Set response content-type to JSON
	w.Header().Set("content-type", "application/json")

	URLt := &url{}

	//Url is given to the server as JSON and now we decode it to a go structure
	var error = json.NewDecoder(request.Body).Decode(URLt)
	if error != nil {
	http.Error(w, http.StatusText(400), 400)
	return
	}

	//making a random unique ID for the track files
	rand.Seed(time.Now().UnixNano())

	track, err := igc.ParseLocation(URLt.URL)
	if err != nil {

	http.Error(w, "Bad request!\nMalformed URL!", 400)
	return
	}

	initialID = rand.Intn(100)
	trackFileDB := trackFile{}
	//checkURL gets the collection, url posted and url from database
	if checkURL(collection, URLt.URL, "url") == 0 {
	//if the check url is 0 then it means that the url posted is not in the database so the insertion is executed
	//we assign the initialID(as string that's why the Sprintf is used
	track.UniqueID = fmt.Sprintf("%d", initialID)
	//from the track file which contains all the data of an igc file we assign values to the trackFileDB object
	trackFileDB = trackFile{track.Pilot,
	track.Date.String(),
	track.GliderType,
	track.GliderID,
	fmt.Sprintf("%f", trackLength(track)),
	URLt.URL,
	track.UniqueID,
	time.Now()}

	lengthTrig, err = collection.Count(context.Background(), nil)
	if err != nil {
	http.Error(w, "", 400)
	return
	}

	//insert that data to the database
	res, err := collection.InsertOne(context.Background(), trackFileDB)
	if err != nil {
	log.Fatal(err)
	}
	id := res.InsertedID
	//id is the objectID of the MongoDB which is always generated as a unique id for every single document
	// if that id is nil(don't have that id) it means that the insertion failed
	if id == nil {
	http.Error(w, "", 500)
	}
	fmt.Fprint(w, "{\n\t\"id\": \""+track.UniqueID+"\"\n}")

	err = triggerWebhook()
	if err != nil {
	http.Error(w, "", 400)
	return
	}
	lengthTrigAfter, err = collection.Count(context.Background(), nil)
	if err != nil {
	http.Error(w, "", 400)
	return
	}

	return
	}

	//analogy: select id from track where urlprejpostit=urlt.url
	//if the checkURL is not false then find the id of that igc file and print it
	filter := bson.NewDocument(bson.EC.String("url", URLt.URL)) //where urlprejpostit=urlt.url
	//decode is used to convert the document from the db to the trackFileDB structure
	//FindOne because we are filtering them by url so it means that if that url is in db it's only added once so we after it's found one url
	//that is the same as the url posted, it doesn't need to keep searching in the db for other urls
	err = collection.FindOne(context.Background(), filter).Decode(&trackFileDB) //select * where urlprejpostit=urlt.url

	if err != nil {
	log.Fatal(err)
	}
	//print only the id of that file
	fmt.Fprint(w, "{\n\t\"id\": \""+trackFileDB.UniqueID+"\"\n}")
}
//getAPIigc returns the id of the igc file if the request method used by the client is POST or returns the ids of igc files already in the db
//if the request method is GET
func getAPIigc(w http.ResponseWriter, request *http.Request) {

	//request.method gives us the method selected by the client, in this api there are two methods
	//that are implemented GET and POST, requests made for other methods will result to an error 501
	//501 is an HTTP  error for not implemented
	switch request.Method {

	case "GET":
		w.Header().Set("content-type", "application/json")

		URLs := mux.Vars(request)
		if len(URLs) != 0 {
			http.Error(w, "400 - Bad Request!", 400)
			return
		}

		trackFileDB := trackFile{}
		//filter is nil because we need all the ids which means we don't need to filter them based on anything
		cur, err := collection.Find(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		//we store the ids on the ids variable
		ids := "["
		//length=number of documents in the db
		length, err1 := collection.Count(context.Background(), nil)
		if err1 != nil {
			log.Fatal(err1)
		}
		i := int64(0) //we use int64 for i because the length is  int64 data and we set the value 0
		//cur.Next returns true if there is a next document in the db, it returns false in the last document
		for cur.Next(context.Background()) {
			//decode the document from the database into the trackFileDB struct(we get only the data we need)
			cur.Decode(&trackFileDB)
			//we add the uniqueID of the trackFileDB into the ids array
			ids += trackFileDB.UniqueID
			//so that after the last member there is no comma
			if i == length-1 {
				break
			}
			ids += ","
			i++
		}
		ids += "]"

		fmt.Fprint(w, ids)


	default:
		http.Error(w, "This method is not implemented!", 501)
		return

	}

}

func getAPIIgcID(w http.ResponseWriter, request *http.Request) {

	w.Header().Set("content-type", "application/json")

	URLt := mux.Vars(request)
	if len(URLt) != 1 {
		http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
		return
	}

	if URLt["id"] == "" {
		http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
		return
	}

	trackFileDB := trackFile{}

	filter := bson.NewDocument(bson.EC.String("uniqueid", URLt["id"]))
	err := collection.FindOne(context.Background(), filter).Decode(&trackFileDB)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, "{\n\"H_date\": \""+trackFileDB.HDate+"\",\n\"pilot\": "+
		"\""+trackFileDB.Pilot+"\",\n\"GliderType\": \""+trackFileDB.Glider+"\",\n\"Glider_ID\": "+
		"\""+trackFileDB.GliderID+"\",\n\"track_length\": \""+trackFileDB.TrackLength+"\""+
		",\n\"track_src_url\": \""+trackFileDB.URL+"\"\n}")

}

func getAPIIgcIDField(w http.ResponseWriter, request *http.Request) {

	URLs := mux.Vars(request)
	if len(URLs) != 2 {
		w.Header().Set("content-type", "application/json")
		http.Error(w, "Error 400 : Bad Request!", http.StatusBadRequest)
		return
	}

	if URLs["id"] == "" {
		w.Header().Set("content-type", "application/json")
		http.Error(w, "Error 400 : Bad Request!\n You did not enter an ID.", http.StatusBadRequest)
		return
	}

	if URLs["field"] == "" {
		w.Header().Set("content-type", "application/json")
		http.Error(w, "Error 400 : Bad Request!\n You did not  enter a field.", http.StatusBadRequest)
		return
	}

	trackFileDB := trackFile{}

	filter := bson.NewDocument(bson.EC.String("uniqueid", URLs["id"]))
	err := collection.FindOne(context.Background(), filter).Decode(&trackFileDB)

	if err != nil {
		log.Fatal(err)
	}
	//url path .../field switch
	switch URLs["field"] {

	case "pilot":
		fmt.Fprint(w, trackFileDB.Pilot)
	case "h_date":
		fmt.Fprint(w, trackFileDB.HDate)
	case "glider":
		fmt.Fprint(w, trackFileDB.Glider)
	case "glider_id":
		fmt.Fprint(w, trackFileDB.GliderID)
	case "track_length":
		fmt.Fprint(w, trackFileDB.TrackLength)
	default:
		http.Error(w, "Not found", 404)
	}

}

func getAPITickerLatest(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, tLatest())

}
func getAPITicker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	sTime := time.Now()
	tLatest := ""
	tStart := ""
	tStop := ""
	tracksStr := "["

	trackFileDB := trackFile{}

	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	//length=number of documents in the collection
	length, err1 := collection.Count(context.Background(), nil)
	if err1 != nil {
		log.Fatal(err1)
	}
	i := int64(0)
	//cur.Next returns true if there is a next document in the db, it returns false in the last document
	for cur.Next(context.Background()) {
		//decoding data from database to trackFileDB struct
		cur.Decode(&trackFileDB)

		//five tracks 0 1 2 3 4
		if i <= 4 {
			tracksStr += trackFileDB.UniqueID

		}
		//tStart is the time of the first track added in db
		if i == 0 {
			tStart = fmt.Sprint(trackFileDB.TimeStamp)
		}
		//the last document added is  length-1 -->tLatest
		if i == length-1 {
			tLatest = fmt.Sprint(trackFileDB.TimeStamp)

		} else if i < 4 {
			tracksStr += ","
		}
		//if there are more than 5 tracks then make the 5th track file added time (i=4) the tStop
		if length > 4 {
			//te requiremets cap=5 01234 :
			if i == 4 {
				tStop = fmt.Sprint(trackFileDB.TimeStamp)

			}
		} else {
			//if there are less then 5 then make the latest added the tstop
			tStop = tLatest
		}

		i++
	}
	tracksStr += "]"
	fmt.Fprint(w, "{\n\"tLatest\": \""+tLatest+"\",\n\"tStart\": "+
		"\""+tStart+"\",\n\"tStop\": \""+tStop+"\",\n\"tracks\": "+
		"\""+tracksStr+"\",\n\"processing\": \""+time.Since(sTime).String()+"\"\n}")

}
func getJ(collection *mongo.Collection, a string) int64 {
	trackFileDB := trackFile{}

	cur, err := collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	var i int64
	var j int64
	//while we have documents in db then:
	for cur.Next(context.Background()) {
		err := cur.Decode(&trackFileDB)
		if err != nil {
			log.Fatal(err)
		}
		//if timestamp given in the url path "a" is in trackFileDB.TimeStamp.String() then save that to the variable j
		//that timestamp and return it
		if trackFileDB.TimeStamp.String() == a {
			j = i
			break
		}
		i++
	}
	return j
}
//getAPITickerTimeStamp is a function that returns the JSON struct representing the ticker for the IGC tracks.
func getAPITickerTimeStamp(w http.ResponseWriter, r *http.Request) {
	URLt := mux.Vars(r)
	if len(URLt) != 1 {
		http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
		return
	}

	if URLt["timestamp"] == "" {
		http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
		return
	}
	w.Header().Set("content-type", "application/json")
	resp,_ := respHandler(URLt["timestamp"])
	fmt.Fprint(w, resp)

}
func respHandler(x string)(string,int64){

	sTime := time.Now()
	tLatest := ""
	tStart := ""
	tStop := ""
	tracksStr := "["

	trackFileDB := trackFile{}

	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	//length is the number of documents in the db
	length, err1 := collection.Count(context.Background(), nil)
	if err1 != nil {
		log.Fatal(err1)
	}
	i := int64(0)
	j := getJ(collection, x)

	for cur.Next(context.Background()) {
		//the data from the db is converted to the trackFileDB structure
		cur.Decode(&trackFileDB)

		if i > j && i <= j+5 {
			tracksStr += trackFileDB.UniqueID

		}

		if i == j+1 {
			tStart = fmt.Sprint(trackFileDB.TimeStamp)
		}

		if i == length-1 {
			tLatest = fmt.Sprint(trackFileDB.TimeStamp)

		} else if i > j && i < j+5 {
			tracksStr += ","
		}

		if length > j+5 {

			if i == j+5 {
				tStop = fmt.Sprint(trackFileDB.TimeStamp)

			}
		} else {

			tStop = tLatest
		}

		i++
	}
	tracksStr += "]"
	resp := "{\n\"tLatest\": \""+ tLatest +"\",\n\"tStart\": "+
		"\""+ tStart +"\",\n\"tStop\": \""+ tStop +"\",\n\"tracks\": "+
		"\""+tracksStr+"\",\n\"processing\": \""+time.Since(sTime).String()+"\"\n}"
	return resp,j
}
//function calculating the total  distance of the flight, from the start point until end point(geographical coordinates)
func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}

//FormatSince that returns the current uptime of the service, format as specified by ISO 8601.
func FormatSince(t time.Time) string {
	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
	)
	ts := time.Since(t)
	sign := time.Duration(1)
	if ts < 0 {
		sign = -1
		ts = -ts
	}
	ts += +Decisecond / 2
	d := sign * (ts / Day)
	ts = ts % Day
	h := ts / time.Hour
	ts = ts % time.Hour
	m := ts / time.Minute
	ts = ts % time.Minute
	s := ts / time.Second
	ts = ts % time.Second
	f := ts / Decisecond
	y := d / 365
	return fmt.Sprintf("P%dY%dD%dH%dM%d.%dS", y, d, h, m, s, f)
}

func tLatest() string {
	trackFileDB := trackFile{}

	cur, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	length, err1 := collection.Count(context.Background(), nil)
	if err1 != nil {
		log.Fatal(err1)
	}
	respons := ""
	i := int64(0)
	for cur.Next(context.Background()) {

		cur.Decode(&trackFileDB)

		if i == length-1 {
			respons = fmt.Sprint(trackFileDB.TimeStamp)
		}

		i++
	}

	return respons

}

func main() {


	router := mux.NewRouter()

	router.HandleFunc("/paragliding/", IGCinfo)
	router.HandleFunc("/paragliding/api", getAPI)
	router.HandleFunc("/paragliding/api/track", getAPIigc)
	router.HandleFunc("/paragliding/api/track", postAPIigc).Methods("POST")
	router.HandleFunc("/paragliding/api/track/{id}", getAPIIgcID)
	router.HandleFunc("/paragliding/api/track/{id}/{field}", getAPIIgcIDField)
	router.HandleFunc("/paragliding/api/ticker/latest", getAPITickerLatest)
	router.HandleFunc("/paragliding/api/ticker", getAPITicker)
	router.HandleFunc("/paragliding/api/ticker/{timestamp}", getAPITickerTimeStamp)
	router.HandleFunc("/paragliding/admin/api/webhook",adminClockTrigger)
	router.HandleFunc("/api/webhook/new_track/", WebHookHandler)
	router.HandleFunc("/api/webhook/new_track/{webhookID}", WebHookHandlerID)
	router.HandleFunc("/admin/api/tracks_count", AdminHandlerGet)
	router.HandleFunc("/admin/api/tracks", AdminHandlerDelete)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
	//if err := http.ListenAndServe(":8080", router); err != nil {
		if err != nil {
		//log.Fatal(err)
			log.Fatal("ListenAndServe: ", err)
	}

}
