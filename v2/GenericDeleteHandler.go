package uhttpcrud

import (
	"net/http"

	"github.com/dunv/uhttp"
)

// Returns an instance of an delete-handler for the configured options
func genericDeleteHandler(options CrudOptions) uhttp.Handler {
	requiredGet := options.DeleteRequiredGet
	if requiredGet == nil {
		requiredGet = uhttp.R{}
	}
	requiredGet[options.IDParameterName] = uhttp.STRING

	return uhttp.NewHandler(
		uhttp.WithPreProcess(options.DeletePreprocess),
		uhttp.WithRequiredGet(requiredGet),
		uhttp.WithOptionalGet(options.DeleteOptionalGet),
		uhttp.WithMiddlewares(options.DeleteMiddleware),
		uhttp.WithDelete(func(r *http.Request, ret *int) interface{} {
			// GetDB
			objectID := uhttp.GetAsString(options.IDParameterName, r)

			// Delete
			err := options.ModelService.Delete(*objectID, r.Context())
			if err != nil {
				return err
			}

			// Answer
			return map[string]string{"msg": "Deleted successfully"}
		}),
	)
}
