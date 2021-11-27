package uhttpcrud

import (
	"net/http"

	"github.com/dunv/uhttp"
)

// Returns an instance of an list-handler for the configured options
func genericListHandler(options CrudOptions) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPreProcess(options.ListPreprocess),
		uhttp.WithRequiredGet(options.ListRequiredGet),
		uhttp.WithOptionalGet(options.ListOptionalGet),
		uhttp.WithMiddlewares(options.ListMiddleware),
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			// Load
			objsFromDb, err := options.ModelService.List(r.Context())
			if err != nil {
				return err
			}

			// Render Response
			return objsFromDb
		}),
	)
}
