package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Home"))

	_, _ = fmt.Println(fmt.Sprintf("Home - number of bytes is "))
}

func About(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello About"))
	_, _ = fmt.Println(fmt.Sprintf("About - number of bytes is "))
}
func main() {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(MiddleWareTest)
	mux.Get("/", Home)
	mux.Get("/About", About)

	fmt.Println("starting server...")
	http.ListenAndServe(":8080", mux)

}
func MiddleWareTest(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Page Hit")
		hf.ServeHTTP(w, r)
		fmt.Println("Page served success")
	})
}
