package db

import (
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

// InitSessionConnect initialize db session  for using later
func InitSessionConnect(hostAndPort string) (*mgo.Session, error) {
	session, err := mgo.Dial(hostAndPort)
	session.SetMode(mgo.Monotonic, true)

	if err != nil {
		panic(err)
	}
	return session, nil
}

// Close will shutdown session connect
func Close() {
	session.Close()
}

// GetDBSession get one connection from pool
func GetDBSession() *mgo.Session {
	return session
}
