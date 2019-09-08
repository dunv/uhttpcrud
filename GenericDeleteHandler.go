package uhttpcrud

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	uauthConfig "github.com/dunv/uauth/config"
	uauthModels "github.com/dunv/uauth/models"
	"github.com/dunv/uhttp"
	uhttpContextKeys "github.com/dunv/uhttp/contextkeys"
	uhttpModels "github.com/dunv/uhttp/models"
)

func genericDeleteHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sanity check: DeleteOthersPermission can only be set if DeletePermission is set
		if options.DeletePermission == nil && options.DeleteOthersPermission != nil {
			uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: DeleteOthersPermission can only be set if DeletePermission is set.")
			return
		}

		// Check permissions
		var user uauthModels.User
		var limitToUser *uauthModels.User
		if options.DeletePermission != nil {
			user = r.Context().Value(uauthConfig.CtxKeyUser).(uauthModels.User)

			// Return nothing, if deletePermission is required but the user does not have it
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

		// Get Params
		params := r.Context().Value(uhttpContextKeys.CtxKeyParams).(map[string]interface{})

		// GetDB
		db := r.Context().Value(dbContextKey).(*mongo.Client)
		service := options.ModelService.CopyAndInit(db, options.Database)

		objectID := params[options.IDParameterName]

		// Delete
		err := service.Delete(objectID, &user, limitToUser != nil)
		if err != nil {
			uhttp.RenderError(w, r, err)
			return
		}

		// Answer
		uhttp.RenderMessageWithStatusCode(w, r, 200, "Deleted successfully")
	})
}

// Returns an instance of an delete-handler for the configured options
func GenericDeleteHandler(options CrudOptions) uhttpModels.Handler {
	return uhttpModels.Handler{
		DeleteHandler:             genericDeleteHandler(options),
		PreProcess:                options.DeletePreprocess,
		AdditionalContextRequired: []uhttpModels.ContextKey{dbContextKey},
		AuthRequired:              true, // We need a user in order to delete an object
		RequiredParams: uhttpModels.Params{ParamMap: map[string]uhttpModels.ParamRequirement{
			options.IDParameterName: uhttpModels.ParamRequirement{AllValues: true},
		}},
	}
}
