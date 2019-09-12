package uhttpcrud

import (
	contextKeys "github.com/dunv/uhttp/contextkeys"
)

var dbContextKey contextKeys.ContextKey

// SetDBContextKey for users and roles
func SetDBContextKey(_dbContextKey contextKeys.ContextKey) {
	dbContextKey = _dbContextKey
}
