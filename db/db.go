package db

import (
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"

	"google.golang.org/grpc"
	"log"
	
)

type DB struct {

}

func (db *DB) NewClient() *dgo.Dgraph {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	d, err := grpc.Dial("https://blue-surf-1290226.us-east-1.aws.cloud.dgraph.io/graphql", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}



