package main

import (
	//. "gopkg.in/ahmetb/go-linq.v3"
	//"fmt"
	"testing"
	//"time"
)

func TestCreateCommand(t *testing.T) {
	doc, command := initInputs()
	command.Action = ActionCreate
	command.Todo.ID = "id3"

	applyCommand(&doc, command)

	if len(doc.Todos) != 3 {
		t.Error("create did not insert todo")
	}
	if todo := doc.Todos[2]; todo.ID != command.Todo.ID {
		t.Error("create did not set ID on inserted todo")
	}
}

func TestUpdateCommand(t *testing.T) {
	doc, command := initInputs()
	command.Action = ActionUpdate
	command.Todo.IsDone = true
	command.Todo.Title = "todo1update"

	applyCommand(&doc, command)
	actual := doc.Todos[0]

	if !actual.IsDone {
		t.Error("update did not set IsDone")
	}
	if actual.Title != "todo1update" {
		t.Error("update did not set title on inserted todo")
	}
}

func TestDeleteCommand(t *testing.T) {
	doc, command := initInputs()
	command.Action = ActionDelete
	applyCommand(&doc, command)

	if len(doc.Todos) != 1 {
		t.Error("delete did not delete todo")
	}
	if todo := doc.Todos[0]; todo.ID != "id2" {
		t.Error("delete deleted the wrong todo")
	}
}

func initInputs() (doc TodosDocument, command TodoCommand) {
	doc = TodosDocument{
		Todos: []ApiTodo{
			ApiTodo{
				ID:    "id1",
				Title: "todo1",
			},
			ApiTodo{
				ID:    "id2",
				Title: "todo2",
			},
		},
	}
	command = TodoCommand{
		Action: ActionDelete,
		Todo: ApiTodo{
			ID: "id1",
		},
	}
	return
}
