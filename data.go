package main

import (
	"log"
	"time"

	"fmt"

	"gopkg.in/mgo.v2"
)

type primaryConnection struct {
	session *mgo.Session
}

func createConnection() primaryConnection {
	fmt.Println("dialing...")
	maxWait := time.Duration(1 * time.Second)
	session, err := mgo.DialWithTimeout("localhost:27017", maxWait)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connected to mongo")

	connection := primaryConnection{session: session}

	return connection
}

func (p *primaryConnection) getTodoCollection() TodoCollection {
	sess := p.session.Copy()
	coll := sess.DB("todomvc").C("todos")
	return TodoCollection{session: sess, collection: coll}
}

type TodoCollection struct {
	session    *mgo.Session
	collection *mgo.Collection
}
