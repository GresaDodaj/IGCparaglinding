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
