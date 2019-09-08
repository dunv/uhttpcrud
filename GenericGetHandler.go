package uhttpcrud

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	uauthConfig "github.com/dunv/uauth/config"
	uauthModels "github.com/dunv/uauth/models"
	"github.com/dunv/uhttp"
	uhttpModels "github.com/dunv/uhttp/models"
	uhttpContextKeys "github.com/dunv/uhttp/contextkeys"
	"github.com/dunv/ulog"
)

func genericGetHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sanity check: GetOthersPermission can only be set if GetPermission is set
		if options.GetPermission == nil && options.GetOthersPermission != nil {
			uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: GetOthersPermission can only be set if GetPermission is set.")
			return
		}

		// Check permissions
		var limitToUser *uauthModels.User
		var tmpUser uauthModels.User
		if options.GetPermission != nil {
			tmpUser = r.Context().Value(uauthConfig.CtxKeyUser).(uauthModels.User)

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

		// Get Params
		params := r.Context().Value(uhttpContextKeys.CtxKeyParams).(map[string]interface{})

		// GetDB
		db := r.Context().Value(dbContextKey).(*mongo.Client)
		service := options.ModelService.CopyAndInit(db, options.Database)

		// Get
		objectID := params[options.IDParameterName]
		var objFromDb interface{}
		objFromDb, err := service.Get(objectID, &tmpUser, limitToUser != nil)

		if err != nil {
			uhttp.RenderError(w, r, fmt.Errorf("Could not find object with ID: '%s'", params[options.IDParameterName].(string)))
			return
		}

		ulog.LogIfError(json.NewEncoder(w).Encode(objFromDb))
		return
	})
}

// Returns an instance of an get-handler for the configured options
func GenericGetHandler(options CrudOptions) uhttpModels.Handler {
	return uhttpModels.Handler{
		GetHandler:                genericGetHandler(options),
		PreProcess:                options.GetPreprocess,
		AdditionalContextRequired: []uhttpModels.ContextKey{dbContextKey},
		AuthRequired:              options.GetPermission != nil,
		RequiredParams: uhttpModels.Params{ParamMap: map[string]uhttpModels.ParamRequirement{
			options.IDParameterName: uhttpModels.ParamRequirement{AllValues: true},
		}},
	}
}
