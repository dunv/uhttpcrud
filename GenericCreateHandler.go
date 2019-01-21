package uhttpcrud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/dunv/uauth"
	"github.com/dunv/umongo"
	"github.com/dunv/uhttp"
)

func genericCreateHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get User
		user := r.Context().Value(auth.CtxKeyUser).(auth.User)
		if options.CreatePermission != nil && !user.CheckPermission(*options.CreatePermission) {
			uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.CreatePermission))
			return
		}

		// Parse body into new "dynamic" object
		model := options.Model
		reflectModel := reflect.New(reflect.TypeOf(model))
		modelInterface := reflectModel.Interface()
		err := json.NewDecoder(r.Body).Decode(modelInterface)
		if err != nil {
			uhttp.RenderError(w, r, fmt.Errorf("Could not decode request body"))
			return
		}

		// Get object from db
		db := r.Context().Value(uhttp.CtxKeyDB).(*mongo.DbSession)
		service := options.ModelService.CopyAndInit(db)

		// Check if all required populated fields are populated (indexes)
		if !service.CheckNotNullable(modelInterface) {
			uhttp.RenderError(w, r, fmt.Errorf("Non-nullable properties are null"))
			return
		}

		// Create (will return an error if already exists)
		err = service.Create(modelInterface, user)
		if err != nil {
			uhttp.RenderError(w, r, err)
			return
		}

		// Answer
		uhttp.RenderMessageWithStatusCode(w, r, 200, "Saved successfully")
	})
}

// GenericCreateHandler <-
func GenericCreateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		Methods:      []string{"OPTIONS", "POST"},
		Handler:      genericCreateHandler(options),
		DbRequired:   true,
		AuthRequired: true, // We need a user in order to create an object
	}
}
