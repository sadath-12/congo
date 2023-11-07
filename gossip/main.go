package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Node struct {
	ID        int
	Data      string
	Peers     map[int]*Node
	gossipLog []string
}

func (n *Node) Gossip(ch chan struct{}) {
	for peerID, peer := range n.Peers {
		go func(peerID int, peer *Node) {
			if err := n.sendDataToPeer(peerID, peer); err != nil {
				log.Printf("Node %d encountered an error sending data to Node %d: %v\n", n.ID, peerID, err)
				n.notifyPeersAboutUnresponsiveNode(peerID)
			}
			ch <- struct{}{}
		}(peerID, peer)
	}
}

func (n *Node) sendDataToPeer(peerID int, peer *Node) error {
	fmt.Printf("Node %d sent data to Node %d: %s\n", n.ID, peerID, n.Data)
	peer.ReceiveData(n.Data)
	n.gossipLog = append(n.gossipLog, fmt.Sprintf("Sent to Node %d", peerID))
	return nil
}

func (n *Node) ReceiveData(data string) {
	n.Data = data
	n.gossipLog = append(n.gossipLog, "Received")
}

func (n *Node) notifyPeersAboutUnresponsiveNode(unresponsiveNodeID int) {
	for peerID, peer := range n.Peers {
		if peerID != unresponsiveNodeID {
			fmt.Printf("Node %d notifying Node %d about unresponsive Node %d\n", n.ID, peerID, unresponsiveNodeID)
			peer.HandleUnresponsiveNode(unresponsiveNodeID)
		}
	}
}

func (n *Node) HandleUnresponsiveNode(unresponsiveNodeID int) {
	fmt.Printf("Node %d received notification about unresponsive Node %d\n", n.ID, unresponsiveNodeID)
}

func (n *Node) writeLogToFile() {
	logsFolder := "logs"
	if _, err := os.Stat(logsFolder); os.IsNotExist(err) {
		if err := os.Mkdir(logsFolder, os.ModePerm); err != nil {
			log.Fatalf("Error creating logs folder: %v", err)
		}
	}

	logFileName := filepath.Join(logsFolder, fmt.Sprintf("node%d.log", n.ID))
	file, err := os.Create(logFileName)
	if err != nil {
		log.Fatalf("Error creating log file for Node %d: %v", n.ID, err)
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)

	for _, entry := range n.gossipLog {
		logger.Println(entry)
	}
}

func main() {
	numNodes := 5

	nodes := make([]*Node, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i] = &Node{
			ID:        i + 1,
			Data:      fmt.Sprintf("Initial data from Node %d", i+1),
			Peers:     make(map[int]*Node, numNodes-1),
			gossipLog: make([]string, 0),
		}
	}

	for i, node := range nodes {
		for j, peer := range nodes {
			if i != j {
				node.Peers[j+1] = peer
			}
		}
	}

	duration := 10 * time.Second
	gossipInterval := 1 * time.Second
	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		ch := make(chan struct{}, numNodes-1)
		for _, node := range nodes {
			go node.Gossip(ch)
		}
		for i := 0; i < numNodes-1; i++ {
			<-ch
		}
		time.Sleep(gossipInterval)
	}

	for _, node := range nodes {
		node.writeLogToFile()
	}
}
