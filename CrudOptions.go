package uhttpcrud

import (
	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

// CrudOptions <-
type CrudOptions struct {
	IDParameterName string
	Database        string

	ListEndpoint         *string
	ListPermission       *uauth.Permission
	ListOthersPermission *uauth.Permission

	GetEndpoint         *string
	GetPermission       *uauth.Permission
	GetOthersPermission *uauth.Permission

	CreateEndpoint         *string
	CreatePermission       *uauth.Permission
	CreateOthersPermission *uauth.Permission

	UpdateEndpoint         *string
	UpdatePermission       *uauth.Permission
	UpdateOthersPermission *uauth.Permission

	DeleteEndpoint         *string
	DeletePermission       *uauth.Permission
	DeleteOthersPermission *uauth.Permission

	ModelService ModelService
	Model        WithID
}

// CreateEndpoints <-
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
