package uhttpcrud

import (
	"net/http"

	contextKeys "github.com/dunv/uhttp/contextkeys"
)

// WithID makes sure the struct in question has a gettable ID property for DB-operations
type WithID interface {
	// Returns the ID of this struct (typically a string or ObjectID, etc.)
	GetID() string
}

func getWithIDFromPostModel(r *http.Request) WithID {
	return r.Context().Value(contextKeys.CtxKeyPostModel).(WithID)
}
