package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

func genericUpdateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.NewHandler(
		uhttp.WithPreProcess(options.UpdatePreprocess),
		uhttp.WithRequiredGet(options.UpdateRequiredGet),
		uhttp.WithOptionalGet(options.UpdateOptionalGet),
		uhttp.WithMiddlewares(uauth.AuthJWT()),
		uhttp.WithPostModel(options.Model, func(r *http.Request, model interface{}, ret *int) interface{} {

			// Sanity check: UpdateOthersPermission can only be set if UpdatePermission is set
			if options.UpdatePermission == nil && options.UpdateOthersPermission != nil {
				*ret = http.StatusInternalServerError
				return map[string]string{"err": "Configuration problem: UpdateOthersPermission can only be set if UpdatePermission is set."}
			}

			// Check permissions
			var user *uauth.User
			var limitToUser *uauth.User
			var err error
			if options.UpdatePermission != nil {
				user, err = uauth.UserFromRequest(r)
				if err != nil {
					return fmt.Errorf("Could not get user (%s)", err)
				}

				// Return nothing, if updatePermission is required but the user does not have it
				if !user.CheckPermission(*options.UpdatePermission) {
					return fmt.Errorf("User does not have the required permission: %s", *options.UpdatePermission)
				}

				// Limit results if UpdateOthersPermission is required but the user does not have it
				if options.UpdateOthersPermission != nil {
					if !user.CheckPermission(*options.UpdateOthersPermission) {
						limitToUser = user
					}
				}
			}

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
			_, err = options.ModelService.Get(idFromModel, limitToUser != nil, r.Context())
			if err != nil {
				return fmt.Errorf("No object with the id %s exists (%s)", idFromModel, err)
			}

			// Actual update
			updatedDocument, err := options.ModelService.Update(modelInterface, limitToUser != nil, r.Context())
			if err != nil {
				return err
			}

			// Answer
			return updatedDocument
		}),
	)
}
