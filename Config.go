package uhttpcrud

import (
	uhttpModels "github.com/dunv/uhttp/models"
)

var dbContextKey uhttpModels.ContextKey

// SetDBContextKey for users and roles
func SetDBContextKey(_dbContextKey uhttpModels.ContextKey) {
	dbContextKey = _dbContextKey
}
