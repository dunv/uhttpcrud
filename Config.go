package uhttpcrud

import (
	"github.com/dunv/uhttp"
)

var dbContextKey uhttp.ContextKey

// SetDBContextKey for users and roles
func SetDBContextKey(_dbContextKey uhttp.ContextKey) {
	dbContextKey = _dbContextKey
}
