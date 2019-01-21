package uhttpcrud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/dunv/auth"
	"github.com/dunv/mongo"
	"github.com/dunv/uhttp"
)

func genericUpdateHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get User
		user := r.Context().Value(auth.CtxKeyUser).(auth.User)
		if options.UpdatePermission != nil && !user.CheckPermission(*options.UpdatePermission) {
			uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.UpdatePermission))
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
		idFromModel := modelInterface.(WithID).GetID()
		if idFromModel == nil || !service.CheckNotNullable(modelInterface) {
			uhttp.RenderError(w, r, fmt.Errorf("Non-nullable properties are null or no ID present"))
			return
		}

		// Check if already exists
		_, err = service.Get(idFromModel)
		if err != nil {
			uhttp.RenderError(w, r, fmt.Errorf("No object with the id %s exists", modelInterface.(WithID).GetID()))
			return
		}

		// Actual update
		err = service.Update(modelInterface, user)
		if err != nil {
			uhttp.RenderError(w, r, err)
			return
		}

		// Answer
		uhttp.RenderMessageWithStatusCode(w, r, 200, "Updated successfully")
	})
}

// GenericUpdateHandler <-
func GenericUpdateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		Methods:      []string{"OPTIONS", "POST"},
		Handler:      genericUpdateHandler(options),
		DbRequired:   true,
		AuthRequired: true, // We need a user in order to update an object
	}
}
