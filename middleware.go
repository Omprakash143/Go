package main

import (
	"fmt"
	"net/http"
)

func MiddleWareTest1(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Page Hit")
		hf.ServeHTTP(w, r)
		// fmt.Println("Page served success")
	})
}
