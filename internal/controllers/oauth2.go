package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bartmika/osin-example/internal/models"
	"github.com/google/uuid"
	"github.com/openshift/osin"
)

// Access token endpoint
func (h *Controller) handleTokenRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("handleTokenRequest|Starting")
	ctx := r.Context()
	resp := h.OAuthServer.NewResponse()
	defer resp.Close()

	var authenticatedUser *models.User

	if ar := h.OAuthServer.HandleAccessRequest(resp, r); ar != nil {
		log.Println("handleTokenRequest|HandleAccessRequest|Starting")

		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			log.Println("handleTokenRequest|HandleAccessRequest|AUTHORIZATION_CODE|Starting")

			// ALGORITHM
			// (1) Check to see if authorize token in request has `UserData` to continue.
			// (2) Check to see if third-party application is running in our system (and was not shutdown by us) to continue.
			// (3) Check to see if user has granted permission to use third-party application in our system to continue.
			// (4) If all above true, grant token access.

			// DEVELOPER NOTE:
			// Our `authorize` endpoint includes the user data so we want to
			// confirm it was included in the authorization code and if it
			// was not then we error!
			if ar.UserData != nil {
				// FOR SECURITY PURPOSES, WE WANT TO CHECK TO SEE IF THE THIRD PARTY
				// APPLICATION IS RUNNING AND HAS NOT BEEN SHUTDOWN BY US OR THE
				// DEVELOPER TO ENSURE NO UNAUTHORIZED ACCESS OCCures.
				if app, _ := h.ApplicationRepo.GetByClientID(ctx, ar.Client.GetId()); app != nil {
					ud := ar.UserData.(models.UserLite)

					// FOR SECURITY PURPOSES, WE WANT TO CHECK THIS APPLICATION
					// HAS BEEN GRANTED ACCESS TO THE USER ACCOUNT.
					if permissionGranted, _ := h.AuthorizedApplicationRepo.CheckIfPermissionGrantedByUserIDAndByApplicationID(ctx, ud.ID, app.ID); permissionGranted {
						log.Println("handleTokenRequest|HandleAccessRequest|AUTHORIZATION_CODE|Authorized")
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
				}
			}

			log.Println("handleTokenRequest|HandleAccessRequest|AUTHORIZATION_CODE|Finished")

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
				log.Println("handleTokenRequest|HandleAccessRequest|REFRESH_TOKEN|Authorized")
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
			user, _ := h.authenticatedUser(ctx, ar.Username, ar.Password)
			if user != nil {
				// The `osin` library requires we set this value to true so
				// the whole system will generate our access token.
				log.Println("handleTokenRequest|HandleAccessRequest|PASSWORD|Authorized")
				ar.Authorized = true

				// Store this for our output later...
				authenticatedUser = user
			} else {
				log.Println("handleTokenRequest|HandleAccessRequest|PASSWORD|Unauthorized")
			}
			log.Println("handleTokenRequest|HandleAccessRequest|PASSWORD|Finished")
		case osin.CLIENT_CREDENTIALS:
			log.Println("handleTokenRequest|HandleAccessRequest|CLIENT_CREDENTIALS|Started")

			if ar.UserData != nil {
				ud := ar.UserData.(models.UserLite)

				// The `osin` library requires we set this value to true so
				// the whole system will generate our access token.
				log.Println("handleTokenRequest|HandleAccessRequest|CLIENT_CREDENTIALS|Authorized")
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
			log.Println("handleTokenRequest|authenticatedUser|ar.UserData|set")
		}

		h.OAuthServer.FinishAccessRequest(resp, r, ar)

		log.Println("handleTokenRequest|HandleAccessRequest|Finished")
	}
	if resp.IsError && resp.InternalError != nil {
		log.Printf("handleTokenRequest|HandleAccessRequest|ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		// The following section is the additional `extra` data we want to
		// return in the token object when we return our token.
		if resp.Output != nil && authenticatedUser != nil {
			resp.Output["user_id"] = authenticatedUser.ID
			resp.Output["user_uuid"] = authenticatedUser.UUID
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
	log.Println("handleTokenRequest|Finished")
}

// Authorization code endpoint
func (h *Controller) handleAuthorizeRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()

	log.Println("handleAuthorizeRequest|Started")
	var authenticatedUser *models.User

	resp := h.OAuthServer.NewResponse()
	defer resp.Close()

	if ar := h.OAuthServer.HandleAuthorizeRequest(resp, r); ar != nil {
		user, _ := h.handleAuthorizationLoginPage(ctx, ar, w, r)
		if user == nil {
			log.Println("handleAuthorizeRequest|handleAuthorizationLoginPage|FAILED - unauthorized") //TODO: Failure GUI
			return
		}
		log.Println("handleAuthorizeRequest|Authorizing|Started")
		if isRunning, err := h.ApplicationRepo.CheckIfRunningByClientID(ctx, ar.Client.GetId()); !isRunning {
			log.Println("handleAuthorizeRequest|ApplicationRepo.CheckIfRunningByClientID|FAILED - not running state!") //TODO: Failure GUI
			log.Println("handleAuthorizeRequest|ApplicationRepo.CheckIfRunningByClientID|FAILED - err", err)           //TODO: Failure GUI
			return
		}
		app, err := h.ApplicationRepo.GetByClientID(ctx, ar.Client.GetId())
		if err != nil {
			log.Println("handleAuthorizeRequest|ApplicationRepo.GetByClientID|err", err) //TODO: err
			return
		}

		// We want to give the user the ability to keep a record of the
		// application they authorized to use on their behalf. First check if
		// the application exists then update or else if D.N.E. then create.
		authApp, err := h.AuthorizedApplicationRepo.GetByUserIDAndApplicationID(ctx, user.ID, app.ID)
		if err != nil {
			log.Println("handleAuthorizeRequest|AuthorizedApplicationRepo.GetByUserIDAndApplicationID|err", err) //TODO: err
			return
		}
		if authApp != nil {
			authApp.State = models.AuthorizedApplicationPermissionGrantedState
			authApp.ModifiedTime = time.Now()
			err = h.AuthorizedApplicationRepo.UpdateByID(ctx, authApp)
			if err != nil {
				log.Println("handleAuthorizeRequest|AuthorizedApplicationRepo.UpdateByID|err", err) //TODO: err
				return
			}
		} else {
			authApp = &models.AuthorizedApplication{
				TenantID:      user.TenantID,
				UUID:          uuid.NewString(),
				ApplicationID: app.ID,
				UserID:        user.ID,
				State:         models.AuthorizedApplicationPermissionGrantedState,
				CreatedTime:   time.Now(),
				ModifiedTime:  time.Now(),
			}
			err = h.AuthorizedApplicationRepo.Insert(ctx, authApp)
			if err != nil {
				log.Println("handleAuthorizeRequest|AuthorizedApplicationRepo.Insert|err", err) //TODO: err
				return
			}
		}

		// The `osin` library requires we set this value to true so
		// the whole system will generate our access token.
		ar.Authorized = true

		// Store this for our output later...
		authenticatedUser = user

		// This is important because we want to store the user account
		// in the store session for the particular access token when.
		// This means that when we retrieve from the oauth storage for
		// the access token, the user profile will be returned as well.
		// That also means we can utilize our oauth storage in other
		// areas of our app like 'middleware' to protect resources.
		ar.UserData = &models.UserLite{
			TenantID: user.TenantID,
			RoleID:   user.RoleID,
			ID:       user.ID,
			UUID:     user.UUID,
			Timezone: user.Timezone,
			Name:     user.Name,
			State:    user.State,
			Email:    user.Email,
		}

		h.OAuthServer.FinishAuthorizeRequest(resp, r, ar)
		log.Println("handleAuthorizeRequest|Authorizing|Finished")
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		// The following section is the additional `extra` data we want to
		// return in the token object when we return our token.
		if resp.Output != nil && authenticatedUser != nil {
			resp.Output["user_id"] = authenticatedUser.ID
			resp.Output["user_uuid"] = authenticatedUser.UUID
			resp.Output["first_name"] = authenticatedUser.FirstName
			resp.Output["last_name"] = authenticatedUser.LastName
			resp.Output["email"] = authenticatedUser.Email
			resp.Output["language"] = authenticatedUser.Language
			resp.Output["role_id"] = authenticatedUser.RoleID
			resp.Output["tenant_id"] = authenticatedUser.TenantID
			log.Println("handleAuthorizeRequest|added extra to token")
		}
	}
	osin.OutputJSON(resp, w, r)

	log.Println("handleAuthorizeRequest|Finished")
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
