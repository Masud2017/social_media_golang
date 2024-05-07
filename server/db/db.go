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
	d, err := grpc.Dial("database:9080", grpc.WithInsecure())
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
		father: [uid] .
		mother: [uid] .
		son: [uid] .
		rel: string @index(exact) .
		user : uid .
		
		req: string @index(exact,term) . 
		req_to: uid .
		request: [uid] .

		req_rel: string @index(exact,term) .
		request_from: [uid] .
		req_from: uid .

		req_from_uid: string @index(exact,term) . 
		req_to_uid : string @index(exact,term) .
		
		type User {
			name: string
			email: string
			password: string
			friend: [Relation]
			father: [Relation]
			mother: [Relation]
			son: [Relation]
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
			req_from_uid: string
		}			
		type RelationRequestFromOther {
			req_rel: string
			req_from: User
			req_to_uid: string
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
			request_from {
				uid
			}
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

	

	if (len(root.GetRelationShipRequestFromOther) > 0) {
		return root.GetRelationShipRequestFromOther[0].RequestFrom
	} else {
		return []models.RelationRequestFromOther{}
	}
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

	

	if (len(root.GetRelationShipRequest) > 0) {
		return root.GetRelationShipRequest[0].Request
	} else {
		return []models.RelationRequest{}
	}
}


func (db *DB) MyRelationList(user_id string) ([]models.Relation, []models.Relation, []models.Relation, []models.Relation) {
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

			father {
				uid
				rel
				user {
					uid
					name
					email
				}
			}

			mother {
				uid
				rel
				user {
					uid
					name
					email
				}
			}

			son {
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

	

	if (len(root.GetRelations) > 0) {
		return root.GetRelations[0].Friend, root.GetRelations[0].Father, root.GetRelations[0].Mother, root.GetRelations[0].Son
	} else {
		return []models.Relation{},[]models.Relation{},[]models.Relation{},[]models.Relation{}
	}
}

/*Accepts request from other user*/
/*
todo : 
-> Add the request to the user.friend/ user.mother/ user.father/ user.son
-> remove the req from both current user and the to user
*/

func delete_at_index(slice []models.RelationRequest, index int) []models.RelationRequest {
    return append(slice[:index], slice[index+1:]...)
}

func delete_at_index_relation_from(slice []models.RelationRequestFromOther, index int) []models.RelationRequestFromOther {
    return append(slice[:index], slice[index+1:]...)
}

func (db *DB) AcceptReq(user_id string,req_id string) bool {
	db.InitSchema()
	mu := &api.Mutation{
		CommitNow: true,
	}
	user := db.Me(user_id)

	query := `
	{
		getRequestFromOtherUser(func: uid(`+req_id+`)) {
			uid
			req_rel
			req_from {
					uid
			name
			email

			request {
				uid
				req_to {
					uid
				}
			}
			request_from {
				uid
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
		GetRequestFromOtherUser []models.RelationRequestFromOther `json:"getRequestFromOtherUser,omitempty"`
	}
	var root Root
	
	if err := json.Unmarshal(resp.GetJson(), &root); err != nil {
		log.Fatal(err)
	}

	if (len(root.GetRequestFromOtherUser) > 0) {
		requestFromOtherUser := root.GetRequestFromOtherUser[0]

		if (requestFromOtherUser.ReqRel == "Friend") {
			// user.Friend = append(user.Friend, models.Relation{
			// 	Rel: requestFromOtherUser.ReqRel,
			// 	User: requestFromOtherUser.ReqFrom,
			// })
			newUser:= models.User{
				Uid : user.Uid,
				Name: user.Name,
				Email: user.Email,
				Password: user.Password,
				Friend : []models.Relation{{
					Rel: requestFromOtherUser.ReqRel,
					User: requestFromOtherUser.ReqFrom,
					},
				},
			}
			
			pb, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}

			mu.SetJson = pb
			
			_ , errMut := db.Client.NewTxn().Mutate(db.ctx, mu)

			if errMut != nil {
				log.Fatal(errMut)
			}


		} else if (requestFromOtherUser.ReqRel == "Father") {
			// user.Father = append(user.Father, models.Relation{
			// 	Rel: requestFromOtherUser.ReqRel,
			// 	User: requestFromOtherUser.ReqFrom,
			// })

			newUser:= models.User{
				Uid : user.Uid,
				Name: user.Name,
				Email: user.Email,
				Password: user.Password,
				Father : []models.Relation{{
					Rel: requestFromOtherUser.ReqRel,
					User: requestFromOtherUser.ReqFrom,
					},
				},
			}

			pb, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}

			mu.SetJson = pb
			
			_ , errMut := db.Client.NewTxn().Mutate(db.ctx, mu)

			if errMut != nil {
				log.Fatal(errMut)
			}


		} else if (requestFromOtherUser.ReqRel == "Mother") {
			// user.Mother = append(user.Mother, models.Relation{
			// 	Rel: requestFromOtherUser.ReqRel,
			// 	User: requestFromOtherUser.ReqFrom,
			// })

			newUser:= models.User{
				Uid : user.Uid,
				Name: user.Name,
				Email: user.Email,
				Password: user.Password,
				Mother : []models.Relation{{
					Rel: requestFromOtherUser.ReqRel,
					User: requestFromOtherUser.ReqFrom,
					},
				},
			}

			pb, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}

			mu.SetJson = pb
			
			_ , errMut := db.Client.NewTxn().Mutate(db.ctx, mu)

			if errMut != nil {
				log.Fatal(errMut)
			}


		} else if (requestFromOtherUser.ReqRel == "Son") {
			// user.Son = append(user.Son, models.Relation{
			// 	Rel: requestFromOtherUser.ReqRel,
			// 	User: requestFromOtherUser.ReqFrom,
			// })

			newUser:= models.User{
				Uid : user.Uid,
				Name: user.Name,
				Email: user.Email,
				Password: user.Password,
				Son : []models.Relation{{
					Rel: requestFromOtherUser.ReqRel,
					User: requestFromOtherUser.ReqFrom,
					},
				},
			}

			pb, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}

			mu.SetJson = pb
			
			_ , errMut := db.Client.NewTxn().Mutate(db.ctx, mu)

			if errMut != nil {
				log.Fatal(errMut)
			}

		}


		

		// ADD FRIEND TO THE OTHER USER TOO
		
		if (requestFromOtherUser.ReqRel == "Friend") {
			// requestFromOtherUser.ReqFrom.Friend = append(requestFromOtherUser.ReqFrom.Friend, models.Relation{
			// 	Rel: requestFromOtherUser.ReqRel,
			// 	User: user,
			// })

			newUser:= models.User{
				Uid : requestFromOtherUser.ReqFrom.Uid,
				Name: requestFromOtherUser.ReqFrom.Name,
				Email: requestFromOtherUser.ReqFrom.Email,
				Password: requestFromOtherUser.ReqFrom.Password,
				Friend : []models.Relation{{
					Rel: requestFromOtherUser.ReqRel,
					User: user,
					},
				},
			}

			pbOtherUser, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}

			muOther := &api.Mutation{
				CommitNow: true,
			}

			muOther.SetJson = pbOtherUser
			
			_ , errMutOther := db.Client.NewTxn().Mutate(db.ctx, muOther)

			if errMutOther != nil {
				log.Fatal(errMutOther)
			}

		} else if (requestFromOtherUser.ReqRel == "Father") {
			// requestFromOtherUser.ReqFrom.Father = append(requestFromOtherUser.ReqFrom.Father, models.Relation{
			// 	Rel: requestFromOtherUser.ReqRel,
			// 	User: user,
			// })

			newUser:= models.User{
				Uid : requestFromOtherUser.ReqFrom.Uid,
				Name: requestFromOtherUser.ReqFrom.Name,
				Email: requestFromOtherUser.ReqFrom.Email,
				Password: requestFromOtherUser.ReqFrom.Password,
				Son : []models.Relation{{
					Rel: "Son",
					User: user,
					},
				},
			}

			pbOtherUser, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}

			muOther := &api.Mutation{
				CommitNow: true,
			}

			muOther.SetJson = pbOtherUser
			
			_ , errMutOther := db.Client.NewTxn().Mutate(db.ctx, muOther)

			if errMutOther != nil {
				log.Fatal(errMutOther)
			}

		} else if (requestFromOtherUser.ReqRel == "Mother") {
			// requestFromOtherUser.ReqFrom.Mother = append(requestFromOtherUser.ReqFrom.Mother, models.Relation{
			// 	Rel: requestFromOtherUser.ReqRel,
			// 	User: user,
			// })

			newUser:= models.User{
				Uid : requestFromOtherUser.ReqFrom.Uid,
				Name: requestFromOtherUser.ReqFrom.Name,
				Email: requestFromOtherUser.ReqFrom.Email,
				Password: requestFromOtherUser.ReqFrom.Password,
				Son : []models.Relation{{
					Rel: "Son",
					User: user,
					},
				},
			}

			pbOtherUser, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}

			muOther := &api.Mutation{
				CommitNow: true,
			}

			muOther.SetJson = pbOtherUser
			
			_ , errMutOther := db.Client.NewTxn().Mutate(db.ctx, muOther)

			if errMutOther != nil {
				log.Fatal(errMutOther)
			}

		} else if (requestFromOtherUser.ReqRel == "Son") {
			// requestFromOtherUser.ReqFrom.Son = append(requestFromOtherUser.ReqFrom.Son, models.Relation{
			// 	Rel: requestFromOtherUser.ReqRel,
			// 	User: user,
			// })

			newUser:= models.User{
				Uid : requestFromOtherUser.ReqFrom.Uid,
				Name: requestFromOtherUser.ReqFrom.Name,
				Email: requestFromOtherUser.ReqFrom.Email,
				Password: requestFromOtherUser.ReqFrom.Password,
				Father : []models.Relation{{
					Rel: "Father",
					User: user,
					},
				},
			}

			pbOtherUser, err := json.Marshal(newUser)
			if err != nil {
				log.Fatal(err)
			}

			muOther := &api.Mutation{
				CommitNow: true,
			}

			muOther.SetJson = pbOtherUser
			
			_ , errMutOther := db.Client.NewTxn().Mutate(db.ctx, muOther)

			if errMutOther != nil {
				log.Fatal(errMutOther)
			}
		}

		



		// Now removing the request from the both user
		// for other user Request
		// for current user RequestFrom
		

		otherUser := requestFromOtherUser.ReqFrom
		
		removeAbleIndex := 0
		
		

		for index, reqItem := range otherUser.Request {
			fmt.Println(reqItem.ReqTo.Uid)
			if (reqItem.ReqTo.Uid == user.Uid) {
				fmt.Printf ("Value of reqItem reqto uid : %s, and user uid  %s\n",reqItem.ReqTo.Uid,user.Uid)
				removeAbleIndex = index
				fmt.Println("Value of remvoeable index : %d and regular index : %d\n",removeAbleIndex,index)
				break
			}
		}

		

		// otherUser.Request = delete_at_index(otherUser.Request,removeAbleIndex)

		// pb2, err := json.Marshal(otherUser)
		type ReMoveRequestStruct struct {
			Uid  string `json:"uid,omitempty"`
			Request models.RelationRequest `json:"request,omitempty"`
		}
		removeAbleReq := ReMoveRequestStruct {
			Uid : otherUser.Uid,
			Request : models.RelationRequest {
				Uid: otherUser.Request[removeAbleIndex].Uid,
			},
		}
		pb2, err := json.Marshal(removeAbleReq)
		if err != nil {
			log.Fatal(err)
		}

		mu2 := &api.Mutation{
			CommitNow: true,
		}
		mu2.DeleteJson = pb2
		
		_, errMut2 := db.Client.NewTxn().Mutate(db.ctx, mu2)

		if errMut2 != nil {
			log.Fatal(errMut2)
		}



		// now remove the request from
		// requestFromOtherUserToRemove := root.GetRequestFromOtherUser[0]

		removeAbleIndex2 := 0
		for index,Item := range user.RequestFrom {
			if (Item.Uid == requestFromOtherUser.Uid) {
				fmt.Printf("Value of item uid : %s and  value of request from other user uid : %s\n",Item.Uid,requestFromOtherUser.Uid)
				removeAbleIndex2 = index
			}
		}

		// user.RequestFrom = delete_at_index_relation_from(user.RequestFrom,removeAbleIndex2)
		
		// pb3, err := json.Marshal(user)
		type ReMoveRequestFromStruct struct {
			Uid  string `json:"uid,omitempty"`
			RequestFrom models.RelationRequestFromOther `json:"request_from,omitempty"`
		}
		fmt.Println("user requestfrom other ",user.RequestFrom[removeAbleIndex2].Uid, " and user id : ",user.Uid)
		removeAbleReqFrom := ReMoveRequestFromStruct {
			Uid : user.Uid,
			RequestFrom : models.RelationRequestFromOther {
				Uid: user.RequestFrom[removeAbleIndex2].Uid,
			},
		}
		pb3, err := json.Marshal(removeAbleReqFrom)
		if err != nil {
			log.Fatal(err)
		}

		mu3 := &api.Mutation{
			CommitNow: true,
		}

		mu3.DeleteJson = pb3
		
		_, errMut3 := db.Client.NewTxn().Mutate(db.ctx, mu3)

		if errMut3 != nil {
			log.Fatal(errMut3)
		}

		return true

	} else {
		// do nothing
		return false
	}


}

/*Cancels request from other user*/
func (db *DB) CancelReq(user_id string,req_id string) bool {
	db.InitSchema()
	
	user := db.Me(user_id)

	query := `
	{
		getRequestFromOtherUser(func: uid(`+req_id+`)) {
			uid
			req_rel
			req_from {
					uid
			name
			email

			request {
				uid
				req_to {
					uid
				}
			}
			request_from {
				uid
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
		GetRequestFromOtherUser []models.RelationRequestFromOther `json:"getRequestFromOtherUser,omitempty"`
	}
	var root Root
	
	if err := json.Unmarshal(resp.GetJson(), &root); err != nil {
		log.Fatal(err)
	}

	if (len(root.GetRequestFromOtherUser) > 0) {
		requestFromOtherUser := root.GetRequestFromOtherUser[0]

		otherUser := requestFromOtherUser.ReqFrom
		
		removeAbleIndex := 0
		
		

		for index, reqItem := range otherUser.Request {
			fmt.Println(reqItem.ReqTo.Uid)
			if (reqItem.ReqTo.Uid == user.Uid) {
				fmt.Printf ("Value of reqItem reqto uid : %s, and user uid  %s\n",reqItem.ReqTo.Uid,user.Uid)
				removeAbleIndex = index
				fmt.Println("Value of remvoeable index : %d and regular index : %d\n",removeAbleIndex,index)
				break
			}
		}

		

		// otherUser.Request = delete_at_index(otherUser.Request,removeAbleIndex)

		// pb2, err := json.Marshal(otherUser)
		type ReMoveRequestStruct struct {
			Uid  string `json:"uid,omitempty"`
			Request models.RelationRequest `json:"request,omitempty"`
		}
		removeAbleReq := ReMoveRequestStruct {
			Uid : otherUser.Uid,
			Request : models.RelationRequest {
				Uid: otherUser.Request[removeAbleIndex].Uid,
			},
		}
		pb2, err := json.Marshal(removeAbleReq)
		if err != nil {
			log.Fatal(err)
		}

		mu2 := &api.Mutation{
			CommitNow: true,
		}
		mu2.DeleteJson = pb2
		
		_, errMut2 := db.Client.NewTxn().Mutate(db.ctx, mu2)

		if errMut2 != nil {
			log.Fatal(errMut2)
		}



		// now remove the request from
		// requestFromOtherUserToRemove := root.GetRequestFromOtherUser[0]

		removeAbleIndex2 := 0
		for index,Item := range user.RequestFrom {
			if (Item.Uid == requestFromOtherUser.Uid) {
				fmt.Printf("Value of item uid : %s and  value of request from other user uid : %s\n",Item.Uid,requestFromOtherUser.Uid)
				removeAbleIndex2 = index
			}
		}

		// user.RequestFrom = delete_at_index_relation_from(user.RequestFrom,removeAbleIndex2)
		
		// pb3, err := json.Marshal(user)
		type ReMoveRequestFromStruct struct {
			Uid  string `json:"uid,omitempty"`
			RequestFrom models.RelationRequestFromOther `json:"request_from,omitempty"`
		}
		fmt.Println("user requestfrom other ",user.RequestFrom[removeAbleIndex2].Uid, " and user id : ",user.Uid)
		removeAbleReqFrom := ReMoveRequestFromStruct {
			Uid : user.Uid,
			RequestFrom : models.RelationRequestFromOther {
				Uid: user.RequestFrom[removeAbleIndex2].Uid,
			},
		}
		pb3, err := json.Marshal(removeAbleReqFrom)
		if err != nil {
			log.Fatal(err)
		}

		mu3 := &api.Mutation{
			CommitNow: true,
		}

		mu3.DeleteJson = pb3
		
		_, errMut3 := db.Client.NewTxn().Mutate(db.ctx, mu3)

		if errMut3 != nil {
			log.Fatal(errMut3)
		}

		return true

	} else {
		return false
	}

	

}

/*Cancels request to other user*/
// func (db *DB) CancelToReq(user_id string,req_id string) {

// }

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




