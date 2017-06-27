package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"
)

// TodosDocument is the model of of the document persisted to Mongo
type TodosDocument struct {
	ID    string
	Todos []ApiTodo
}

// ApiTodo is the model of todo stored in Mongo and used by client
type ApiTodo struct {
	Title  string `json:"title"`
	ID     string `json:"id"`
	IsDone bool   `json:"completed"`
}

type todoStore struct {
	primaryConnection primaryConnection
}

func (store *todoStore) SaveTodos(id string, todos []ApiTodo) error {

	coll := store.primaryConnection.getTodoCollection()
	defer coll.session.Close()

	document := TodosDocument{}
	err := coll.collection.Find(bson.M{"id": id}).One(&document)

	notFound := err != nil

	document.ID = id
	document.Todos = todos

	if notFound {
		log.Printf("inserting new document with id %s\n", id)
		err = coll.collection.Insert(document)
	} else {
		log.Printf("updating document with id %s\n", id)
		err = coll.collection.Update(bson.M{"id": id}, document)
	}

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func (store *todoStore) LoadTodos(id string) (todos []ApiTodo, err error) {

	coll := store.primaryConnection.getTodoCollection()
	defer coll.session.Close()

	document := TodosDocument{}

	err = coll.collection.Find(bson.M{"id": id}).One(&document)

	if err != nil {
		// TOOD: better way to check error meaning
		if fmt.Sprint(err) == "not found" {
			// id doesn't exist yet, so return empty todos
			return make([]ApiTodo, 0), nil
		}
		fmt.Printf("error retrieving todos with id '%s'", id)
		log.Fatal(err)
	}

	return document.Todos, err
}
