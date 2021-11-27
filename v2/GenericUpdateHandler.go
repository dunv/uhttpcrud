package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uhttp"
)

func genericUpdateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPreProcess(options.UpdatePreprocess),
		uhttp.WithRequiredGet(options.UpdateRequiredGet),
		uhttp.WithOptionalGet(options.UpdateOptionalGet),
		uhttp.WithMiddlewares(options.UpdateMiddleware),
		uhttp.WithPostModel(options.Model, func(r *http.Request, model interface{}, ret *int) interface{} {
			modelInterface := getWithIDFromPostModel(r)

			// Check if all required populated fields are populated (indexes)
			idFromModel, err := modelInterface.GetID()
			if err != nil {
				return fmt.Errorf("could not getID (%s)", err)
			}

			if idFromModel == "" || !options.ModelService.Validate(modelInterface) {
				return fmt.Errorf("Non-nullable properties are null or no ID present")
			}

			// Check if already exists
			_, err = options.ModelService.Get(idFromModel, r.Context())
			if err != nil {
				return fmt.Errorf("No object with the id %s exists (%s)", idFromModel, err)
			}

			// Actual update
			updatedDocument, err := options.ModelService.Update(modelInterface, r.Context())
			if err != nil {
				return err
			}

			// Answer
			return updatedDocument
		}),
	)
}
