package db

import (
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"log"
	"context"
	
	"fmt"
	"time"
)

type DB struct {
	Client *dgo.Dgraph

}

func (db *DB) NewClient() {
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	db.Client =  dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func (db *DB) InitSchema() {
	op := &api.Operation{}

	op.Schema = `
				name: string @index(exact) .
				email: string @index(exact,term).
				Friend: [uid] .
				type User {
					name: string
					email: string
					password: string
					Friend: [Relation]
				}
				type Relation {
					relation_name: string
					user: User
				}
				
			`
	
	

}




