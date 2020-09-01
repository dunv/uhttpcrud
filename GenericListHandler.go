package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

// Returns an instance of an list-handler for the configured options
func genericListHandler(options CrudOptions) uhttp.Handler {
	var middlewares []uhttp.Middleware
	if options.ListPermission != nil {
		middlewares = []uhttp.Middleware{uauth.AuthJWT()}
	}
	return uhttp.NewHandler(
		uhttp.WithPreProcess(options.ListPreprocess),
		uhttp.WithMiddlewares(middlewares),
		uhttp.WithRequiredGet(options.ListRequiredGet),
		uhttp.WithOptionalGet(options.ListOptionalGet),
		uhttp.WithGet(func(r *http.Request, ret *int) interface{} {
			// Sanity check: ListOthersPermission can only be set if ListPermission is set
			if options.ListPermission == nil && options.ListOthersPermission != nil {
				*ret = http.StatusInternalServerError
				return map[string]string{"err": "Configuration problem: ListOthersPermission can only be set if ListPermission is set."}
			}

			// Check permissions
			var limitToUser *uauth.User
			var tmpUser *uauth.User
			var err error
			if options.ListPermission != nil {
				// Return nothing, if listPermission is required but the user does not have it
				tmpUser, err = uauth.UserFromRequest(r)
				if err != nil {
					return fmt.Errorf("Could not get user (%s)", err)
				}

				if !tmpUser.CheckPermission(*options.ListPermission) {
					return fmt.Errorf("User does not have the required permission: %s", *options.ListPermission)
				}

				// Limit results if ListOthersPermission is required but the user does not have it
				if options.ListOthersPermission != nil {
					if !tmpUser.CheckPermission(*options.ListOthersPermission) {
						limitToUser = tmpUser
					}
				}
			}

			// Load
			objsFromDb, err := options.ModelService.List(limitToUser != nil, r.Context())
			if err != nil {
				return err
			}

			// Render Response
			return objsFromDb
		}),
	)
}
