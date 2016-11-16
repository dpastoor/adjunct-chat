package main

import (
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("client connected!")
		so.Emit("chat", "hello from server")
		so.Join("edits")
		// I don't fully understand the behavior, but it seems like
		// the only way to send to both is to emit, when the event was not triggered
		// by a given client, else emitting will only send back
		so.On("edit", func(msg string) {
			log.Println("recieved message", msg)
			//so.Emit("chat", msg) // this effectively echos back the message only to the sender in the chat channel
			// emulate a delay in network
			so.BroadcastTo("edits", "update", msg) // sends update to all OTHER clients
		})
		so.On("edit:knit", func(msg string) {
			log.Println("recieved knit message", msg)
			//so.Emit("chat", msg) // this effectively echos back the message only to the sender in the chat channel
			so.BroadcastTo("edits", "update:knit", msg) // sends update to all OTHER clients
		})
		so.On("disconnection", func() {
			log.Println("disconnected from chat")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	mux.Handle("/socket.io/", server)
	mux.Handle("/assets", http.FileServer(http.Dir("./assets")))

	handler := cors.Default().Handler(mux)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	// Insert the middleware
	handler = c.Handler(handler)

	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", handler))
}

//'Access-Control-Allow-Credentials' must be true
