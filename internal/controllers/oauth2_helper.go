package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/openshift/osin"

	"github.com/bartmika/osin-example/internal/models"
	"github.com/bartmika/osin-example/internal/utils"
)

func (h *Controller) authenticatedUser(ctx context.Context, email string, password string) (*models.User, error) {
	// Lookup the user in our database, else return a `400 Bad Request` error.
	user, err := h.UserRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Incorrect email or password")
	}

	// Verify the inputted password and hashed password match.
	passwordMatch := utils.CheckPasswordHash(password, user.PasswordHash)
	if passwordMatch == false {
		return nil, errors.New("Incorrect email or password")
	}
	return user, nil
}

// handleAuthorizationLoginPage will be a web-page when the user makes a GET
// request and handles user login verification when making POST requests.
func (h *Controller) handleAuthorizationLoginPage(ctx context.Context, ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) (*models.User, error) {
	r.ParseForm()
	if r.Method == "POST" {
		return h.authenticatedUser(ctx, r.FormValue("login"), r.FormValue("password"))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Override the JSON header.
	w.Write([]byte("<html><body>"))

	w.Write([]byte("<h1>Third-Party Application Authorization</h1>"))
	w.Write([]byte("<p>You are about to authorize this application. This application will be granted <b>full access</b> to your account. If you would like to proceed, please enter your credentials now or close this window if you want to cancel.</p>"))

	w.Write([]byte(fmt.Sprintf("<form action=\"/authorize?%s\" method=\"POST\">", r.URL.RawQuery)))

	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

	return nil, nil
}
