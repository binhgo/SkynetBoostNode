package main

import (
	"github.com/globalsign/mgo"
)

func CreateSession() (*mgo.Session, error)  {
	session, err := mgo.Dial("mongodb:27017")
	if err != nil {
		return nil, err
	}

	return session, nil
}