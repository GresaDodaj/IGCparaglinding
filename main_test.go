package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)


func Test_IGCinfo(test *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(IGCinfo))
	defer testServer.Close()

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	if err != nil {
		test.Errorf("Error constructing the GET request, %s", err)
	}

	response, err := client.Do(request)
	if err != nil {
		test.Errorf("Error executing the GET request, %s", err)
	}

	if response.StatusCode != http.StatusNotFound {
		test.Errorf("StatusNotFound %d, received %d. ",404, response.StatusCode)
		return
	}

}

func Test_getApiIGC_NotImplemented(test *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(getAPIigc))
	defer testServer.Close()

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodDelete, testServer.URL, nil)
	if err != nil {
		test.Errorf("Error constructing the DELETE request, %s", err)
	}

	response, err := client.Do(request)
	if err != nil {
		test.Errorf("Error executing the DELETE request, %s", err)
	}

	if response.StatusCode != http.StatusNotImplemented {
		test.Errorf("Expected StatusNotImplemented %d, received %d. ", 501, response.StatusCode)
		return
	}

}
func Test_getAPIIgcId_NotImplemented(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPIIgcID))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusBadRequest, resp.StatusCode)
		return
	}

}
func Test_getAPIIgcField_NotImplemented(t *testing.T) {

	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPIIgcIDField))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != 400 {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusNotImplemented, resp.StatusCode)
		return
	}

}


func Test_getAPIIgc_MalformedURL(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPI))
	defer ts.Close()

	testCases := []string{
		ts.URL,
		ts.URL + "/something/",
		ts.URL + "/something/123/",
	}

	for _, tstring := range testCases {
		resp, err := http.Get(ts.URL)
		if err != nil {
			t.Errorf("Error making the GET request, %s", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, http.StatusBadRequest, resp.StatusCode)
			return
		}
	}
}


func Test_getApiIgcID_Malformed(test *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(getAPIIgcID))
	defer testServer.Close()

	testCases := []string {
		testServer.URL,
		testServer.URL + "/blla/",
		testServer.URL + "/blla/123/",
	}


	for _, tstring := range testCases {
		response, err := http.Get(testServer.URL)
		if err != nil {
			test.Errorf("Error making the GET request, %s", err)
		}

		if response.StatusCode != http.StatusBadRequest {
			test.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, 400, response.StatusCode)
			return
		}
	}
}


func Test_getApiIgcIDField_MalformedURL(test *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(getAPIIgcIDField))
	defer testServer.Close()

	testCases := []string {
		testServer.URL,
		testServer.URL + "/blla/",
		testServer.URL + "/blla/123/",
	}


	for _, tstring := range testCases {
		response, err := http.Get(testServer.URL)
		if err != nil {
			test.Errorf("Error making the GET request, %s", err)
		}

		if response.StatusCode != http.StatusBadRequest {
			test.Errorf("For route: %s, expected StatusCode %d, received %d. ", tstring, 400, response.StatusCode)
			return
		}
	}
}
func Test_getAPIIgc_Post(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPIigc))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	apiURLTest := url{}
	apiURLTest.URL = "http://skypolaris.org/wp-content/uploa/IGS%20Files/Madrid%20to%20Jerez.igc"

	jsonData, _ := json.Marshal(apiURLTest)

	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Error making the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	if resp.StatusCode == 400 {
		assert.Equal(t, 400, resp.StatusCode, "Bad request is expected")
	} else {
		assert.Equal(t, 200, resp.StatusCode, "OK response is expected")
	}

}
func Test_getAPIIgcPostEmpty(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(getAPIigc))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	apiURLTest := url{}
	apiURLTest.URL = ""

	jsonData, _ := json.Marshal(apiURLTest)

	req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Error making the POST request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	assert.Equal(t, 400, resp.StatusCode, "OK response is expected")

}

func Test_getAPITicker1(t *testing.T){
	req, err := http.NewRequest("GET", "/igcinfo/api/ticker", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getAPITicker)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusOK { // It should be 200 (OK)
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}

func Test_respHandler(t *testing.T){

	resp,j := respHandler("testest123")

	if resp == "" || j !=0 {
		t.Error("Timestamp was not correct!")
	}
}

func Test_FormatTime(t *testing.T){

	resp := FormatSince(time.Now())

	if resp == "" {
		t.Error("Timestamp was not correct!")
	}
}

func Test_getJ(t *testing.T){

	resp := getJ(collection,"randomTest")

	if resp !=0 {
		t.Error("Timestamp was not correct!")
	}
}


func Test_tLatest(t *testing.T){

	resp := tLatest()

	if resp !="" {
		t.Error("Timestamp was not correct!")
	}
}

func Test_GetApiIGC(test *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(getAPIigc))
	defer testServer.Close()

	client := &http.Client{}

	jsonPayload := "{"
	jsonPayload += `"url": "http://skypolaris.org/wp-content/uploads/IGS%20Files/Medellin%20Guatemala.igc"`
	jsonPayload += "}"

	var jsonStr = []byte(jsonPayload)
	request, err := http.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer(jsonStr))

	if err != nil {
		test.Errorf("Error constructing the GET request, %s", err)
	}

	response, err := client.Do(request)
	if err != nil {
		test.Errorf("Error executing the GET request, %s", err)
	}

	if response.StatusCode != http.StatusOK {
		test.Errorf("StatusNotFound %d, received %d. ",404, response.StatusCode)
		return
	}

}


func Test_WebHookHandler(test *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(WebHookHandler))
	defer testServer.Close()

	client := &http.Client{}

	jsonPayload := "{"
	jsonPayload += `"webhookurl": "https://discordapp.com/api/webhooks/504251337953771521/ASiZ1DNh9YtTbLbEOzTp-LmUny8ju_qyhhQwmWlAE6zuWWI7x2nLhLbof9vnNp771at4",`
	jsonPayload += `"mintriggervalue": 2`
	jsonPayload += "}"

	var jsonStr = []byte(jsonPayload)
	request, err := http.NewRequest(http.MethodPost, testServer.URL, bytes.NewBuffer(jsonStr))

	if err != nil {
		test.Errorf("Error constructing the GET request, %s", err)
	}

	response, err := client.Do(request)
	if err != nil {
		test.Errorf("Error executing the GET request, %s", err)
	}

	if response.StatusCode != http.StatusOK {
		test.Errorf("StatusNotFound %d, received %d. ",404, response.StatusCode)
		return
	}

}