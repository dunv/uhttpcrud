package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

func genericUpdateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		PostModel:     options.Model,
		PreProcess:    options.UpdatePreprocess,
		RequiredGet:   options.UpdateRequiredGet,
		OptionalGet:   options.UpdateOptionalGet,
		AddMiddleware: uauth.AuthJWT(), // We need a user in order to update an object
		PostHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Sanity check: UpdateOthersPermission can only be set if UpdatePermission is set
			if options.UpdatePermission == nil && options.UpdateOthersPermission != nil {
				uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: UpdateOthersPermission can only be set if UpdatePermission is set.")
				return
			}

			// Check permissions
			var user *uauth.User
			var limitToUser *uauth.User
			var err error
			if options.UpdatePermission != nil {
				user, err = uauth.UserFromRequest(r)
				if err != nil {
					uhttp.RenderError(w, r, fmt.Errorf("Could not get user (%s)", err))
					return
				}

				// Return nothing, if updatePermission is required but the user does not have it
				if !user.CheckPermission(*options.UpdatePermission) {
					uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.UpdatePermission))
					return
				}

				// Limit results if UpdateOthersPermission is required but the user does not have it
				if options.UpdateOthersPermission != nil {
					if !user.CheckPermission(*options.UpdateOthersPermission) {
						limitToUser = user
					}
				}
			}

			modelInterface := getWithIDFromPostModel(r)

			// Check if all required populated fields are populated (indexes)
			idFromModel, err := modelInterface.(WithID).GetID()
			if err != nil {
				uhttp.RenderError(w, r, fmt.Errorf("could not getID (%s)", err))
				return
			}

			if idFromModel == "" || !options.ModelService.Validate(modelInterface) {
				uhttp.RenderError(w, r, fmt.Errorf("Non-nullable properties are null or no ID present"))
				return
			}

			// Check if already exists
			_, err = options.ModelService.Get(idFromModel, limitToUser != nil, r.Context())
			if err != nil {
				uhttp.RenderError(w, r, fmt.Errorf("No object with the id %s exists (%s)", idFromModel, err))
				return
			}

			// Actual update
			updatedDocument, err := options.ModelService.Update(modelInterface, limitToUser != nil, r.Context())
			if err != nil {
				uhttp.RenderError(w, r, err)
				return
			}

			// Answer
			uhttp.Render(w, r, updatedDocument)
		}),
	}
}
