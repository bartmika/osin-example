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
	Config                        map[string]string
	OAuthServer                   *osin.Server
	OAuthStorage                  *OSINRedisStorage
	TenantRepo                    models.TenantRepository
	UserRepo                      models.UserRepository
	ApplicationRepo               models.ApplicationRepository
	ApplicationLiteRepo           models.ApplicationLiteRepository
	AuthorizedApplicationRepo     models.AuthorizedApplicationRepository
	AuthorizedApplicationLiteRepo models.AuthorizedApplicationLiteRepository
	SessionManager                *session.SessionManager
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

	// --- APPLICATION --- //
	case n == 3 && p[2] == "applications" && r.Method == http.MethodGet:
		h.applicationsListEndpoint(w, r)
	case n == 3 && p[2] == "applications" && r.Method == http.MethodPost:
		h.applicationCreateEndpoint(w, r)
	case n == 4 && p[2] == "application" && r.Method == http.MethodGet:
		h.applicationGetEndpoint(w, r, p[3])
	case n == 4 && p[2] == "application" && r.Method == http.MethodPut:
		h.applicationUpdateEndpoint(w, r, p[3])
	case n == 4 && p[2] == "application" && r.Method == http.MethodDelete:
		h.applicationDeleteEndpoint(w, r, p[3])

	// --- OAuth 2.0 --- //
	case n == 1 && p[0] == "token":
		h.handleTokenRequest(w, r)
	case n == 1 && p[0] == "authorize":
		h.handleAuthorizeRequest(w, r)
	case n == 1 && p[0] == "info":
		h.handleInfoRequest(w, r)

	// --- CATCH ALL: D.N.E. ---
	default:
		log.Println("HandleRequests|Page D.N.E.")
		http.NotFound(w, r)
	}
}
