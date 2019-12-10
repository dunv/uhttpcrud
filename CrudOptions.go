package uhttpcrud

import (
	"context"

	uauthPermissions "github.com/dunv/uauth/permissions"
	"github.com/dunv/uhttp"
	"github.com/dunv/uhttp/params"
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
	ListPermission       *uauthPermissions.Permission
	ListOthersPermission *uauthPermissions.Permission
	ListRequiredGet      params.R
	ListOptionalGet      params.R

	// GetEndpoint is the http-endpoint-name for get queries
	// If == nil there will be no get-enpoint
	GetEndpoint         *string
	GetPreprocess       func(context.Context) error
	GetPermission       *uauthPermissions.Permission
	GetOthersPermission *uauthPermissions.Permission
	GetRequiredGet      params.R
	GetOptionalGet      params.R

	// CreateEndpoint is the http-endpoint-name for create queries
	// If == nil there will be no create-enpoint
	CreateEndpoint         *string
	CreatePreprocess       func(context.Context) error
	CreatePermission       *uauthPermissions.Permission
	CreateOthersPermission *uauthPermissions.Permission
	CreateRequiredGet      params.R
	CreateOptionalGet      params.R

	// UpdateEndpoint is the http-endpoint-name for update queries
	// If == nil there will be no update-enpoint
	UpdateEndpoint         *string
	UpdatePreprocess       func(context.Context) error
	UpdatePermission       *uauthPermissions.Permission
	UpdateOthersPermission *uauthPermissions.Permission
	UpdateRequiredGet      params.R
	UpdateOptionalGet      params.R

	// DeleteEndpoint is the http-endpoint-name for delete queries
	// If == nil there will be no delete-enpoint
	DeleteEndpoint         *string
	DeletePreprocess       func(context.Context) error
	DeletePermission       *uauthPermissions.Permission
	DeleteOthersPermission *uauthPermissions.Permission
	DeleteRequiredGet      params.R
	DeleteOptionalGet      params.R

	// ModelService will be called upon for all database interactions
	ModelService ModelService

	// Model will be used to parse and validate models given to create/update handlers
	Model WithID
}

// CreateEndpoints adds all handlers configured in CrudOptions using the uhttp-framework
func (o CrudOptions) CreateEndpoints() {
	if o.GetEndpoint != nil {
		uhttp.Handle(*o.GetEndpoint, genericGetHandler(o))
	}
	if o.ListEndpoint != nil {
		uhttp.Handle(*o.ListEndpoint, genericListHandler(o))
	}
	if o.CreateEndpoint != nil {
		uhttp.Handle(*o.CreateEndpoint, genericCreateHandler(o))
	}
	if o.UpdateEndpoint != nil {
		uhttp.Handle(*o.UpdateEndpoint, genericUpdateHandler(o))
	}
	if o.DeleteEndpoint != nil {
		uhttp.Handle(*o.DeleteEndpoint, genericDeleteHandler(o))
	}
}
