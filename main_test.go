package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSessionMidlleware(t *testing.T) {
	fmt.Println("Testing session midddleware...")
	var h http.Handler
	hh := SessionMidlleware(h)

	switch v := hh.(type) {
	case http.Handler:
		// correct handler
	default:
		t.Error(fmt.Sprintf("Type is not http.Handler, but got type:  %T", v))
	}
}

func TestMiddleWareTest(t *testing.T) {
	fmt.Println("Testing main midddleware...")

	var h http.Handler
	hh := MiddleWareTest(h)

	switch v := hh.(type) {
	case http.Handler:
		// correct handler
	default:
		t.Error(fmt.Sprintf("Type is not http.Handler, but got type:  %T", v))
	}
}

var testData = []struct {
	MethodName string
	MethodType string
	Params     string
	StatusCode int
	ServerURL  string
}{
	{"Home", "GET", "", http.StatusOK, "/"},
	{"About", "POST", "", http.StatusOK, "/About"},
}

func TestHandlers(t *testing.T) {
	fmt.Println("Testing Handlers...")
	routes := getRoutes()
	h := httptest.NewTLSServer(routes)
	for _, test := range testData {
		if test.MethodType == "GET" {
			fmt.Println(h.URL + test.ServerURL)
			// Used main server URL to test
			resp, err := h.Client().Get("http://127.0.0.1:8080" + test.ServerURL)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != test.StatusCode {
				t.Log("Status code did not match...")
				t.Fatal("Status code error...")
			}
			// success
		} else {
			values := url.Values{}
			resp, err := h.Client().PostForm("http://127.0.0.1:8080"+test.ServerURL, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != test.StatusCode {
				t.Log("Status code did not match...")
				t.Fatal("Status code error...")
			}
			// success
		}
	}

	defer h.Close()
}
