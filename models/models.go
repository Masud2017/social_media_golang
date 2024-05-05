package models

import (
	"encoding/json"
)

type User {
	Uid      string     `json:"uid,omitempty"`
	Name     string     `json:"name,omitempty"`
	Email	 string     `json:"email,omitempty"`
	Password string     `json:"password,omitempty"`
	
	Relations  []Relation   `json:"relations,omitempty"`
	
}

type Relation {
	RelationName string `json:"relation_name,omitempty"`
	User User `json:"user,omitempty"`
}