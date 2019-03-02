package uhttpcrud

import (
	"github.com/dunv/uauth"
	"github.com/dunv/umongo"
	"gopkg.in/mgo.v2/bson"
)

// ModelService <-
type ModelService interface {
	CopyAndInit(db *umongo.DbSession) ModelService
	GetIndexProperties() string
	GetByIndexProperties(interface{}) (interface{}, error)
	CheckNotNullable(interface{}) bool
	Get(ID *bson.ObjectId, user *uauth.User) (interface{}, error)
	List(user *uauth.User) (interface{}, error)
	Create(obj interface{}, user uauth.User) (*bson.ObjectId, error)
	Update(obj interface{}, user uauth.User) error
	Delete(id bson.ObjectId, user *uauth.User) error
}
