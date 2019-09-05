package uhttpcrud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
	"go.mongodb.org/mongo-driver/mongo"
)

func genericUpdateHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sanity check: UpdateOthersPermission can only be set if UpdatePermission is set
		if options.UpdatePermission == nil && options.UpdateOthersPermission != nil {
			uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: UpdateOthersPermission can only be set if UpdatePermission is set.")
			return
		}

		// Check permissions
		var user uauth.User
		var limitToUser *uauth.User
		if options.UpdatePermission != nil {
			user = r.Context().Value(uauth.CtxKeyUser).(uauth.User)

			// Return nothing, if updatePermission is required but the user does not have it
			if !user.CheckPermission(*options.UpdatePermission) {
				uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.UpdatePermission))
				return
			}

			// Limit results if UpdateOthersPermission is required but the user does not have it
			if options.UpdateOthersPermission != nil {
				if !user.CheckPermission(*options.UpdateOthersPermission) {
					limitToUser = &user
				}
			}
		}

		// Parse body into new "dynamic" object
		model := options.Model
		reflectModel := reflect.New(reflect.TypeOf(model))
		modelInterface := reflectModel.Interface()
		err := json.NewDecoder(r.Body).Decode(modelInterface)
		if err != nil {
			// uhttp.RenderError(w, r, fmt.Errorf("Could not decode request body"))
			uhttp.RenderError(w, r, err)
			return
		}

		// Get object from db
		db := r.Context().Value(dbContextKey).(*mongo.Client)
		service := options.ModelService.CopyAndInit(db, options.Database)

		// Check if all required populated fields are populated (indexes)
		idFromModel := modelInterface.(WithID).GetID()
		if idFromModel == nil || !service.Validate(modelInterface) {
			uhttp.RenderError(w, r, fmt.Errorf("Non-nullable properties are null or no ID present"))
			return
		}

		// Check if already exists
		if limitToUser != nil {
			_, err = service.Get(idFromModel, limitToUser)
		} else {
			_, err = service.Get(idFromModel, nil)
		}
		if err != nil {
			uhttp.RenderError(w, r, fmt.Errorf("No object with the id %s exists", modelInterface.(WithID).GetID()))
			return
		}

		// Actual update
		updatedDocument, err := service.Update(modelInterface, &user)
		if err != nil {
			uhttp.RenderError(w, r, err)
			return
		}

		// Answer
		uhttp.CheckAndLogError(json.NewEncoder(w).Encode(updatedDocument))
	})
}

func GenericUpdateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		PostHandler:  genericUpdateHandler(options),
		PostModel:    options.Model,
		PreProcess:   options.UpdatePreprocess,
		DbRequired:   []uhttp.ContextKey{dbContextKey},
		AuthRequired: true, // We need a user in order to update an object
	}
}
