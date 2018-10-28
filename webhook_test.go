package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_webhookID(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(WebHookHandler))
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
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusOK %d, received %d. ", http.StatusBadRequest, resp.StatusCode)
		return
	}

}

func Test_WebHookHandlerID(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(WebHookHandlerID))
	defer ts.Close()

	client := &http.Client{}
	webhookInfo := WEBHOOKForm{}

	webhookInfo.WEBHOOKURL = "https://discordapp.com/api/webhooks/504251337953771521/ASiZ1DNh9YtTbLbEOzTp-LmUny8ju_qyhhQwmWlAE6zuWWI7x2nLhLbof9vnNp771at4"

	jsonData, _ := json.Marshal(webhookInfo)

	request, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Errorf("Error making the POST request, %s", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected Bad Request %d, received %d. ", http.StatusBadRequest, resp.StatusCode)
		return
	}

	if resp.StatusCode == 400 {
		assert.Equal(t, 400, resp.StatusCode, "Bad Request  response is expected")
	} else {
		assert.Equal(t, 200, resp.StatusCode, "OK response is expected")
	}

}
func Test_adminAPITracksCount(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(AdminHandlerGet))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}

	req, err = http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusBadRequest, resp.StatusCode)
		return
	}

}
func Test_adminAPITracksDelete(t *testing.T) {
	// Create a request to pass to our handler
	// There are no query parameters, that's why the third parameter is nil
	req, err := http.NewRequest("DELETE", "/paragliding/admin/api/tracks", nil)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder to record the response
	resRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(AdminHandlerDelete)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(resRecorder, req)

	// Check the status code
	if resRecorder.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", resRecorder.Code, http.StatusOK)
	}
}
func Test_adminAPITracks(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(AdminHandlerDelete))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the DELETE request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the DELETE request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}

	req, err = http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusNotImplemented %d, received %d. ", http.StatusBadRequest, resp.StatusCode)
		return
	}

}
func Test_getAPITickerLatest(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPITickerLatest))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}

	req, err = http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}

}
func Test_getAPITicker(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPITicker))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected StatusOK %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}

	req, err = http.NewRequest(http.MethodPost, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the POST request, %s", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Error executing the POST request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected Status OK %d, received %d. ", http.StatusOK, resp.StatusCode)
		return
	}

}
func Test_getAPITickerTimestamp(t *testing.T) {
	// instantiate mock HTTP server (just for the purpose of testing
	ts := httptest.NewServer(http.HandlerFunc(getAPITickerTimeStamp))
	defer ts.Close()

	//create a request to our mock HTTP server
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Errorf("Error constructing the GET request, %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Error executing the GET request, %s", err)
	}

	//check if the response from the handler is what we except
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected StatusBadRequest %d, received %d. ", http.StatusBadRequest, resp.StatusCode)
		return
	}

}
func Test_mongoConnect(t *testing.T) {
	if conn := connectToDB("track"); conn == nil {
		t.Error("No connection")
	}
}
func Test_mongoConnectWebhooks(t *testing.T) {
	if conn := connectToDB("webhooks"); conn == nil {
		t.Error("No connection")
	}
}


func Test_urlInMong(t *testing.T) {
	URLt := &url{}
	urlExists := checkUrl(connectToDB("track"), URLt.URL,"url")
	if urlExists !=0{
		t.Error("Track should not exist")
	}
}
func Test_webhookInMongo(t *testing.T) {
	URLt := &url{}
	urlExists := checkUrl(connectToDB("webhooks"), URLt.URL,"url")
	if urlExists !=0{
		t.Error("Webhook should not exist")
	}
}
func Test_whTrigger(t *testing.T){

	err := triggerWebhook()

	if err != nil {
		t.Error(err)
	}
}

func Test_whcTrigger(t *testing.T){

	err := triggerWebhookPeriod()

	if err != nil {
		t.Error(err)
	}
}