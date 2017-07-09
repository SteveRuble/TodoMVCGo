package main

import (
	"fmt"
	"log"

	. "gopkg.in/ahmetb/go-linq.v3"
	"gopkg.in/mgo.v2/bson"
)

const ActionCreate string = "CREATE"
const ActionDelete string = "DELETE"
const ActionUpdate string = "UPDATE"

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

type TodoCommand struct {
	ListID string  `json:"listID"`
	Action string  `json:"action"`
	Todo   ApiTodo `json:"todo"`
}

type todoStore struct {
	primaryConnection primaryConnection
}

type singleListStore struct {
	primaryConnection primaryConnection
	listID            string
	commands          chan TodoCommand
	done              chan bool
}

func makeSingleListCommandChannelFactory(pc primaryConnection) CommandChannelFactory {

	return func(listID string) CommandChannel {
		var commandPipeline = make(CommandChannel)
		go run(commandPipeline, pc, listID)
		return commandPipeline
	}

}

func run(c CommandChannel, pc primaryConnection, listID string) {
	coll := pc.getTodoCollection()
	defer coll.session.Close()

	for {
		select {
		case command, ok := <-c:

			if !ok {
				// hub has closed our channel, we're done here
				return
			}

			document := TodosDocument{}
			err := coll.collection.Find(bson.M{"id": listID}).One(&document)
			notFound := err != nil
			document.ID = listID

			applyCommand(&document, command)
			// drain any other commands for this list
			n := len(c)
			for i := 0; i < n; i++ {
				applyCommand(&document, <-c)
			}

			if notFound {
				log.Printf("inserting new document with id %s\n", listID)
				err = coll.collection.Insert(document)
			} else {
				log.Printf("updating document with id %s\n", listID)
				err = coll.collection.Update(bson.M{"id": listID}, document)
			}

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func applyCommand(doc *TodosDocument, command TodoCommand) error {

	switch command.Action {
	case ActionCreate:
		doc.Todos = append(doc.Todos, command.Todo)
	case ActionDelete:
		tmp := make([]ApiTodo, len(doc.Todos))
		From(doc.Todos).Where(func(x interface{}) bool {
			return x.(ApiTodo).ID != command.Todo.ID
		}).ToSlice(&tmp)
		doc.Todos = tmp
	case ActionUpdate:
		for i, todo := range doc.Todos {
			if todo.ID == command.Todo.ID {
				todo.IsDone = command.Todo.IsDone
				todo.Title = command.Todo.Title
				doc.Todos[i] = todo
			}
		}
	}

	return nil
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
