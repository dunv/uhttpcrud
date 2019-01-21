package uhttpcrud

import (
	"gopkg.in/mgo.v2/bson"
)

// WithID <-
type WithID interface {
	GetID() *bson.ObjectId
}
