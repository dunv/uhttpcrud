package uhttpcrud

import (
	"context"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

type CrudOptions struct {
	IDParameterName string
	Database        string

	ListEndpoint         *string
	ListPreprocess       func(context.Context) error
	ListPermission       *uauth.Permission
	ListOthersPermission *uauth.Permission

	GetEndpoint         *string
	GetPreprocess       func(context.Context) error
	GetPermission       *uauth.Permission
	GetOthersPermission *uauth.Permission

	CreateEndpoint         *string
	CreatePreprocess       func(context.Context) error
	CreatePermission       *uauth.Permission
	CreateOthersPermission *uauth.Permission

	UpdateEndpoint         *string
	UpdatePreprocess       func(context.Context) error
	UpdatePermission       *uauth.Permission
	UpdateOthersPermission *uauth.Permission

	DeleteEndpoint         *string
	DeletePreprocess       func(context.Context) error
	DeletePermission       *uauth.Permission
	DeleteOthersPermission *uauth.Permission

	ModelService ModelService
	Model        WithID
}

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
