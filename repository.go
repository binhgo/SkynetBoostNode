package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type NodeRepository struct {
	mgoSession *mgo.Session
}

func GetRepo(mainMgoSession *mgo.Session) *NodeRepository {
	repo := &NodeRepository{mainMgoSession.Clone()}
	return repo
}

func (repo *NodeRepository) Insert(node *Node) (*Node, error) {
	node.MongoID = bson.NewObjectId()
	err := repo.Collection().Insert(node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (repo *NodeRepository) GetAll() ([]*Node, error) {
	var nodes []*Node
	err := repo.Collection().Find(nil).All(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (repo *NodeRepository) Update(node *Node) error {
	err := repo.Collection().UpdateId(node.MongoID, node)
	if err != nil {
		return err
	}
	return nil
}

func (repo *NodeRepository) Delete(node *Node) error {
	err := repo.Collection().RemoveId(node.MongoID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *NodeRepository) Collection() *mgo.Collection {
	return repo.mgoSession.DB("SkynetPeers").C("Node")
}

func (repo *NodeRepository) Close() {
	repo.mgoSession.Close()
}
