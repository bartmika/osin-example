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

		// The following code will lookup the URL path in a whitelist and
		// if the visited path matches then we will skip URL protection.
		// We do this because a majority of API endpoints are protected
		// by authorization.

		urlSplit := ctx.Value("url_split").([]string)
		skipPath := map[string]bool{
			"version":       true,
			"hello":         true,
			"register":      true,
			"login":         true,
			"refresh-token": true,
			"submit":        true,
			"appauth":       true,
		}

		// If the URL cannot be split into the size we want then skip running
		// this middleware. Then check to see if we our URL contains the white
		// listed URL and if it does then we skip protection for this request.

		ctx = context.WithValue(ctx, "is_skip_authentication", true)
		if len(urlSplit) <= 2 {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}
		if skipPath[urlSplit[2]] || skipPath[urlSplit[1]] {
			log.Println("OSINProcessorMiddleware|Skipping|Whitelisted path detected")
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}
		ctx = context.WithValue(ctx, "is_skip_authentication", false)

		// The following code was written with help from the following link:
		// https://github.com/openshift/osin/blob/e57c5735994989e8c4c95d2ddf8b90e5b4345de3/info.go#L16

		r.ParseForm()

		// STEP 1:
		// Get the token or flow to the next middleware.
		bearer := osin.CheckBearerAuth(r)
		if bearer == nil {
			log.Println("OSINProcessorMiddleware|Skipping|No bearer, flowing to the next middleware.")
			ctx = context.WithValue(ctx, "is_authorized", false)
			fn(w, r.WithContext(ctx)) // Flow to the next middleware without anything done.
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
