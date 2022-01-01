package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bartmika/osin-example/internal/models"
	"github.com/bartmika/osin-example/internal/utils"
	"github.com/openshift/osin"
)

// Access token endpoint
func (h *Controller) handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("handleTokenRequest|starting...")
	resp := h.OAuthServer.NewResponse()
	defer resp.Close()

	var authenticatedUser *models.User

	if ar := h.OAuthServer.HandleAccessRequest(resp, r); ar != nil {
		log.Println("handleTokenRequest|HandleAccessRequest|starting...")

		switch ar.Type {
		case osin.AUTHORIZATION_CODE:

			//TODO: IMPL.
			ar.Authorized = true

		case osin.REFRESH_TOKEN:
			log.Println("handleTokenRequest|HandleAccessRequest|REFRESH_TOKEN|Starting")
			// DEVELOPERS NOTE:
			// We are only going to authorize the `refresh token` if there is
			// user data associated with the request. Once the condition is
			// true then we will save the authenticated user.
			if ar.UserData != nil && ar.AccessData != nil {
				ud := ar.UserData.(models.UserLite)

				// The `osin` library requires we set this value to true so
				// the whole system will generate our access token.
				ar.Authorized = true

				// Store this for our output later...
				authenticatedUser = &models.User{
					ID:       ud.ID,
					UUID:     ud.UUID,
					TenantID: ud.TenantID,
					Name:     ud.Name,
					State:    ud.State,
					RoleID:   ud.RoleID,
					Timezone: ud.Timezone,
					Language: ud.Language,
				}
			}
			log.Println("handleTokenRequest|HandleAccessRequest|REFRESH_TOKEN|Finished")
		case osin.PASSWORD:
			log.Println("handleTokenRequest|HandleAccessRequest|PASSWORD|Started")
			// DEVELOPERS NOTE:
			// THIS IS WHERE WE WANT TO HANDLE THE USER LOGIC FROM OUR SYSTEM.
			user, _ := h.authenticatedUser(context.Background(), ar.Username, ar.Password)
			if user != nil {
				// The `osin` library requires we set this value to true so
				// the whole system will generate our access token.
				ar.Authorized = true

				// Store this for our output later...
				authenticatedUser = user
			}
			log.Println("handleTokenRequest|ar.Authorized:", ar.Authorized)
			log.Println("handleTokenRequest|HandleAccessRequest|PASSWORD|Finished")
		case osin.CLIENT_CREDENTIALS:
			log.Println("handleTokenRequest|HandleAccessRequest|CLIENT_CREDENTIALS|Started")

			if ar.UserData != nil {
				ud := ar.UserData.(models.UserLite)

				// The `osin` library requires we set this value to true so
				// the whole system will generate our access token.
				ar.Authorized = true

				// Store this for our output later...
				authenticatedUser = &models.User{
					ID:       ud.ID,
					UUID:     ud.UUID,
					TenantID: ud.TenantID,
					Name:     ud.Name,
					State:    ud.State,
					RoleID:   ud.RoleID,
					Timezone: ud.Timezone,
					Language: ud.Language,
				}
				log.Println("handleTokenRequest|HandleAccessRequest|set user data")
			} else {
				log.Println("handleTokenRequest|HandleAccessRequest|no user data set")
			}
			log.Println("handleTokenRequest|HandleAccessRequest|CLIENT_CREDENTIALS|Finished")
		case osin.ASSERTION:
			//TODO: IMPL.
			if ar.AssertionType == "urn:osin.example.complete" && ar.Assertion == "osin.data" {
				ar.Authorized = true
			}
		}

		// This is important because we want to store the user account
		// in the store session for the particular access token when.
		// This means that when we retrieve from the oauth storage for
		// the access token, the user profile will be returned as well.
		// That also means we can utilize our oauth storage in other
		// areas of our app like 'middleware' to protect resources.
		if authenticatedUser != nil {
			ar.UserData = &models.UserLite{
				TenantID: authenticatedUser.TenantID,
				RoleID:   authenticatedUser.RoleID,
				ID:       authenticatedUser.ID,
				UUID:     authenticatedUser.UUID,
				Timezone: authenticatedUser.Timezone,
				Name:     authenticatedUser.Name,
				State:    authenticatedUser.State,
				Email:    authenticatedUser.Email,
			}
		}

		h.OAuthServer.FinishAccessRequest(resp, r, ar)

		log.Println("handleTokenRequest|HandleAccessRequest|finished")
	}
	if resp.IsError && resp.InternalError != nil {
		log.Printf("handleTokenRequest|HandleAccessRequest|ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		// The following section is the additional `extra` data we want to
		// return in the token object when we return our token.
		if resp.Output != nil && authenticatedUser != nil {
			resp.Output["user_id"] = authenticatedUser.ID
			resp.Output["first_name"] = authenticatedUser.FirstName
			resp.Output["last_name"] = authenticatedUser.LastName
			resp.Output["email"] = authenticatedUser.Email
			resp.Output["language"] = authenticatedUser.Language
			resp.Output["role_id"] = authenticatedUser.RoleID
			resp.Output["tenant_id"] = authenticatedUser.TenantID
			log.Println("handleTokenRequest|added extra to token")
		}

	}
	osin.OutputJSON(resp, w, r)
	log.Println("handleTokenRequest|finished")
}

// Authorization code endpoint
func (h *Controller) handleAuthorizeRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("handleAuthorizeRequest|starting...")

	resp := h.OAuthServer.NewResponse()
	defer resp.Close()

	if ar := h.OAuthServer.HandleAuthorizeRequest(resp, r); ar != nil {
		if !utils.HandleLoginPage(ar, w, r) {
			return
		}
		ar.UserData = struct{ Login string }{Login: "test"}
		ar.Authorized = true
		h.OAuthServer.FinishAuthorizeRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		resp.Output["custom_parameter"] = 187723
	}
	osin.OutputJSON(resp, w, r)

	log.Println("handleAuthorizeRequest|finished")
}

// Info endpoint
func (h *Controller) handleInfoRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("handleInfoRequest|starting")
	resp := h.OAuthServer.NewResponse()
	defer resp.Close()

	if ir := h.OAuthServer.HandleInfoRequest(resp, r); ir != nil {
		h.OAuthServer.FinishInfoRequest(resp, r, ir)
	}
	osin.OutputJSON(resp, w, r)

	log.Println("handleInfoRequest|finished")
}
