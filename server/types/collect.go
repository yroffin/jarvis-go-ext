package types

import (
	"gopkg.in/mgo.v2/bson"
)

// CollectResource : connector struct
type CollectResource struct {
	Collections []string `json:"collections,omitempty"`
	Data        []bson.M `json:"data,omitempty"`
}
