package main

import (
	"log"
	"net"
	"time"

	"github.com/globalsign/mgo"
)

var DeleteQueue chan *Node
var InsertQueue chan *Node
var UpdateQueue chan *Node
var GetQueue chan net.Conn

func main() {

	DeleteQueue = make(chan *Node, 20)
	InsertQueue = make(chan *Node, 20)
	UpdateQueue = make(chan *Node, 20)
	GetQueue = make(chan net.Conn, 20)

	session, err := CreateSession()
	defer session.Close()

	if err != nil {
		log.Fatalln("Cannot create mgo session. Panic")
	}

	InitServer(":6868", InsertQueue, UpdateQueue, GetQueue)
	go InsertNewPeerFromQueue(session)

	go CheckPeerStatus(session)
	go RemoveInActivePeers(session)

	// handle request
}

func InsertNode(node *Node, session *mgo.Session) (*Node, error) {
	repo := GetRepo(session)
	defer repo.Close()

	n, err := repo.Insert(node)
	if err != nil {
		log.Printf("Error INSERT: %s", err)
		return nil, err
	}

	return n, nil
}

func GetAllNodes(session *mgo.Session) ([]*Node, error) {
	repo := GetRepo(session)
	defer repo.Close()

	nodes, err := repo.GetAll()
	if err != nil {
		log.Printf("Error GETALL: %s", err)
		return nil, err
	}

	return nodes, nil
}

func DeleteNode(session *mgo.Session, node *Node) error {
	repo := GetRepo(session)
	defer repo.Close()

	err := repo.Delete(node)
	if err != nil {
		log.Printf("Error DeleteNode: %s", err)
		return err
	}

	return nil
}

func CheckPeerStatus(session *mgo.Session) {
	for {
		time.Sleep(time.Second * 60)

		nodes, err := GetAllNodes(session)
		if err != nil {
			log.Printf("Error CheckPeerStatus.\n Time: %s\n ERR: %s\n", time.Now().String(), err)
		}

		for _, n := range nodes {
			if n.UpdatedTime.Add(time.Hour * 24).Before(time.Now()) {
				// 24 hours pass seen last update
				// save into channel -> another go routine read from channel and delete later
				DeleteQueue <- n
			}
		}
	}
}

func RemoveInActivePeers(session *mgo.Session) {
	for {
		node := <-DeleteQueue
		err := DeleteNode(session, node)
		if err != nil {
			log.Printf("Error RemoveInActivePeers.\n Time: %s\n ERR: %s\n", time.Now().String(), err)
			continue
		}
	}
}

func InsertNewPeerFromQueue(session *mgo.Session) {
	for {
		node := <-InsertQueue
		n, err := InsertNode(node, session)
		if err != nil {
			log.Printf("Error RemoveInActivePeers.\n Time: %s\n ERR: %s\n", time.Now().String(), err)
			continue
		}

		log.Printf("Inserted new peer: %s\n", n.MongoID)
	}
}
