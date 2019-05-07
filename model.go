package main

import (
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
}

func NewNode(ip string, port string, serverPort string, registerTime time.Time, updatedTime time.Time, isActive bool) *Node {
	node := &Node{bson.NewObjectId(), ip, port, serverPort, registerTime, updatedTime, isActive}
	return node
}
