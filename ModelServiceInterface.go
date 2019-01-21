package uhttpcrud

import (
	"github.com/dunv/umongo"
	"gopkg.in/mgo.v2/bson"
)

// ModelService <-
type ModelService interface {
	CopyAndInit(db *umongo.DbSession) ModelService
	GetIndexProperties() string
	GetByIndexProperties(interface{}) (interface{}, error)
	CheckNotNullable(interface{}) bool
	Get(ID *bson.ObjectId) (interface{}, error)
	List() (interface{}, error)
	Create(obj interface{}, user auth.User) error
	Update(obj interface{}, user auth.User) error
	Delete(id bson.ObjectId) error
}
