package controllers

import (
	"context"
	"log"
	"net/http"

	"github.com/bartmika/osin-example/internal/models"
	"github.com/openshift/osin"
)

func (h *Controller) OSINProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r.ParseForm()

		// DEVELOPERS NOTE:
		// The following code was written with help from the following link:
		// https://github.com/openshift/osin/blob/e57c5735994989e8c4c95d2ddf8b90e5b4345de3/info.go#L16

		// STEP 1:
		// Get the token or flow to the next middleware.
		bearer := osin.CheckBearerAuth(r)
		if bearer == nil {
			log.Println("OSINProcessorMiddleware|Skipping|No bearer, flowing to the next middleware.")
			// Flow to the next middleware without anything done.
			ctx = context.WithValue(ctx, "is_authorized", false)
			fn(w, r.WithContext(ctx))
			return
		}
		log.Println("OSINProcessorMiddleware|Processing...")

		if bearer.Code == "" {
			http.Error(w, "code is nil", http.StatusUnauthorized)
			return
		}

		// load access data
		accessData, err := h.OAuthStorage.LoadAccess(bearer.Code)

		if accessData == nil {
			if err != nil {
				http.Error(w, "failed to load access data", http.StatusUnauthorized)
				return
			}

			http.Error(w, "access data is nil", http.StatusUnauthorized)
			return
		}

		if accessData.UserData == nil {
			http.Error(w, "user data is nil", http.StatusUnauthorized)
			return
		}

		// Extract the values and set data type.
		user := accessData.UserData.(models.UserLite)

		// Save to context.
		ctx = context.WithValue(ctx, "is_authorized", true)
		ctx = context.WithValue(ctx, "user_tenant_id", user.TenantID)
		ctx = context.WithValue(ctx, "user_role_id", user.RoleID)
		ctx = context.WithValue(ctx, "user_id", user.ID)
		ctx = context.WithValue(ctx, "user_uuid", user.UUID)
		ctx = context.WithValue(ctx, "user_timezone", user.Timezone)
		ctx = context.WithValue(ctx, "user_name", user.Name)

		// Flow to the next middleware without anything done.
		fn(w, r.WithContext(ctx))
	}
}
