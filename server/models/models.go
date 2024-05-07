package models

type User struct {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Email	 string     `json:"email,omitempty"`
	Password string     `json:"password,omitempty"`
	
	Friend  []Relation   `json:"friend,omitempty"`
	Father  []Relation   `json:"father,omitempty"`
	Mother  []Relation   `json:"mother,omitempty"`
	Son  []Relation   `json:"son,omitempty"`
	// DType    []string   `json:"dgraph.type,omitempty"`
	Request []RelationRequest `json:"request,omitempty"`
	RequestFrom []RelationRequestFromOther `json:"request_from,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`


	
}

type Relation struct{
	Rel string `json:"rel,omitempty"`
	User User `json:"user,omitempty"`
	DType []string `json:"dgraph.type,omitempty"`
	// DType    []string   `json:"dgraph.type,omitempty"`
}

/*
These request will be done by the user him/her self
*/
type RelationRequest struct {
	Uid string `json:"uid,omitempty"`
	ReqRel string `json:"req_rel,omitempty"`
	ReqTo User `json:"req_to,omitempty"` // user id for the User who will be the relative
	ReqFromUid string `json:"req_from_uid,omitempty"`
}

/*
This request will come from others users to the current user
*/
type RelationRequestFromOther struct {
	Uid string `json:"uid,omitempty"`
	ReqRel string `json:"req_rel,omitempty"`
	ReqFrom User `json:"req_from,omitempty"` // user id for the User who will be the relative
	ReqFromUid string `json:"req_from_uid,omitempty"`
}