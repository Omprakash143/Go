package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal("Application did not start...")
	}
}
func MiddleWareTest(hf http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Page Hit")
		hf.ServeHTTP(w, r)
		fmt.Println("Page served success")
	})
}

func SessionMidlleware(hf http.Handler) http.Handler {
	return session.LoadAndSave(hf)
}

func run() error {

	// sessions
	session = scs.New()
	session.Lifetime = time.Second * 30
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	muxRoutes := getRoutes()
	fmt.Println("starting server...")
	http.ListenAndServe(":8080", muxRoutes)

	return nil
}

func getRoutes() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(MiddleWareTest)
	mux.Use(SessionMidlleware)

	mux.Get("/", Home)
	mux.Post("/About", About)
	return mux
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Home"))
	ip := r.RemoteAddr
	session.Put(r.Context(), "ip", ip)
	_, _ = fmt.Println(fmt.Sprintf("Home - number of bytes is "))
}

func About(w http.ResponseWriter, r *http.Request) {

	ip := session.GetString(r.Context(), "ip")
	if len(ip) == 0 {
		w.Write([]byte("Hello About Not sure of Your ip"))
	} else {
		w.Write([]byte("Hello About Your ip is " + ip))
	}
	_, _ = fmt.Println(fmt.Sprintf("About - number of bytes is "))
}
