package db

import (
	"github.com/dgraph-io/dgo/v210"
	"google.golang.org/grpc"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masud2017/social_media_golang/models"
)

type DB struct {
	Client *dgo.Dgraph
	SchemaOp *api.Operation
	ctx context.Context

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
		email: string @index(exact,term) .
		password: string @index(term) .
		friend: [uid] .
		rel: string @index(exact) .
		user : uid .
		
		req: string @index(exact,term) . 
		req_to: uid .
		request: [uid] .

		req_rel: string @index(exact,term) .
		request_from: [uid] .
		req_from: uid .
		
		type User {
			name: string
			email: string
			password: string
			friend: [Relation]
			request: [RelationRequest]
			request_from:[RelationRequestFromOther]
		}
		type Relation {
			rel: string
			user: User
		}	
		type RelationRequest {
			req_rel: string
			req_to: User
		}			
		type RelationRequestFromOther {
			req_rel: string
			req_from: User
		}
	`

	db.SchemaOp = op

	ctx := context.Background()
	if err := db.Client.Alter(ctx, db.SchemaOp); err != nil {
		log.Fatal(err)
	}

	db.ctx = ctx
}

func isEmailUnique(email string,ctx context.Context,client *dgo.Dgraph) bool {
	query := `
	{
		findUserByEmail(func: allofterms(email,"`+email+`")) {
			uid
			name
			email
			password
		}
	}
	`

	resp, err := client.NewTxn().Query(ctx,query)
	if err != nil {
		log.Fatal(err)
	}

	
	type Root struct {
		FindUserByEmail []models.User `json:"findUserByEmail,omitempty"`
	}
	var root Root
	
	if err := json.Unmarshal(resp.GetJson(), &root); err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
	if (len(root.FindUserByEmail) == 0) {
		return false
	} else {
		if (root.FindUserByEmail[0].Email == email) {
			fmt.Println("This user does exist ..")
			
			return true;
		}
	}

	return false;
}

func (db *DB) SignupUser(user *models.User) bool {
	db.InitSchema()

	if (isEmailUnique(user.Email,db.ctx,db.Client)) {
		return false
	}



	mu := &api.Mutation{
		CommitNow: true,
	}
	pb, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb
	response, err := db.Client.NewTxn().Mutate(db.ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)

	return true

}

func (db *DB) GetUserList() []models.User {
	db.InitSchema()
	
	query := `
	{
		getAllUsers(func: has(name)) {
			uid
			name
			email
			password
		}
	}
	`

	resp, err := db.Client.NewTxn().Query(db.ctx,query)
	if err != nil {
		log.Fatal(err)
	}

	
	type Root struct {
		GetAllUsers []models.User `json:"getAllUsers,omitempty"`
	}
	var root Root
	
	if err := json.Unmarshal(resp.GetJson(), &root); err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
	

	return root.GetAllUsers
}


func (db *DB) Me(uid string) models.User {
	db.InitSchema()
	
	query := `
	{
		getMe(func: uid(`+uid+`)) {
			uid
			name
			email
			password
		}
	}
	`

	resp, err := db.Client.NewTxn().Query(db.ctx,query)
	if err != nil {
		log.Fatal(err)
	}

	
	type Root struct {
		GetMe []models.User `json:"getMe,omitempty"`
	}
	var root Root
	
	if err := json.Unmarshal(resp.GetJson(), &root); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Printing the value of resp :")

	fmt.Println(resp)


	return root.GetMe[0]
}

func (db *DB) RequestForRelationship(relReq models.RelationRequest,me models.User) models.RelationRequest {
	db.InitSchema()

	mu := &api.Mutation{
		CommitNow: true,
	}

	// me.Request = relReq
	me.Request = append(me.Request, relReq)
	

	pb, err := json.Marshal(me)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb

	fmt.Println("Printing the value of pb : ",pb)
	response, err := db.Client.NewTxn().Mutate(db.ctx, mu)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.GetJson())

	// after that need to populate other user info too

	reqFrom := models.RelationRequestFromOther{
		Uid : "_:req_from",
		ReqRel : relReq.ReqRel,
		ReqFrom : models.User{
			Uid : me.Uid,
			Name: me.Name,
			Email: me.Email,
		},
	}
	reqToUser := relReq.ReqTo
	reqToUser.RequestFrom = append(reqToUser.RequestFrom,reqFrom)

	pb2, err := json.Marshal(reqToUser)
	if err != nil {
		log.Fatal(err)
	}

	mu.SetJson = pb2
	
	db.Client.NewTxn().Mutate(db.ctx, mu)
	

	return relReq
}


/*
This function will return all the relation ship requests from other users 
*/
func (db *DB) RelationShipRequests(user_id string) []models.RelationRequestFromOther {
	db.InitSchema()
	
	query := `
	{
		getRelationShipRequestFromOther(func: uid(`+user_id+`)) {
			request_from {
				uid
				req_rel
				req_from {
					uid
					name
					email
				}
			}
		}
	}
	`

	resp, err := db.Client.NewTxn().Query(db.ctx,query)
	if err != nil {
		log.Fatal(err)
	}

	
	type Root struct {
		GetRelationShipRequestFromOther []models.User `json:"getRelationShipRequestFromOther,omitempty"`
	}
	var root Root
	
	if err := json.Unmarshal(resp.GetJson(), &root); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Printing the value of resp :")

	fmt.Println(resp)

	

	return root.GetRelationShipRequestFromOther[0].RequestFrom
}

/*
This function will return all the relationship requests that this user made to other users 
*/
func (db *DB) MyRelationShipRequests(user_id string) []models.RelationRequest {
	db.InitSchema()
	
	query := `
	{
		getRelationShipRequest(func: uid(`+user_id+`)) {
			request {
				uid
				req_rel
				req_to {
					uid
					name
					email
				}
			}
		}
	}
	`

	resp, err := db.Client.NewTxn().Query(db.ctx,query)
	if err != nil {
		log.Fatal(err)
	}

	
	type Root struct {
		GetRelationShipRequest []models.User `json:"getRelationShipRequest,omitempty"`
	}
	var root Root
	
	if err := json.Unmarshal(resp.GetJson(), &root); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Printing the value of resp :")

	fmt.Println(resp)

	

	return root.GetRelationShipRequest[0].Request
}


func (db *DB) MyRelationList(user_id string) []models.Relation {
	db.InitSchema()
	
	query := `
	{
		getRelations(func: uid(`+user_id+`)) {
			friend {
				uid
				rel
				user {
					uid
					name
					email
				}
			}
		}
	}
	`

	resp, err := db.Client.NewTxn().Query(db.ctx,query)
	if err != nil {
		log.Fatal(err)
	}

	
	type Root struct {
		GetRelations []models.User `json:"getRelations,omitempty"`
	}
	var root Root
	
	if err := json.Unmarshal(resp.GetJson(), &root); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Printing the value of resp :")

	fmt.Println(resp)

	

	return root.GetRelations[0].Friend
}

func (db *DB) AcceptReq() {
	
}


// func (db *DB) InitSchema(user *models.User) {
// 	op := &api.Operation{}

// 	p:= user

// 	op.Schema = `
// 		name: string @index(exact) .
// 		email: string @index(exact,term) .
// 		password: string @index(term) .
// 		Friend: [uid] .
// 		rel: string @index(exact) .
// 		user : uid .
// 		type User {
// 			name: string
// 			email: string
// 			password: string
// 			Friend: [Relation]
// 		}
// 		type Relation {
// 			rel: string
// 			user: User
// 		}				
// 	`

// 	ctx := context.Background()
// 	if err := db.Client.Alter(ctx, op); err != nil {
// 		log.Fatal(err)
// 	}


// 	mu := &api.Mutation{
// 		CommitNow: true,
// 	}
// 	pb, err := json.Marshal(p)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	mu.SetJson = pb
// 	response, err := db.Client.NewTxn().Mutate(ctx, mu)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Assigned uids for nodes which were created would be returned in the response.Uids map.
// 	variables := map[string]string{"$id1": response.Uids["masud"]}
// 	q := `query Me($id1: string){
// 		me(func: uid($id1)) {
// 			name
// 			email
// 			password
// 			friend @filter(eq(name, "Md")){
// 				rel
// 				user_id
// 			}
// 		}
// 	}`

// 	resp, err := db.Client.NewTxn().QueryWithVars(ctx, q, variables)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	type Root struct {
// 		Me []models.User `json:"me"`
// 	}

// 	var r Root
// 	err = json.Unmarshal(resp.Json, &r)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	out, _ := json.MarshalIndent(r, "", "\t")
// 	fmt.Printf("%s\n", out)

// }




