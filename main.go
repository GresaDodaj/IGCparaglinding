package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"log"
	"math/rand"
	"net/http"
	"os"

	"time"

	"github.com/gorilla/mux"
	"github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/mongo"
	//"github.com/mongodb/mongo-go-driver/mongo"
)

//time since the server starts
var startTime = time.Now()
var urlMap = make(map[int]string)
var mapID int
var initialID int
var uniqueId int
var collection = connectToDB()

type url struct {
	URL string `json:"url"`
}
//per mos me i marr krejt te dhanat nga track file e bojna ket strukture  per te cilat te dhena na duhen
type trackFile struct {
	Pilot string
	H_date string
	Glider string
	GliderID string
	TrackLength string
	Url string
	UniqueID string
}

//saves the igc files tracks
var IGC_files []Track

//Struct that saves the ID and igcTrack data
type Track struct {
	ID        string    `json:"ID"`
	IGC_Track igc.Track `json:"igcTrack"`
}

//Struct that saves meta information about the server
type MetaInfo struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

func checkUrl(collection *mongo.Collection,url string)int64{

	filter := bson.NewDocument(bson.EC.String("url",""+url+""))
	length, err := collection.Count(context.Background(),filter)

	if err != nil {
		log.Fatal(err)
	}

	return length
}

// this function returns true if the index is not found and false otherwise

func connectToDB()*mongo.Collection{
	client, err := mongo.NewClient("mongodb://localhost:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("paraglidingDB").Collection("track")
	
	return collection
}
func findIndex(x map[int]string, y int) bool {
	for k, _ := range x {
		if k == y {
			return false
		}
	}
	return true
}

//this function the key of the string if the map contains it, or -1 if the map does not contain the string
func searchMap(x map[int]string, y string) int {

	for k, v := range x {
		if v == y {
			return k
		}
	}
	return -1
}

//this function gets the port assigned by heroku
func GetAddr() string {
	var port = os.Getenv("PORT")

	if port == "" {
		port = "8080"
		fmt.Println("No port  variable detected, defaulting to " + port)
	}
	return ":" + port
}

func IGCinfo(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "Error 404: Page not found!", http.StatusNotFound)
	return
}

func GETapi(w http.ResponseWriter, request *http.Request) {
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

//request is what we get from the client
func getApiIGC(w http.ResponseWriter, request *http.Request) {

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

		trackIDs := make([]string, 0, 0)

		for i := range IGC_files {
			trackIDs = append(trackIDs, IGC_files[i].ID)
		}

		json.NewEncoder(w).Encode(trackIDs)

	case "POST":
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

		//mapID = searchMap(urlMap, URLt.URL)
		initialID = rand.Intn(100)
		trackFileDB := trackFile{}
		//tash ktu me qit checkurl funksion e merr si parameter collection edhe urln qe e kena shkru ne post
		if checkUrl(collection,URLt.URL)==0{
			//nese sosht ne db qajo url atehere e kthen zero edhe ekzekutohet inserti
			//mas = veq pe kthejme uniqueid ne string
			//track osht e mariushit
			track.UniqueID = fmt.Sprintf("%d",initialID)
			//tash ktu ne trackfiledb i shtojme prej track te mariushit veq infot qe na duhen
			trackFileDB = trackFile{track.Pilot,
			track.Date.String(),
			track.GliderType,
			track.GliderID,
			fmt.Sprintf("%f", trackLength(track)),
			URLt.URL,
			track.UniqueID}

			//me insert i shtijme qato te dhana ne databaze
			res, err := collection.InsertOne(context.Background(), trackFileDB)
			if err != nil {
				log.Fatal(err)
			}
			id := res.InsertedID
//id osht per objectID e mongos nese osht nil dmth insertimi nuk osht bo me sukses se gjithe
			if id == nil {
				http.Error(w, "", 500)
			}
			fmt.Fprint(w, "{\n\t\"id\": \""+track.UniqueID+"\"\n}")
			return
		}else{

			//gjema id e rreshtit ne db qe e ka urln te barabart me URLt.URL
			//select id from track where urlprejpostit=urlt.url


			filter := bson.NewDocument(bson.EC.String("url",URLt.URL))//where urlprejpostit=urlt.url
			//decodde osht perdor per me kthy rreshtin e dbs ne strukture trackFileDB
			err := collection.FindOne(context.Background(),filter).Decode(&trackFileDB) //select * where urlprejpostit=urlt.url

			if err != nil {
				log.Fatal(err)
			}

			fmt.Fprint(w, "{\n\t\"id\": \""+trackFileDB.UniqueID+"\"\n}")//tash e bojna print veq id se qajo findone query ja shoqeron
			//vlerat prej databazes strukres trackFileDB

		}


		//if mapID == -1 {
		//	if findIndex(urlMap, initialID) {
		//		uniqueId = initialID
		//		urlMap[uniqueId] = URLt.URL
		//
		//		igcFile := Track{}
		//		igcFile.ID = strconv.Itoa(uniqueId)
		//		igcFile.IGC_Track = track
		//		IGC_files = append(IGC_files, igcFile)
		//		fmt.Fprint(w, "{\n\t\"id\": \""+igcFile.ID+"\"\n}")
		//
		//		res, err := collection.InsertOne(context.Background(), track)
		//		if err != nil {
		//			log.Fatal(err)
		//		}
		//		id := res.InsertedID
		//
		//		if id == nil {
		//			http.Error(w, "", 500)
		//		}
		//
		//		return
		//	} else {
		//		rand.Seed(time.Now().UnixNano())
		//		uniqueId = rand.Intn(100)
		//		urlMap[uniqueId] = URLt.URL
		//		igcFile := Track{}
		//		igcFile.ID = strconv.Itoa(uniqueId)
		//		igcFile.IGC_Track = track
		//		IGC_files = append(IGC_files, igcFile)
		//		fmt.Fprint(w, "{\n\t\"id\": \""+igcFile.ID+"\"\n}")
		//
		//		res, err := collection.InsertOne(context.Background(), track)
		//		if err != nil {
		//			log.Fatal(err)
		//		}
		//		id := res.InsertedID
		//
		//		if id == nil {
		//			http.Error(w, "", 500)
		//		}
		//
		//		return
		//	}
		//} else {
		//	uniqueId = searchMap(urlMap, URLt.URL)
		//	fmt.Fprint(w, "{\n\t\"id\": \""+fmt.Sprintf("%d", uniqueId)+"\"\n}")
		//	return
		//}

	default:
		http.Error(w, "This method is not implemented!", 501)
		return

	}

}

func getApiIgcID(w http.ResponseWriter, request *http.Request) {

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

	filter := bson.NewDocument(bson.EC.String("uniqueid",URLt["id"]))//where uniqueid=urlt["id"] qikjo mas barazimit osht mux variabla te url aty ..../64 qikjo id
	// decodde osht perdor per me kthy rreshtin e dbs ne strukture trackFileDB
	err := collection.FindOne(context.Background(),filter).Decode(&trackFileDB)

	if err != nil {
		log.Fatal(err)
	}


	fmt.Fprint(w, "{\n\"H_date\": \""+trackFileDB.H_date+"\",\n\"pilot\": " +
				"\""+trackFileDB.Pilot+"\",\n\"GliderType\": \""+trackFileDB.Glider+"\",\n\"Glider_ID\": " +
				"\""+trackFileDB.GliderID+"\",\n\"track_length\": \""+trackFileDB.TrackLength+"\"\n}")

	//for i := range IGC_files {
	//	//The requested meta information about a particular track based on the ID given in the url
	//	//checking if the meta information about it is in memory if so the meta information will be returned
	//	//otherwise it will return error 404, not found
	//	if IGC_files[i].ID == URLt["id"] {
	//		tDate := IGC_files[i].IGC_Track.Date.String()
	//		tPilot := IGC_files[i].IGC_Track.Pilot
	//		tGlider := IGC_files[i].IGC_Track.GliderType
	//		tGliderId := IGC_files[i].IGC_Track.GliderID
	//		tTrackLength := fmt.Sprintf("%f", trackLength(IGC_files[i].IGC_Track))
	//		w.Header().Set("content-type", "application/json")
	//		fmt.Fprint(w, "{\n\"H_date\": \""+tDate+"\",\n\"pilot\": " +
	//			"\""+tPilot+"\",\n\"GliderType\": \""+tGlider+"\",\n\"Glider_ID\": " +
	//			"\""+tGliderId+"\",\n\"track_length\": \""+tTrackLength+"\"\n}")
	//	} else {
	//		http.Error(w, "", 404)
	//	}
	//}

}

func getApiIgcIDField(w http.ResponseWriter, request *http.Request) {

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

	filter := bson.NewDocument(bson.EC.String("uniqueid",URLs["id"]))//where uniqueid=urlt["id"] qikjo mas barazimit osht mux variabla te url aty ..../64 qikjo id
	// decodde osht perdor per me kthy rreshtin e dbs ne strukture trackFileDB
	err := collection.FindOne(context.Background(),filter).Decode(&trackFileDB)

	if err != nil {
		log.Fatal(err)
	}
   //te url te pathi qe e shkrujna ../field e merr qita edhe me switch qka tosht e qet 
	switch URLs["field"] {

	case "pilot": fmt.Fprint(w,trackFileDB.Pilot)
	case "h_date": fmt.Fprint(w,trackFileDB.H_date)
	case "glider": fmt.Fprint(w,trackFileDB.Glider)
	case "glider_id": fmt.Fprint(w,trackFileDB.GliderID)
	case "track_length": fmt.Fprint(w,trackFileDB.TrackLength)
	default:
		http.Error(w,"Not found",404)
	}

	//for i := range IGC_files {
	//
	//	if IGC_files[i].ID == URLs["id"] {
	//
	//		mapping := map[string]string{
	//			"pilot":        IGC_files[i].IGC_Track.Pilot,
	//			"glider":       IGC_files[i].IGC_Track.GliderType,
	//			"glider_id":    IGC_files[i].IGC_Track.GliderID,
	//			"track_length": fmt.Sprintf("%f", trackLength(IGC_files[i].IGC_Track)),
	//			"h_date":       IGC_files[i].IGC_Track.Date.String(),
	//		}
	//
	//		field := URLs["field"]
	//		field = strings.ToLower(field)
	//
	//		if val, ok := mapping[field]; ok {
	//			fmt.Fprint(w, val)
	//		} else {
	//
	//			http.Error(w, "", 404)
	//
	//			return
	//		}
	//
	//	}
	//
	//}
}

//function calculating the total  distance of the flight, from the start point until end point(geographical coordinates)
func trackLength(track igc.Track) float64 {

	totalDistance := 0.0

	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	return totalDistance
}

// function that returns the current uptime of the service, format as specified by ISO 8601.
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

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/paragliding/", IGCinfo)
	router.HandleFunc("/paragliding/api", GETapi)
	router.HandleFunc("/paragliding/api/track", getApiIGC)
	//qikjo {id} osht mux.vars
	router.HandleFunc("/paragliding/api/track/{id}", getApiIgcID)
	//ktu edhe field osht njo prej mux.vars
	router.HandleFunc("/paragliding/api/track/{id}/{field}", getApiIgcIDField)

	//err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if err := http.ListenAndServe(":8080", router); err != nil {
		//if err != nil {
		log.Fatal(err)
		//	log.Fatal("ListenAndServe: ", err)
	}

}
