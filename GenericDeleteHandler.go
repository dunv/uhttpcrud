package uhttpcrud

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
	"github.com/dunv/umongo"
)

func genericDeleteHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get User
		user := r.Context().Value(uauth.CtxKeyUser).(uauth.User)
		if options.DeletePermission != nil && !user.CheckPermission(*options.DeletePermission) {
			uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.DeletePermission))
			return
		}

		// Get Params
		params := r.Context().Value(uhttp.CtxKeyParams).(map[string]interface{})

		// GetDB
		db := r.Context().Value(uhttp.CtxKeyDB).(*umongo.DbSession)
		service := options.ModelService.CopyAndInit(db)

		// Delete
		err := service.Delete(bson.ObjectIdHex(params[options.IDParameterName].(string)))
		if err != nil {
			uhttp.RenderError(w, r, err)
			return
		}

		// Answer
		uhttp.RenderMessageWithStatusCode(w, r, 200, "Deleted successfully")
	})
}

// GenericDeleteHandler <-
func GenericDeleteHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		Methods:      []string{"OPTIONS", "DELETE"},
		Handler:      genericDeleteHandler(options),
		DbRequired:   true,
		AuthRequired: true, // We need a user in order to delete an object
		RequiredParams: uhttp.Params{ParamMap: map[string]uhttp.ParamRequirement{
			options.IDParameterName: uhttp.ParamRequirement{AllValues: true},
		}},
	}
}
