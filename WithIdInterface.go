package uhttpcrud

import "go.mongodb.org/mongo-driver/bson/primitive"

// WithID <-
type WithID interface {
	GetID() *primitive.ObjectID
}
