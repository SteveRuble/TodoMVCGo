package main

import (
	"encoding/json"
	"net/http"
)

const APIPREFIX string = "/api/todos/"

func main() {

	hub := newHub()
	go hub.run()

	handler := makeAPIHandler()

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	mux.Handle("/api/todos/", handler)

	mux.Handle("/", http.FileServer(http.Dir("./frontend")))

	http.ListenAndServe(":8080", mux)

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
