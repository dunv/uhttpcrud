package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uauth"
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
		uhttp.WithMiddlewares(uauth.AuthJWT()), // We need a user in order to delete an object
		uhttp.WithRequiredGet(requiredGet),
		uhttp.WithOptionalGet(options.DeleteOptionalGet),
		uhttp.WithDelete(func(r *http.Request, ret *int) interface{} {
			// Sanity check: DeleteOthersPermission can only be set if DeletePermission is set
			if options.DeletePermission == nil && options.DeleteOthersPermission != nil {
				*ret = http.StatusInternalServerError
				return map[string]string{"err": "Configuration problem: DeleteOthersPermission can only be set if DeletePermission is set."}
			}

			// Check permissions
			var user *uauth.User
			var limitToUser *uauth.User
			var err error
			if options.DeletePermission != nil {
				// Return nothing, if deletePermission is required but the user does not have it
				user, err = uauth.UserFromRequest(r)
				if err != nil {
					return fmt.Errorf("Could not get user (%s)", err)
				}

				if !user.CheckPermission(*options.DeletePermission) {
					return fmt.Errorf("User does not have the required permission: %s", *options.DeletePermission)
				}

				// Limit results if DeleteOthersPermission is required but the user does not have it
				if options.DeleteOthersPermission != nil {
					if !user.CheckPermission(*options.DeleteOthersPermission) {
						limitToUser = user
					}
				}
			}

			// GetDB
			objectID := uhttp.GetAsString(options.IDParameterName, r)

			// Delete
			err = options.ModelService.Delete(*objectID, limitToUser != nil, r.Context())
			if err != nil {
				return err
			}

			// Answer
			return map[string]string{"msg": "Deleted successfully"}
		}),
	)
}
