package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

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

	w.Write([]byte(fmt.Sprintf("LOGIN %s <br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"/authorize?%s\" method=\"POST\">", r.URL.RawQuery)))

	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))

	w.Write([]byte("</form>"))

	w.Write([]byte("</body></html>"))

	return nil, nil
}

func (h *Controller) handleAuthorizationPermissionRequest(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	defer r.Body.Close()

	r.ParseForm()

	code := r.FormValue("code")

	w.Header().Set("Content-Type", "text/html; charset=utf-8") // Override the JSON header.
	w.Write([]byte("<html><body>"))
	w.Write([]byte("APP AUTH - CODE<br/>"))
	defer w.Write([]byte("</body></html>"))

	if code == "" {
		w.Write([]byte("Nothing to do"))
		return
	}

	jr := make(map[string]interface{})
	//
	// build access code url
	aurl := fmt.Sprintf("/token?grant_type=authorization_code&client_id=1234&client_secret=aabbccdd&state=xyz&redirect_uri=%s&code=%s",
		url.QueryEscape("http://localhost:8000/appauth/code"), url.QueryEscape(code))
	log.Println(aurl)

	// if parse, download and parse json
	if r.FormValue("doparse") == "1" {
		err := utils.DownloadAccessToken(fmt.Sprintf("http://localhost:14000%s", aurl),
			&osin.BasicAuth{"1234", "aabbccdd"}, jr)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.Write([]byte("<br/>"))
		}
	}

	// show json error
	if erd, ok := jr["error"]; ok {
		w.Write([]byte(fmt.Sprintf("ERROR: %s<br/>\n", erd)))
	}

	// show json access token
	if at, ok := jr["access_token"]; ok {
		w.Write([]byte(fmt.Sprintf("ACCESS TOKEN: %s<br/>\n", at)))
	}

	w.Write([]byte(fmt.Sprintf("FULL RESULT: %+v<br/>\n", jr)))

	// output links
	w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Goto Token URL</a><br/>", aurl)))

	cururl := *r.URL
	curq := cururl.Query()
	curq.Add("doparse", "1")
	cururl.RawQuery = curq.Encode()
	w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Download Token</a><br/>", cururl.String())))

	if rt, ok := jr["refresh_token"]; ok {
		rurl := fmt.Sprintf("/appauth/refresh?code=%s", rt)
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Refresh Token</a><br/>", rurl)))
	}

	if at, ok := jr["access_token"]; ok {
		rurl := fmt.Sprintf("/appauth/info?code=%s", at)
		w.Write([]byte(fmt.Sprintf("<a href=\"%s\">Info</a><br/>", rurl)))
	}
}
