package uhttpcrud

import (
	"github.com/dunv/uauth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ModelService <-
type ModelService interface {
	CopyAndInit(db *mongo.Client, database string) ModelService
	GetIndexProperties() string
	GetByIndexProperties(interface{}) (interface{}, error)
	CheckNotNullable(interface{}) bool
	Get(ID *primitive.ObjectID, user *uauth.User) (interface{}, error)
	List(user *uauth.User) (interface{}, error)
	Create(obj interface{}, user uauth.User) (*primitive.ObjectID, error)
	Update(obj interface{}, user uauth.User) error
	Delete(id primitive.ObjectID, user *uauth.User) error
}
