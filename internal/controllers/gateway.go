package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/bartmika/osin-example/internal/idos"
	"github.com/bartmika/osin-example/internal/models"
	"github.com/bartmika/osin-example/internal/utils"
	"github.com/bartmika/osin-example/internal/validators"
)

// To run this API, try running in your console:
// $ http post 127.0.0.1:8000/api/v1/register email="fherbert@dune.com" password="the-spice-must-flow" name="Frank Herbert"
func (h *Controller) registerEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initialize our array which will store all the results from the remote server.
	var requestData idos.RegisterRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform our validation and return validation error on any issues detected.
	isValid, errStr := validators.ValidateRegisterRequest(&requestData)
	if isValid == false {
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// // For debugging purposes, print our output so you can see the code working.
	// fmt.Println(requestData.Name)
	// fmt.Println(requestData.Email)
	// fmt.Println(requestData.Password)

	// Lookup the email and if it is not unique we need to generate a `400 Bad Request` response.
	if userFound, _ := h.UserRepo.CheckIfExistsByEmail(ctx, requestData.Email); userFound {
		http.Error(w, "Email alread exists", http.StatusBadRequest)
		return
	}

	//
	// Create our tenant organization.
	//

	t := &models.Tenant{
		UUID:         uuid.NewString(),
		Name:         uuid.NewString(),
		State:        models.TenantActiveState,
		Timezone:     "UTC",
		Language:     requestData.Language,
		CreatedTime:  time.Now().UTC(),
		ModifiedTime: time.Now().UTC(),
	}

	// Save our new user account.
	if err := h.TenantRepo.Insert(ctx, t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t, err = h.TenantRepo.GetByName(ctx, t.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//
	// Create our user account.
	//

	// Secure our password.
	passwordHash, err := utils.HashPassword(requestData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := &models.User{
		UUID:         uuid.NewString(),
		TenantID:     t.ID, // TenantID
		Email:        requestData.Email,
		FirstName:    requestData.FirstName,
		LastName:     requestData.LastName,
		Name:         requestData.FirstName + " " + requestData.LastName,
		LexicalName:  requestData.LastName + ", " + requestData.FirstName,
		PasswordHash: passwordHash,
		State:        models.UserActiveState,
		RoleID:       models.UserTenantAdminRoleID,
		Timezone:     "UTC",
		Language:     requestData.Language,
		CreatedTime:  time.Now().UTC(),
		ModifiedTime: time.Now().UTC(),
	}

	// Save our new user account.
	if err := h.UserRepo.Insert(ctx, m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate our response.
	responseData := idos.RegisterResponseIDO{
		Message: "You have successfully registered an account.",
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil { // [2]
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// To run this API, try running in your console:
// $ http post 127.0.0.1:8000/api/v1/login email="fherbert@dune.com" password="the-spice-must-flow"
func (h *Controller) loginEndpoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	ctx := r.Context()

	var requestData idos.LoginRequestIDO

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform our validation and return validation error on any issues detected.
	isValid, errStr := validators.ValidateLoginRequest(&requestData)
	if isValid == false {
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// DEVELOPERS NOTE:
	// WE WILL TAKE ADVANTAGE OF OUR OAUTH SERVER AND MAKE THE CALL TO THE
	// PASSWORD GRANT TO GET OUR TOKENS.
	log.Println("Beginning Password Based Authorization")

	client := &oauth2.Config{
		ClientID:     "1234",
		ClientSecret: "aabbccdd",
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8000/authorize",
			TokenURL: "http://localhost:8000/token",
		},
		RedirectURL: "http://localhost:8000/appauth/code",
	}

	// NOTE: https://pkg.go.dev/golang.org/x/oauth2#Config.PasswordCredentialsToken
	token, err := client.PasswordCredentialsToken(ctx, requestData.Email, requestData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Finshed Password Based Authorization")

	// Finally return success.
	responseData := idos.LoginResponseIDO{
		FirstName:    token.Extra("first_name").(string),
		LastName:     token.Extra("last_name").(string),
		Email:        token.Extra("email").(string),
		RoleID:       int8(token.Extra("role_id").(float64)),
		TenantID:     uint64(token.Extra("tenant_id").(float64)),
		Language:     token.Extra("language").(string),
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// To run this API, try running in your console:
// $ http post 127.0.0.1:8000/api/v1/refresh-token value="xxx"
func (h *Controller) postRefreshToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var requestData idos.RefreshTokenRequestIDO

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// For debugging purposes, print our output so you can see the code working.
	log.Println(requestData.Value)

	// ctx := r.Context()

	// // Verify our refresh token.
	// sessionUUID, err := utils.ProcessJWTToken(h.SecretSigningKeyBin, requestData.Value)
	// if err != nil {
	// 	http.Error(w, "Unauthorized - refresh token expired or invalid", http.StatusUnauthorized)
	// 	return
	// }
	//
	// // Lookup our user profile in the session or return 500 error.
	// user, err := h.SessionManager.GetUser(ctx, sessionUUID)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Generate our JWT token.
	// newSessionUUID := uuid.NewString()
	// newSessionExpiryTime := time.Hour * 24 * 7 // 1 week
	// accessToken, refreshToken, err := utils.GenerateJWTTokenPair(h.SecretSigningKeyBin, newSessionUUID, newSessionExpiryTime)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Save our new session.
	// err = h.SessionManager.SaveUser(ctx, newSessionUUID, user, newSessionExpiryTime)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Finally return success.
	// responseData := idos.RefreshTokenResponse{
	// 	AccessToken:  accessToken,
	// 	RefreshToken: refreshToken,
	// }
	// if err := json.NewEncoder(w).Encode(&responseData); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

func (h *Controller) profileEndpoint(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	// userID := uint64(ctx.Value("user_id").(uint64))
	// user, err := h.UserRepo.GetByID(ctx, userID)
	//
	// // Start our session.
	// sessionExpiryTime := time.Hour * 24 * 7 // 1 week
	// sessionUUID := uuid.NewString()
	// err = h.SessionManager.SaveUser(ctx, sessionUUID, user, sessionExpiryTime)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Generate our JWT token.
	// accessToken, refreshToken, err := utils.GenerateJWTTokenPair(h.SecretSigningKeyBin, sessionUUID, sessionExpiryTime)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// // Update our results.
	// user.AccessToken = accessToken
	// user.RefreshToken = refreshToken
	// user.PasswordHash = ""
	//
	// // Return our serialized result.
	// if err := json.NewEncoder(w).Encode(&user); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

func (h *Controller) getVersion(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OSIN Example v1.0"))
}

// API endpoint used to test out making posts to our API gateway. Here is an example:
//
// http post 127.0.0.1:8000/api/v1/greeting name=Bart
//
func (h *Controller) postGreeting(w http.ResponseWriter, r *http.Request) {
	var requestData idos.GreetingRequestIDO
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := idos.GreetingResponseIDO{
		Message: "Hello," + requestData.Name,
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil { // [2]
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
