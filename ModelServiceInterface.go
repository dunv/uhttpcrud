package uhttpcrud

import (
	"github.com/dunv/auth"
	"github.com/dunv/mongo"
	"gopkg.in/mgo.v2/bson"
)

// ModelService <-
type ModelService interface {
	CopyAndInit(db *mongo.DbSession) ModelService
	GetIndexProperties() string
	GetByIndexProperties(interface{}) (interface{}, error)
	CheckNotNullable(interface{}) bool
	Get(ID *bson.ObjectId) (interface{}, error)
	List() (interface{}, error)
	Create(obj interface{}, user auth.User) error
	Update(obj interface{}, user auth.User) error
	Delete(id bson.ObjectId) error
}
