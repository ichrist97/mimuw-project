package database

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

var Session *gocql.Session

// automatically called before main() when import as package in go
func init() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "mimuwapi"
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = time.Second * 10

	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Cassandra")
}
