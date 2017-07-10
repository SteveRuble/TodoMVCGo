package main

import (
	//. "gopkg.in/ahmetb/go-linq.v3"
	//"fmt"
	"encoding/json"
	"testing"
	"time"
)

const (
	id1 = "TEST"
)

var (
	testChannels map[string]CommandChannel
	sut          *Hub
)

func setUpHubTests() {
	if sut != nil {
		sut.done <- true
	}
	testChannels = make(map[string]CommandChannel)
	sut = newHub()
	go sut.run(testChannelFactory)
}

func testChannelFactory(listID string) CommandChannel {
	testChannels[listID] = make(CommandChannel, 100)
	return testChannels[listID]
}

func TestHubCreatesCommandChannels(t *testing.T) {
	setUpHubTests()
	client := NewClient(id1)
	sut.register <- client
	yield()
	if _, ok := testChannels[id1]; !ok {
		t.Error("command channel not created")
	}
}

func TestHubClosesCommandChannelsWhenAllClientsUnregistered(t *testing.T) {
	setUpHubTests()
	client1 := NewClient(id1)
	client2 := NewClient(id1)
	sut.register <- client1
	sut.register <- client2

	if len(testChannels) != 1 {
		t.Errorf("testChannels: %v", testChannels)
		t.Error("command channel should only be created once per ID")
	}

	sut.unregister <- client1
	sut.unregister <- client2

	if _, ok := <-testChannels[id1]; ok {
		t.Error("command channel not closed")
	}
}

func TestHubForwardsCommandsToCommandChannel(t *testing.T) {
	setUpHubTests()
	client := NewClient(id1)

	sut.register <- client

	commandChannel := testChannels[id1]

	expected := TodoCommand{Action: "ok", ListID: id1}

	bytes, _ := json.Marshal(expected)

	sut.broadcast <- bytes

	timeout := time.After(time.Millisecond)

	for {
		select {
		case _ = <-commandChannel:
			// command was forwarded
			return
		case _ = <-timeout:
			t.Error("command was not forwarded")
			return
		}
	}
}

func yield() {
	_ = <-time.After(time.Millisecond)
}
