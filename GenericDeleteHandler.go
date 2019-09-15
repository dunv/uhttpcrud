package uhttpcrud

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dunv/uauth"
	uauthModels "github.com/dunv/uauth/models"
	"github.com/dunv/uhttp"
	uhttpModels "github.com/dunv/uhttp/models"
	"github.com/dunv/uhttp/params"
)

// Returns an instance of an delete-handler for the configured options
func genericDeleteHandler(options CrudOptions) uhttpModels.Handler {
	return uhttpModels.Handler{
		PreProcess:   options.DeletePreprocess,
		AuthRequired: true, // We need a user in order to delete an object
		RequiredGet: params.R{
			options.IDParameterName: params.STRING,
		},
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
			objectID := params.GetAsString(options.IDParameterName, r)

			// Delete
			err := service.Delete(*objectID, &user, limitToUser != nil)
			if err != nil {
				uhttp.RenderError(w, r, err)
				return
			}

			// Answer
			uhttp.RenderMessageWithStatusCode(w, r, 200, "Deleted successfully")
		}),
	}
}
