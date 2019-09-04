package uhttpcrud

import (
	"github.com/dunv/uauth"
	"go.mongodb.org/mongo-driver/mongo"
)

// ModelService <-
type ModelService interface {
	CopyAndInit(db *mongo.Client, database string) ModelService
	GetIndexProperties() string
	GetByIndexProperties(interface{}) (interface{}, error)
	CheckNotNullable(interface{}) bool
	Get(ID interface{}, user *uauth.User) (interface{}, error)
	List(user *uauth.User) (interface{}, error)
	Create(obj interface{}, user uauth.User) (interface{}, error)
	Update(obj interface{}, user uauth.User) error
	Delete(id interface{}, user *uauth.User) error
}
