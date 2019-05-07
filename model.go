package main

import (
	"net"
	"time"

	"github.com/globalsign/mgo/bson"
)

type Node struct {
	MongoID      bson.ObjectId `bson:"_id,omitempty"`
	IP           string
	Port         string
	ServerPort   string
	RegisterTime time.Time
	UpdatedTime  time.Time
	IsActive     bool

	// connection to send data back to client
	Conn net.Conn
}

func NewNode(ip string, port string, serverPort string, registerTime time.Time, updatedTime time.Time, isActive bool) *Node {
	node := &Node{nil, ip, port, serverPort, registerTime, updatedTime, isActive, nil}
	return node
}