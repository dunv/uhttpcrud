package uhttpcrud

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dunv/uauth"
	uauthModels "github.com/dunv/uauth/models"
	"github.com/dunv/uhttp"
	uhttpModels "github.com/dunv/uhttp/models"
	"github.com/dunv/ulog"
	"go.mongodb.org/mongo-driver/mongo"
)

// Returns an instance of an list-handler for the configured options
func genericListHandler(options CrudOptions) uhttpModels.Handler {
	var middleware *uhttpModels.Middleware
	if options.ListPermission != nil {
		middleware = uauth.AuthJWT()
	}
	return uhttpModels.Handler{
		PreProcess:    options.ListPreprocess,
		AddMiddleware: middleware,
		GetHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Sanity check: ListOthersPermission can only be set if ListPermission is set
			if options.ListPermission == nil && options.ListOthersPermission != nil {
				uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: ListOthersPermission can only be set if ListPermission is set.")
				return
			}

			// Check permissions
			var limitToUser *uauthModels.User
			var tmpUser uauthModels.User
			if options.ListPermission != nil {
				// Return nothing, if listPermission is required but the user does not have it
				tmpUser = uauth.User(r)
				if !tmpUser.CheckPermission(*options.ListPermission) {
					uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.ListPermission))
					return
				}

				// Limit results if ListOthersPermission is required but the user does not have it
				if options.ListOthersPermission != nil {
					if !tmpUser.CheckPermission(*options.ListOthersPermission) {
						limitToUser = &tmpUser
					}
				}
			}

			// GetDB
			db := r.Context().Value(dbContextKey).(*mongo.Client)
			service := options.ModelService.CopyAndInit(db, options.Database)

			// Load
			objsFromDb, err := service.List(&tmpUser, limitToUser != nil)
			if err != nil {
				uhttp.RenderError(w, r, err)
				return
			}

			// Render Response
			ulog.LogIfError(json.NewEncoder(w).Encode(objsFromDb))
			return
		}),
	}
}
