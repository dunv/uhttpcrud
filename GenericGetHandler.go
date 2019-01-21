package uhttpcrud

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/dunv/auth"
	"github.com/dunv/mongo"
	"github.com/dunv/uhttp"
)

func genericGetHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if options.GetPermission != nil {
			// Get User
			user := r.Context().Value(auth.CtxKeyUser).(auth.User)
			if !user.CheckPermission(*options.GetPermission) {
				uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.ListPermission))
				return
			}
		}

		// Get Params
		params := r.Context().Value(uhttp.CtxKeyParams).(map[string]interface{})

		// GetDB
		db := r.Context().Value(uhttp.CtxKeyDB).(*mongo.DbSession)
		service := options.ModelService.CopyAndInit(db)

		// Check if already exists
		objectID := bson.ObjectIdHex(params[options.IDParameterName].(string))
		objectFromDb, err := service.Get(&objectID)
		if err != nil {
			uhttp.RenderError(w, r, fmt.Errorf("Could not find object with ID: '%s'", params[options.IDParameterName].(string)))
			return
		}

		json.NewEncoder(w).Encode(objectFromDb)
		return
	})
}

// GenericGetHandler <-
func GenericGetHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		Methods:      []string{"GET"},
		Handler:      genericGetHandler(options),
		DbRequired:   true,
		AuthRequired: options.GetPermission != nil,
		RequiredParams: uhttp.Params{ParamMap: map[string]uhttp.ParamRequirement{
			options.IDParameterName: uhttp.ParamRequirement{AllValues: true},
		}},
	}
}
