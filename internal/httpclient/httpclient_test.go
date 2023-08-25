package httpclient

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customHeader := r.Header["X-Custom-Header"]
		if len(customHeader) != 1 {
			w.WriteHeader(http.StatusBadRequest)
		} else if customHeader[0] != "Hello, World!" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	input := fmt.Sprintf(`GET "%s" X-Custom-Header: "Hello, World!"`, server.URL)
	hr := Hit(input)

	if hr.Err != nil {
		t.Fail()
	} else if hr.ResponseHeaders[0] != "200 OK" {
		t.Fail()
	}
}

func TestBinaryResponseBody(t *testing.T) {
	nonPrintableResponse := "\n\nRESPONSE CONTAINS NON-PRINTABLE CHARACTERS.\n"

	var tests = []struct {
		name  string
		input string
	}{
		{
			"Bytes from a JPG file",
			"/9j/4AAQSkZJRgABAQAAAQABAAD//gA7Q1JFQVRPUjogZ2QtanBlZyB2MS4wICh1c2luZyBJSkcgSlBFRyB2NjIpLCBxdWFsaXR5ID0gOTQK/9sAQwACAQECAQECAgICAgICAgMFAwMDAwMGBAQDBQcGBwcHBgcHCAkLCQgICggHBwoNCgoLDAwMDAc",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputBytes, err := base64.RawStdEncoding.DecodeString(test.input)
			if err != nil {
				panic(err)
			}
			if formatResponseBody(inputBytes) != nonPrintableResponse {
				t.Fail()
			}
		})
	}
}

func TestPlainResponseBody(t *testing.T) {
	nonPrintableResponse := "\n\nRESPONSE CONTAINS NON-PRINTABLE CHARACTERS.\n"

	var tests = []struct {
		name  string
		input string
	}{
		{
			"Text",
			"Hello, World!",
		},
		{
			"HTML",
			`<!DOCTYPE html><html lang="en"> <head> <meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>My Website</title> <link rel="stylesheet" href="./style.css">
				<link rel="icon" href="./favicon.ico" type="image/x-icon"> </head>
				<body> <main> <h1>Welcome to My Website</h1>
				</main><script src="index.js"></script> </body></html>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if formatResponseBody([]byte(test.input)) == nonPrintableResponse {
				t.Fail()
			}
		})
	}
}

func TestRedirects(t *testing.T) {
	var tests = []struct {
		name             string
		input            string
		expectedResponse string
	}{
		{
			"No Redirects",
			`GET "http://www.ramitmittal.com"`,
			"301 Moved Permanently",
		},
		{
			"With Redirects",
			`GET "http://www.ramitmittal.com" -location`,
			"200 OK",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if hr := Hit(test.input); hr == nil || hr.Err != nil {
				t.Fail()
			} else if hr.ResponseHeaders[0] != test.expectedResponse {
				t.Log(hr.ResponseHeaders[0])
				t.Fail()
			}
		})
	}
}
