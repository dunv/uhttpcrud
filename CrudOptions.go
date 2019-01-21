package uhttpcrud

import (
	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

// CrudOptions <-
type CrudOptions struct {
	IDParameterName  string
	ListPermission   *uauth.Permission
	ListEndpoint     *string
	GetPermission    *uauth.Permission
	GetEndpoint      *string
	CreatePermission *uauth.Permission
	CreateEndpoint   *string
	UpdatePermission *uauth.Permission
	UpdateEndpoint   *string
	DeletePermission *uauth.Permission
	DeleteEndpoint   *string
	ModelService     ModelService
	Model            WithID
}

// CreateEndpoints <-
func (o CrudOptions) CreateEndpoints() {
	uhttp.Handle(*o.GetEndpoint, GenericGetHandler(o))
	uhttp.Handle(*o.ListEndpoint, GenericListHandler(o))
	uhttp.Handle(*o.CreateEndpoint, GenericCreateHandler(o))
	uhttp.Handle(*o.UpdateEndpoint, GenericUpdateHandler(o))
	uhttp.Handle(*o.DeleteEndpoint, GenericDeleteHandler(o))
}
