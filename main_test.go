package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

	testServer := httptest.NewServer(http.HandlerFunc(getApiIGC))
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



func Test_getApiIgcID_Malformed(test *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(getApiIgcID))
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

	testServer := httptest.NewServer(http.HandlerFunc(getApiIgcIDField))
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

