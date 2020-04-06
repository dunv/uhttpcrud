package uhttpcrud

import (
	"fmt"
	"net/http"

	"github.com/dunv/uauth"
	"github.com/dunv/uhttp"
)

// Returns an instance of an update-handler for the configured options
func genericCreateHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		PostModel:     options.Model,
		PreProcess:    options.CreatePreprocess,
		RequiredGet:   options.CreateRequiredGet,
		OptionalGet:   options.CreateOptionalGet,
		AddMiddleware: uauth.AuthJWT(), // We need a user in order to create an object
		PostHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Sanity check: CreateOthersPermission can only be set if CreatePermission is set
			if options.CreatePermission == nil && options.CreateOthersPermission != nil {
				uhttp.RenderMessageWithStatusCode(w, r, 500, "Configuration problem: CreateOthersPermission can only be set if CreatePermission is set.")
				return
			}

			// Get User
			user, err := uauth.UserFromRequest(r)
			if err != nil {
				uhttp.RenderError(w, r, fmt.Errorf("Could not get user (%s)", err))
				return
			}

			if options.CreatePermission != nil && !user.CheckPermission(*options.CreatePermission) {
				uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.CreatePermission))
				return
			}

			modelInterface := getWithIDFromPostModel(r)

			// Check if all required populated fields are populated (indexes)
			if !options.ModelService.Validate(modelInterface) {
				uhttp.RenderError(w, r, fmt.Errorf("Non-nullable properties are null"))
				return
			}

			// Create (will return an error if already exists)
			createdDocument, err := options.ModelService.Create(modelInterface, r.Context())
			if err != nil {
				uhttp.RenderError(w, r, err)
				return
			}

			// Answer
			uhttp.Render(w, r, createdDocument)
		}),
	}
}
