package raft

import (
	"log"
	"testing"
)

func logMsg(data interface{}) {
	switch msg := data.(type) {
	case *ElectMessage:
		log.Printf("[MSG][Elect] leader = %s term = %d\n", msg.Server, msg.Term)
	case *CommitMessage:
		log.Printf("[MSG][Commit] server = %s term = %d commitIndex = %d\n", msg.Server, msg.Term, msg.CommitIndex)
	default:
		log.Printf("unknown message")
	}
}

func Test_Election(t *testing.T) {
	observeChan := make(chan interface{}, 16)
	peers := make([]string, 3)
	peers[0] = "server1"
	peers[1] = "server2"
	peers[2] = "server3"
	c := NewCluster(peers, observeChan)
	c.Start()
	msg := <-observeChan
	electMsg := msg.(*ElectMessage)
	logMsg(electMsg)
	c.Stop()
}

func Test_ElectionAfterCrash(t *testing.T) {
	observeChan := make(chan interface{}, 16)
	peers := make([]string, 3)
	peers[0] = "server1"
	peers[1] = "server2"
	peers[2] = "server3"
	c := NewCluster(peers, observeChan)
	c.Start()
	msg := <-observeChan
	electMsg := msg.(*ElectMessage)
	logMsg(electMsg)

	server := c.Get(electMsg.Server)
	server.Stop()

	msg = <-observeChan
	electMsg = msg.(*ElectMessage)
	logMsg(electMsg)

	c.Stop()
}

func Test_ElectionHighTerm(t *testing.T) {
	observeChan := make(chan interface{}, 16)
	peers := make([]string, 3)
	peers[0] = "server1"
	peers[1] = "server2"
	peers[2] = "server3"
	c := NewCluster(peers, observeChan)

	server1 := c.Get("server1")
	server1.log.append(&LogEntry{1, 1, "1"})
	server1.raft.CurrentTerm = 2
	c.Start()
	msg := <-observeChan
	logMsg(msg)
	c.Stop()
}

func Test_AppendSync(t *testing.T) {
	observeChan := make(chan interface{}, 16)
	peers := make([]string, 3)
	peers[0] = "server1"
	peers[1] = "server2"
	peers[2] = "server3"
	c := NewCluster(peers, observeChan)
	c.Start()
	msg := <-observeChan
	logMsg(msg)
	leader := c.GetLeader()
	leader.sendCommandAsync("get")
	leader.sendCommandAsync("set")
	leader.sendCommandAsync("add")
	for i := 0; i < 3; i++ {
		msg = <-observeChan
		logMsg(msg)
	}
	c.Stop()
}
