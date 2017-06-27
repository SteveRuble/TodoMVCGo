package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const APIPREFIX string = "/api/todos/"

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "API")
}

func main() {

	initFrontEnd()

	initApi()

	http.ListenAndServe(":8080", nil)
}

func initApi() {

	connection := createConnection()

	store := todoStore{primaryConnection: connection}

	todosHandler := makeTodosHandler(store)

	http.HandleFunc(APIPREFIX, todosHandler)
}

func initFrontEnd() {
	wd, _ := os.Getwd()
	dir := wd + "\\frontend\\"
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)
}

func makeTodosHandler(store todoStore) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		encoder := json.NewEncoder(w)

		id := req.URL.Path[len(APIPREFIX):]

		switch req.Method {
		case http.MethodGet:
			{
				todos, err := store.LoadTodos(id)
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
					err = store.SaveTodos(id, todos)
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
}
