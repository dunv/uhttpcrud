package uhttpcrud

import (
	"context"
)

// ModelService is an interface which all CRUD-http-handlers will use
type ModelService interface {
	// Validate should validate the model which is created/updated
	// - called from the createHandler, before it calls service.Create
	// - called from the updateHandler, before it calls service.Update
	Validate(interface{}) bool

	// Get retrieves a document by its ID (typically a string or ObjectID etc.)
	// If limitToUser is true the service should only return documents which belong to the user
	Get(ID string, ctx context.Context) (interface{}, error)

	// List retrieves all documents which this user has access to
	// If limitToUser is true the service should only return documents which belong to the user
	// ctx will contain the request-context. This way we can pass on filters from the reuqest
	// into the service
	List(ctx context.Context) (interface{}, error)

	// Create creates a document in the database and returns the new document
	// If permissions are implemented, the service should make this created document belong to the
	// user passed into this method
	Create(obj interface{}, ctx context.Context) (interface{}, error)

	// Update updates a document. It is up to the implementer to get the ID-property, etc.
	// It returns the updated document
	// If limitToUser is true the service should check if this user is allowed to modify this document
	Update(obj interface{}, ctx context.Context) (interface{}, error)

	// Delete deletes a document by its ID (typically a string or ObjectID etc.)
	// If limitToUser is true the service should check if this user is allowed to delete this document
	Delete(ID string, ctx context.Context) error
}
