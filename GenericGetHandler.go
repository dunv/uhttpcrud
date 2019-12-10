package uhttpcrud

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dunv/uauth"
	uauthModels "github.com/dunv/uauth/models"
	"github.com/dunv/uhttp"
	"github.com/dunv/uhttp/params"
)

// Returns an instance of an get-handler for the configured options
func genericGetHandler(options CrudOptions) uhttp.Handler {
	var middleware *uhttp.Middleware
	if options.GetPermission != nil {
		middleware = uauth.AuthJWT()
	}

	requiredGet := options.GetRequiredGet
	if requiredGet == nil {
		requiredGet = params.R{}
	}
	requiredGet[options.IDParameterName] = params.STRING

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
			var limitToUser *uauthModels.User
			var tmpUser uauthModels.User
			if options.GetPermission != nil {
				tmpUser = uauth.User(r)

				// Return nothing, if listPermission is required but the user does not have it
				if !tmpUser.CheckPermission(*options.GetPermission) {
					uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.GetPermission))
					return
				}

				// Limit results if ListOthersPermission is required but the user does not have it
				if options.GetOthersPermission != nil {
					if !tmpUser.CheckPermission(*options.GetOthersPermission) {
						limitToUser = &tmpUser
					}
				}
			}

			// GetDB
			db := r.Context().Value(dbContextKey).(*mongo.Client)
			service := options.ModelService.CopyAndInit(db, options.Database)

			// Get
			objectID := params.GetAsString(options.IDParameterName, r)
			var objFromDb interface{}
			objFromDb, err := service.Get(*objectID, &tmpUser, limitToUser != nil, r.Context())

			if err != nil {
				uhttp.RenderError(w, r, fmt.Errorf("Could not find object with ID: '%s'", *objectID))
				return
			}

			uhttp.Render(w, r, objFromDb)
		}),
	}
}
