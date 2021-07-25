package uhttpcrud

import (
	"context"
	"errors"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

type CrudOptions struct {
	// IDParameterName is the name of the GET parameter which is used
	// for getting and deleting documents
	IDParameterName string

	// ListEndpoint is the http-endpoint-name for list queries
	// If == nil there will be no list-enpoint
	ListEndpoint         *string
	ListPreprocess       func(context.Context) error
	ListPermission       *uauth.Permission
	ListOthersPermission *uauth.Permission
	ListRequiredGet      uhttp.R
	ListOptionalGet      uhttp.R

	// GetEndpoint is the http-endpoint-name for get queries
	// If == nil there will be no get-enpoint
	GetEndpoint         *string
	GetPreprocess       func(context.Context) error
	GetPermission       *uauth.Permission
	GetOthersPermission *uauth.Permission
	GetRequiredGet      uhttp.R
	GetOptionalGet      uhttp.R

	// CreateEndpoint is the http-endpoint-name for create queries
	// If == nil there will be no create-enpoint
	CreateEndpoint         *string
	CreatePreprocess       func(context.Context) error
	CreatePermission       *uauth.Permission
	CreateOthersPermission *uauth.Permission
	CreateRequiredGet      uhttp.R
	CreateOptionalGet      uhttp.R

	// UpdateEndpoint is the http-endpoint-name for update queries
	// If == nil there will be no update-enpoint
	UpdateEndpoint         *string
	UpdatePreprocess       func(context.Context) error
	UpdatePermission       *uauth.Permission
	UpdateOthersPermission *uauth.Permission
	UpdateRequiredGet      uhttp.R
	UpdateOptionalGet      uhttp.R

	// DeleteEndpoint is the http-endpoint-name for delete queries
	// If == nil there will be no delete-enpoint
	DeleteEndpoint         *string
	DeletePreprocess       func(context.Context) error
	DeletePermission       *uauth.Permission
	DeleteOthersPermission *uauth.Permission
	DeleteRequiredGet      uhttp.R
	DeleteOptionalGet      uhttp.R

	// ModelService will be called upon for all database interactions
	ModelService ModelService

	// Model will be used to parse and validate models given to create/update handlers
	Model WithID
}

// CreateEndpoints adds all handlers configured in CrudOptions using the uhttp-framework
func (o CrudOptions) CreateEndpoints(u *uhttp.UHTTP) error {
	if o.GetEndpoint != nil {
		if o.ModelService == nil || o.IDParameterName == "" {
			return errors.New("crudOptions.ModelService and crudOptions.IDParameterName is required when using GetEndpoint")
		}
		u.Handle(*o.GetEndpoint, genericGetHandler(o))
	}
	if o.ListEndpoint != nil {
		if o.ModelService == nil {
			return errors.New("crudOptions.ModelService is required when using ListEndpoint")
		}
		u.Handle(*o.ListEndpoint, genericListHandler(o))
	}
	if o.CreateEndpoint != nil {
		if o.ModelService == nil || o.Model == nil {
			return errors.New("crudOptions.ModelService and crudOptions.Model are required when using CreateEndpoint")
		}
		u.Handle(*o.CreateEndpoint, genericCreateHandler(o))
	}
	if o.UpdateEndpoint != nil {
		if o.ModelService == nil || o.Model == nil {
			return errors.New("crudOptions.ModelService, crudOptions.Model are required when using UpdateEndpoint")
		}
		u.Handle(*o.UpdateEndpoint, genericUpdateHandler(o))
	}
	if o.DeleteEndpoint != nil {
		if o.ModelService == nil || o.IDParameterName == "" {
			return errors.New("crudOptions.ModelService, crudOptions.IDParameterName are required when using DeleteEndpoint")
		}
		u.Handle(*o.DeleteEndpoint, genericDeleteHandler(o))
	}
	return nil
}
