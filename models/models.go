package models

type User struct {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Email	 string     `json:"email,omitempty"`
	Password string     `json:"password,omitempty"`
	
	Friend  []Relation   `json:"friend,omitempty"`
	// DType    []string   `json:"dgraph.type,omitempty"`

	
}

type Relation struct{
	Rel string `json:"rel,omitempty"`
	User User `json:"user,omitempty"`
	// DType    []string   `json:"dgraph.type,omitempty"`
}

type RelationRequest {
	Uid string `json:"uid,omitempty"`
	ReqFor string `json:"req_for,omitempty"` // user id for the User who will own this relation
	Rel string `json:"rel,omitempty"`
	ReqTo string `json:"req_to,omitempty"` // user id for the User who will be the relative
}