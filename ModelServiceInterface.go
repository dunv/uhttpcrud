package uhttpcrud

import (
	"github.com/dunv/uauth"
	"go.mongodb.org/mongo-driver/mongo"
)

// ModelService is an interface which all CRUD-http-handlers will use
type ModelService interface {

	// CopyAndInit allows us to clone a service
	CopyAndInit(db *mongo.Client, database string) ModelService

	// Validate should validate the model which is created/updated
	// - called from the createHandler, before it calls service.Create
	// - called from the updateHandler, before it calls service.Update
	Validate(interface{}) bool

	// Get retrieves a document by its ID (typically a string or ObjectID etc.)
	// If user is != nil the service should only return documents which belong to the user
	Get(ID interface{}, user *uauth.User, limitToUser bool) (interface{}, error)

	// List retrieves all documents which this user has access to
	// If user is != nil the service should only return documents which belong to the user
	List(user *uauth.User, limitToUser bool) (interface{}, error)

	// Create creates a document in the database and returns the new document
	// If permissions are implemented, the service should make this created document belong to the
	// user passed into this method
	Create(obj interface{}, user *uauth.User) (interface{}, error)

	// Update updates a document. It is up to the implementer to get the ID-property, etc.
	// It returns the updated document
	// If user is != nil the service should check if this user is allowed to modify this document
	Update(obj interface{}, user *uauth.User, limitToUser bool) (interface{}, error)

	// Delete deletes a document by its ID (typically a string or ObjectID etc.)
	// If user is != nil the service should check if this user is allowed to delete this document
	Delete(id interface{}, user *uauth.User, limitToUser bool) error
}
