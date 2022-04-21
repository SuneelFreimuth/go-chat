package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		registerClient(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveClient(w, r)
	})
	fmt.Println("Starting server on port 5000.")
	http.ListenAndServe("localhost:5000", nil)
}

func registerClient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Registering new client...")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading HTTP connection to Websocket: %s", err)
		return
	}

	client := Client{conn}

	client.RequestIdentification()

	messageType, p, err := conn.ReadMessage()
	if err != nil {
		fmt.Printf("Error reading message from client: %s\n", err)
	}
	fmt.Println("Received message of type", messageType, "->", p)

	if err := conn.WriteJSON(Message{"Other person", "Hello!"}); err != nil {
		fmt.Println("Error sending message to client: %s\n", err)
		return
	}
	fmt.Println("Message sent successfully!")
}

func serveClient(w http.ResponseWriter, r *http.Request) {
	filepath := fmt.Sprintf("client/%s", r.URL.Path)
	file, err := os.Open(filepath)
	if err != nil {
		// TODO: Assumes file doesn't exist, prob unsafe assumption.
		fmt.Println("Failed to find", filepath)
		http.Error(w, fmt.Sprintf("Could not find file %s", filepath), 404)
		return
	}
	fmt.Printf("Serving %s\n", filepath)
	http.ServeContent(w, r, r.URL.Path, time.Now(), file)
}

type MessageIn struct {
	Author string `json:"author"`
	Body string `json:"body"`
}

type MessageOut struct {
	Author string `json:"author"`
	Body string `json:"body"`
	// TODO: Should be a more useful type like time.Time; currently just the number of
	// milliseconds
	ServerTimestamp int64 `json:"serverTimestamp"`
}

type RoomManager struct {
	Rooms map[Room][]Client
	
}

type Room struct {
	// If nil, the Room is anonymous.
	Name string

	Id string
}

type Client struct {
	Conn *websocket.Conn
	
}

// type Hub struct {
// 	rooms map[string]Room
// }

// type Room struct {
// 	Name string
// 	Members map[string]User
// }

// type User struct {
// 	Name string
// 	Id string
// 	Client Client
// }