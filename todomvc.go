package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

const APIPREFIX string = "/api/todos/"

// wsHandler implements the Handler Interface
type wsHandler struct{}

func main() {

	handler := makeAPIHandler()

	hub := newHub()
	channelFactory := makeSingleListCommandChannelFactory(handler.store.primaryConnection)
	go hub.run(channelFactory)

	mux := http.NewServeMux()

	//ws := wsHandler{}

	//mux.Handle("/ws", ws)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("got request")
	})

	mux.Handle("/api/todos/", handler)

	mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	http.ListenAndServe("localhost:8080", mux)

}

type apiHandler struct {
	store todoStore
}

func makeAPIHandler() apiHandler {

	connection := createConnection()

	store := todoStore{primaryConnection: connection}

	handler := apiHandler{store: store}

	return handler
}

func (a apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	encoder := json.NewEncoder(w)

	id := req.URL.Path[len(APIPREFIX):]

	switch req.Method {
	case http.MethodGet:
		{
			todos, err := a.store.LoadTodos(id)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				encoder.Encode(todos)
			}
		}
	case http.MethodPost:
		{
			todos := []ApiTodo{}
			decoder := json.NewDecoder(req.Body)
			err := decoder.Decode(&todos)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				err = a.store.SaveTodos(id, todos)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusAccepted)
				}
			}
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got ws request")
	// upgrader is needed to upgrade the HTTP Connection to a websocket Connection
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	//Upgrading HTTP Connection to websocket connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("error upgrading %s", err)
		return
	}
	_ = wsConn
	//handle your websockets with wsConn
}
