package database

import (
	"fmt"
	"time"
	"github.com/gocql/gocql"
)

var Session *gocql.Session

func initConnection() {
	var err error 
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 10

	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Cassandra")
	defer session.Close()
}