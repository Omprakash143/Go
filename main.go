package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
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

	hub := newHub()
	go hub.run()
	mux.Get("/", Home)
	mux.Post("/About", About)
	mux.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsAccess(hub, w, r)
	})
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

// web socket ------

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Endpoint to serve serving web-socket requests...
func wsAccess(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTp/HTTPs connection to web-socket connection
	Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to created web-socket connection from HTTPs")
		return
	}
	fmt.Println("Web-socket connected...")

	client := &Client{
		hub,
		ws,
		make(chan []byte, 256),
	}
	hub.Register <- client
	fmt.Println("Client registered...")

	go client.cliWS()
	go client.serverWS()

}

// creating HUB to broadcast messages to clients

type Hub struct {
	Clients map[*Client]bool

	Register chan *Client

	Unregister chan *Client

	Broadcast chan []byte
}

func newHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// creating clients

type Client struct {
	HubClient *Hub
	WsConn    *websocket.Conn
	Buffer    chan []byte
}

func (hub *Hub) run() {
	fmt.Println("Hub started........")
	for {
		select {
		case c := <-hub.Register:
			hub.Clients[c] = true
		case c1 := <-hub.Unregister:
			delete(hub.Clients, c1)
			close(c1.Buffer)
		case msg := <-hub.Broadcast:
			for cli, _ := range hub.Clients {
				select {
				case cli.Buffer <- msg:
				default:
					close(cli.Buffer)
					delete(hub.Clients, cli)
				}
			}
		}
	}
}

func (cli *Client) cliWS() {
	fmt.Println("Client started........")
	defer func() {
		cli.HubClient.Unregister <- cli
		cli.WsConn.Close()
	}()
	for {
		_, message, err := cli.WsConn.ReadMessage()
		if err != nil {
			fmt.Println("Failed to read data from Client")
			cli.WsConn.Close()
			break
		}
		cli.HubClient.Broadcast <- message
	}
}

func (cli *Client) serverWS() {
	fmt.Println("server started........")
	defer func() {
		cli.WsConn.Close()
	}()
	for {
		select {
		case message, ok := <-cli.Buffer:
			if !ok {
				cli.WsConn.WriteMessage(websocket.CloseMessage, []byte("Hub closed the channel..."))
				return
			}
			err := cli.WsConn.WriteMessage(1, message)
			if err != nil {
				fmt.Println("Failed to write message to Client...")
				break
			}
		}
	}
}
