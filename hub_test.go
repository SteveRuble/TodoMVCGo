package main

import (
	//. "gopkg.in/ahmetb/go-linq.v3"
	//"fmt"
	"testing"
	//"time"
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
	if testChannels != nil {
		for _, c := range testChannels {
			close(c)
		}
	}
	testChannels = make(map[string]CommandChannel)
	sut = newHub()
	go sut.run(testChannelFactory)
}

func testChannelFactory(listID string) CommandChannel {
	testChannels[listID] = make(CommandChannel)
	return testChannels[listID]
}

func TestHubCreatesCommandChannels(t *testing.T) {
	setUpHubTests()
	client := NewClient(id1)
	sut.register <- client
	if _, ok := testChannels[id1]; !ok {
		t.Error("command channel not created")
	}
}
