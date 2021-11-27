package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

// Returns an instance of an update-handler for the configured options
func genericCreateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPreProcess(options.CreatePreprocess),
		uhttp.WithRequiredGet(options.CreateRequiredGet),
		uhttp.WithOptionalGet(options.CreateOptionalGet),
		uhttp.WithMiddlewares(uauth.AuthJWT()),
		uhttp.WithPostModel(options.Model, func(r *http.Request, model interface{}, ret *int) interface{} {
			// Sanity check: CreateOthersPermission can only be set if CreatePermission is set
			if options.CreatePermission == nil && options.CreateOthersPermission != nil {
				*ret = http.StatusInternalServerError
				return map[string]string{"err": "Configuration problem: CreateOthersPermission can only be set if CreatePermission is set."}
			}

			// Get User
			user, err := uauth.UserFromRequest(r)
			if err != nil {
				return fmt.Errorf("Could not get user (%s)", err)
			}

			if options.CreatePermission != nil && !user.CheckPermission(*options.CreatePermission) {
				return fmt.Errorf("User does not have the required permission: %s", *options.CreatePermission)
			}

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
