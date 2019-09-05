package uhttpcrud

import (
	"context"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

type CrudOptions struct {
	// IDParameterName is the name of the GET parameter which is used
	// for getting and deleting documents
	IDParameterName string

	// Database is the name of the database which contains the documents to be modified (will be passed to all service-calls)
	Database string

	// ListEndpoint is the http-endpoint-name for list queries
	// If == nil there will be no list-enpoint
	ListEndpoint         *string
	ListPreprocess       func(context.Context) error
	ListPermission       *uauth.Permission
	ListOthersPermission *uauth.Permission

	// GetEndpoint is the http-endpoint-name for get queries
	// If == nil there will be no get-enpoint
	GetEndpoint         *string
	GetPreprocess       func(context.Context) error
	GetPermission       *uauth.Permission
	GetOthersPermission *uauth.Permission

	// CreateEndpoint is the http-endpoint-name for create queries
	// If == nil there will be no create-enpoint
	CreateEndpoint         *string
	CreatePreprocess       func(context.Context) error
	CreatePermission       *uauth.Permission
	CreateOthersPermission *uauth.Permission

	// UpdateEndpoint is the http-endpoint-name for update queries
	// If == nil there will be no update-enpoint
	UpdateEndpoint         *string
	UpdatePreprocess       func(context.Context) error
	UpdatePermission       *uauth.Permission
	UpdateOthersPermission *uauth.Permission

	// DeleteEndpoint is the http-endpoint-name for delete queries
	// If == nil there will be no delete-enpoint
	DeleteEndpoint         *string
	DeletePreprocess       func(context.Context) error
	DeletePermission       *uauth.Permission
	DeleteOthersPermission *uauth.Permission

	// ModelService will be called upon for all database interactions
	ModelService ModelService

	// Model will be used to parse and validate models given to create/update handlers
	Model WithID
}

// CreateEndpoints adds all handlers configured in CrudOptions using the uhttp-framework
func (o CrudOptions) CreateEndpoints() {
	if o.GetEndpoint != nil {
		uhttp.Handle(*o.GetEndpoint, GenericGetHandler(o))
	}
	if o.ListEndpoint != nil {
		uhttp.Handle(*o.ListEndpoint, GenericListHandler(o))
	}
	if o.CreateEndpoint != nil {
		uhttp.Handle(*o.CreateEndpoint, GenericCreateHandler(o))
	}
	if o.UpdateEndpoint != nil {
		uhttp.Handle(*o.UpdateEndpoint, GenericUpdateHandler(o))
	}
	if o.DeleteEndpoint != nil {
		uhttp.Handle(*o.DeleteEndpoint, GenericDeleteHandler(o))
	}
}
