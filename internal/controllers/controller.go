package controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/openshift/osin"

	"github.com/bartmika/osin-example/internal/models"
	"github.com/bartmika/osin-example/internal/session"
	"github.com/bartmika/osin-example/internal/utils"
)

type Controller struct {
	SecretSigningKeyBin []byte
	OAuthServer         *osin.Server
	OAuthStorage        *OSINRedisStorage
	TenantRepo          models.TenantRepository
	UserRepo            models.UserRepository
	SessionManager      *session.SessionManager
}

func (h *Controller) HandleRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get our URL paths which are slash-seperated.
	ctx := r.Context()
	p := ctx.Value("url_split").([]string)
	n := len(p)
	log.Println("HandleRequests|P, N:", p, n) // For debugging purposes only.

	_ = utils.DumpRequest(os.Stdout, "HandleRequests|dump", r) // For debugging purposes only.

	switch {
	// --- GATEWAY & PROFILE & DASHBOARD --- //
	case n == 3 && p[2] == "login" && r.Method == http.MethodPost:
		h.loginEndpoint(w, r)
	case n == 3 && p[2] == "register" && r.Method == http.MethodPost:
		h.registerEndpoint(w, r)
	case n == 3 && p[2] == "refresh-token" && r.Method == http.MethodPost:
		h.postRefreshToken(w, r)

	// --- TENANT --- //
	case n == 4 && p[2] == "tenant" && r.Method == http.MethodGet:
		h.tenantGetEndpoint(w, r, p[3])

	// --- OAuth 2.0 --- //
	case n == 1 && p[0] == "token":
		h.handleTokenRequest(w, r)
	case n == 1 && p[0] == "authorize":
		h.handleAuthorizeRequest(w, r)
	case n == 1 && p[0] == "info":
		h.handleInfoRequest(w, r)
	case n == 2 && p[0] == "appauth" && p[1] == "code":
		h.handleAuthorizationPermissionRequest(w, r)

	// --- CATCH ALL: D.N.E. ---
	default:
		http.NotFound(w, r)
	}
}
