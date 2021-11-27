package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uhttp"
)

// Returns an instance of an update-handler for the configured options
func genericCreateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPreProcess(options.CreatePreprocess),
		uhttp.WithRequiredGet(options.CreateRequiredGet),
		uhttp.WithOptionalGet(options.CreateOptionalGet),
		uhttp.WithMiddlewares(options.CreateMiddleware),
		uhttp.WithPostModel(options.Model, func(r *http.Request, model interface{}, ret *int) interface{} {
			modelInterface := getWithIDFromPostModel(r)

			// Check if all required populated fields are populated (indexes)
			if !options.ModelService.Validate(modelInterface) {
				return fmt.Errorf("Non-nullable properties are null")
			}

			// Create (will return an error if already exists)
			createdDocument, err := options.ModelService.Create(modelInterface, r.Context())
			if err != nil {
				return err
			}

			// Answer
			return createdDocument
		}),
	)
}
