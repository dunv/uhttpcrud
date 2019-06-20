package uhttpcrud

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

func genericGetHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sanity check: GetOthersPermission can only be set if GetPermission is set
		if options.GetPermission == nil && options.GetOthersPermission != nil {
			uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: GetOthersPermission can only be set if GetPermission is set.", nil)
			return
		}

		// Check permissions
		var limitToUser *uauth.User
		if options.GetPermission != nil {
			tmpUser := r.Context().Value(uauth.CtxKeyUser).(uauth.User)

			// Return nothing, if listPermission is required but the user does not have it
			if !tmpUser.CheckPermission(*options.GetPermission) {
				uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.GetPermission), nil)
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
		params := r.Context().Value(uhttp.CtxKeyParams).(map[string]interface{})

		// GetDB
		db := r.Context().Value(dbContextKey).(*mongo.Client)
		service := options.ModelService.CopyAndInit(db, options.Database)

		// Get
		objectID, err := primitive.ObjectIDFromHex(params[options.IDParameterName].(string))
		if err != nil {
			uhttp.RenderError(w, r, fmt.Errorf("Could not parse ID: '%s'", params[options.IDParameterName].(string)), nil)
			return
		}
		var objFromDb interface{}
		if limitToUser != nil { // This user obj will be != nil if GetOthersPermission is required, but the user does not have it
			objFromDb, err = service.Get(&objectID, limitToUser)
		} else {
			objFromDb, err = service.Get(&objectID, nil)
		}

		if err != nil {
			uhttp.RenderError(w, r, fmt.Errorf("Could not find object with ID: '%s'", params[options.IDParameterName].(string)), nil)
			return
		}

		json.NewEncoder(w).Encode(objFromDb)
		return
	})
}

// GenericGetHandler <-
func GenericGetHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		GetHandler:   genericGetHandler(options),
		DbRequired:   []uhttp.ContextKey{dbContextKey},
		AuthRequired: options.GetPermission != nil,
		RequiredParams: uhttp.Params{ParamMap: map[string]uhttp.ParamRequirement{
			options.IDParameterName: uhttp.ParamRequirement{AllValues: true},
		}},
	}
}
