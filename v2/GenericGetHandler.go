package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uhttp"
)

// Returns an instance of an get-handler for the configured options
func genericGetHandler(options CrudOptions) uhttp.Handler {
	requiredGet := options.GetRequiredGet
	if requiredGet == nil {
		requiredGet = uhttp.R{}
	}
	requiredGet[options.IDParameterName] = uhttp.STRING

	return uhttp.NewHandler(
		uhttp.WithPreProcess(options.GetPreprocess),
		uhttp.WithRequiredGet(requiredGet),
		uhttp.WithOptionalGet(options.GetOptionalGet),
		uhttp.WithMiddlewares(options.GetMiddleware),
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			// Get
			objectID := uhttp.GetAsString(options.IDParameterName, r)
			var objFromDb interface{}
			objFromDb, err := options.ModelService.Get(*objectID, r.Context())
			if err != nil {
				return fmt.Errorf("Could not find object with ID: '%s' (%s)", *objectID, err)
			}

			return objFromDb
		}),
	)
}
