package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

// Returns an instance of an get-handler for the configured options
func genericGetHandler(options CrudOptions) uhttp.Handler {
	var middleware *uhttp.Middleware
	if options.GetPermission != nil {
		middleware = uauth.AuthJWT()
	}

	requiredGet := options.GetRequiredGet
	if requiredGet == nil {
		requiredGet = uhttp.R{}
	}
	requiredGet[options.IDParameterName] = uhttp.STRING

	return uhttp.Handler{
		PreProcess:    options.GetPreprocess,
		AddMiddleware: middleware,
		RequiredGet:   requiredGet,
		OptionalGet:   options.GetOptionalGet,
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Sanity check: GetOthersPermission can only be set if GetPermission is set
			if options.GetPermission == nil && options.GetOthersPermission != nil {
				uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: GetOthersPermission can only be set if GetPermission is set.")
				return
			}

			// Check permissions
			var limitToUser *uauth.User
			var tmpUser *uauth.User
			var err error
			if options.GetPermission != nil {
				tmpUser, err = uauth.UserFromRequest(r)
				if err != nil {
					uhttp.RenderError(w, r, fmt.Errorf("Could not get user (%s)", err))
					return
				}

				// Return nothing, if listPermission is required but the user does not have it
				if !tmpUser.CheckPermission(*options.GetPermission) {
					uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.GetPermission))
					return
				}

				// Limit results if ListOthersPermission is required but the user does not have it
				if options.GetOthersPermission != nil {
					if !tmpUser.CheckPermission(*options.GetOthersPermission) {
						limitToUser = tmpUser
					}
				}
			}

			// Get
			objectID := uhttp.GetAsString(options.IDParameterName, r)
			var objFromDb interface{}
			objFromDb, err = options.ModelService.Get(*objectID, limitToUser != nil, r.Context())

			if err != nil {
				uhttp.RenderError(w, r, fmt.Errorf("Could not find object with ID: '%s'", *objectID))
				return
			}

			uhttp.Render(w, r, objFromDb)
		}),
	}
}
