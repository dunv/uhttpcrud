package uhttpcrud

// WithID makes sure the struct in question has a gettable ID property for DB-operations
type WithID interface {
	// Returns the ID of this struct (typically a string or ObjectID, etc.)
	GetID() interface{}
}
