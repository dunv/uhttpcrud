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

func genericCreateHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sanity check: CreateOthersPermission can only be set if CreatePermission is set
		if options.CreatePermission == nil && options.CreateOthersPermission != nil {
			uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: CreateOthersPermission can only be set if CreatePermission is set.")
			return
		}

		// Get User
		user := r.Context().Value(uauth.CtxKeyUser).(uauth.User)
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
		db := r.Context().Value(dbContextKey).(*mongo.Client)
		service := options.ModelService.CopyAndInit(db, options.Database)

		// Check if all required populated fields are populated (indexes)
		if !service.CheckNotNullable(modelInterface) {
			uhttp.RenderError(w, r, fmt.Errorf("Non-nullable properties are null"))
			return
		}

		// Create (will return an error if already exists)
		ID, err := service.Create(modelInterface, user)
		if err != nil {
			uhttp.RenderError(w, r, err)
			return
		}

		// Answer
		responseModel := map[string]interface{}{
			"id": ID,
		}
		uhttp.CheckAndLogError(json.NewEncoder(w).Encode(responseModel))
	})
}

func GenericCreateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		PostHandler:  genericCreateHandler(options),
		PreProcess:   options.CreatePreprocess,
		DbRequired:   []uhttp.ContextKey{dbContextKey},
		AuthRequired: true, // We need a user in order to create an object
	}
}
