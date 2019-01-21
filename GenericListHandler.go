package uhttpcrud

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dunv/uhttp"
	"github.com/dunv/umongo"
)

func genericListHandler(options CrudOptions) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if options.ListPermission != nil {
			// Get User
			user := r.Context().Value(auth.CtxKeyUser).(auth.User)
			if !user.CheckPermission(*options.ListPermission) {
				uhttp.RenderError(w, r, fmt.Errorf("User does not have the required permission: %s", *options.ListPermission))
				return
			}
		}

		// GetDB
		db := r.Context().Value(uhttp.CtxKeyDB).(*umongo.DbSession)
		service := options.ModelService.CopyAndInit(db)

		// Load
		objsFromDb, err := service.List()
		if err != nil {
			uhttp.RenderError(w, r, err)
			return
		}

		json.NewEncoder(w).Encode(objsFromDb)
		return
	})
}

// GenericListHandler <-
func GenericListHandler(options CrudOptions) uhttp.Handler {
	return uhttp.Handler{
		Methods:      []string{"GET"},
		Handler:      genericListHandler(options),
		DbRequired:   true,
		AuthRequired: options.ListPermission != nil,
	}
}
