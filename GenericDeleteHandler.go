package uhttpcrud

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dunv/uauth"
	uauthModels "github.com/dunv/uauth/models"
	"github.com/dunv/uhttp"
)

// Returns an instance of an delete-handler for the configured options
func genericDeleteHandler(options CrudOptions) uhttp.Handler {
	requiredGet := options.DeleteRequiredGet
	if requiredGet == nil {
		requiredGet = uhttp.R{}
	}
	requiredGet[options.IDParameterName] = uhttp.STRING

	return uhttp.Handler{
		PreProcess:    options.DeletePreprocess,
		AddMiddleware: uauth.AuthJWT(), // We need a user in order to delete an object
		RequiredGet:   requiredGet,
		OptionalGet:   options.DeleteOptionalGet,
		DeleteHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Sanity check: DeleteOthersPermission can only be set if DeletePermission is set
			if options.DeletePermission == nil && options.DeleteOthersPermission != nil {
				uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: DeleteOthersPermission can only be set if DeletePermission is set.")
				return
			}

			// Check permissions
			var user uauthModels.User
			var limitToUser *uauthModels.User
			if options.DeletePermission != nil {
				// Return nothing, if deletePermission is required but the user does not have it
				user = uauth.User(r)
				if !user.CheckPermission(*options.DeletePermission) {
					uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.DeletePermission))
					return
				}

				// Limit results if DeleteOthersPermission is required but the user does not have it
				if options.DeleteOthersPermission != nil {
					if !user.CheckPermission(*options.DeleteOthersPermission) {
						limitToUser = &user
					}
				}
			}

			// GetDB
			db := r.Context().Value(dbContextKey).(*mongo.Client)
			service := options.ModelService.CopyAndInit(db, options.Database)
			objectID := uhttp.GetAsString(options.IDParameterName, r)

			// Delete
			err := service.Delete(*objectID, &user, limitToUser != nil, r.Context())
			if err != nil {
				uhttp.RenderError(w, r, err)
				return
			}

			// Answer
			uhttp.RenderMessageWithStatusCode(w, r, 200, "Deleted successfully")
		}),
	}
}
