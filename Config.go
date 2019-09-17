package uhttpcrud

var dbContextKey string

// SetDBContextKey for users and roles
func SetDBContextKey(_dbContextKey string) {
	dbContextKey = _dbContextKey
}
