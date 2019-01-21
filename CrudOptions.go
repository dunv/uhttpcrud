package uhttpcrud

import (
	"github.com/dunv/uhttp"
)

// CrudOptions <-
type CrudOptions struct {
	IDParameterName  string
	ListPermission   *auth.Permission
	ListEndpoint     *string
	GetPermission    *auth.Permission
	GetEndpoint      *string
	CreatePermission *auth.Permission
	CreateEndpoint   *string
	UpdatePermission *auth.Permission
	UpdateEndpoint   *string
	DeletePermission *auth.Permission
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
