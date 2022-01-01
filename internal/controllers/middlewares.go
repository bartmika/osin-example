package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

// Middleware will split the full URL path into slash-sperated parts and save to
// the context to flow downstream in the app for this particular request.
func (h *Controller) URLProcessorMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Split path into slash-separated parts, for example, path "/foo/bar"
		// gives p==["foo", "bar"] and path "/" gives p==[""]. Our API starts with
		// "/api", as a result we will start the array slice at "2".
		p := strings.Split(r.URL.Path, "/")[1:]

		// log.Println(r.URL.Path, p) // For debugging purposes only.

		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()
		ctx = context.WithValue(ctx, "url_split", p)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

func (h *Controller) IPAddressMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the IPAddress. Code taken from: https://stackoverflow.com/a/55738279
		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		// Save our IP address to the context.
		ctx := r.Context()
		ctx = context.WithValue(ctx, "IPAddress", IPAddress)
		fn(w, r.WithContext(ctx)) // Flow to the next middleware.
	}
}

// The purpose of this middleware is to return a `401 unauthorized` error if
// the user is not authorized and visiting a protected URL.
func (h *Controller) ProtectedURLsMiddleware(fn http.HandlerFunc) http.HandlerFunc {
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
		}

		// log.Println(urlSplit)

		// DEVELOPERS NOTE:
		// If the URL cannot be split into the size we want then skip running
		// this middleware.
		if len(urlSplit) <= 2 {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			return
		}

		if skipPath[urlSplit[2]] {
			fn(w, r.WithContext(ctx)) // Flow to the next middleware.
		} else {
			// Get our authorization information.
			isAuthorized, ok := ctx.Value("is_authorized").(bool)

			// Either accept continuing execution or return 401 error.
			if ok && isAuthorized {
				fn(w, r.WithContext(ctx)) // Flow to the next middleware.
			} else {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
	}
}

func (h *Controller) PaginationMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()

		// Setup our variables for the paginator.
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pageTokenString := r.FormValue("page_token")
		pageSizeString := r.FormValue("page_size")

		// Convert to unsigned 64-bit integer.
		pageToken, err := strconv.ParseUint(pageTokenString, 10, 64)
		if err != nil {
			// DEVELOPERS NOTE: ALWAYS DEFINE 100 IF NOT SPECIFIED OR ERROR.
			pageToken = 0
		}
		pageSize, err := strconv.ParseUint(pageSizeString, 10, 64)
		if err != nil {
			// DEVELOPERS NOTE: ALWAYS DEFINE 100 IF NOT SPECIFIED OR ERROR.
			pageSize = 100
		}

		// Attach the 'page' parameter value to our context to be used.
		ctx = context.WithValue(ctx, "pageTokenParm", pageToken)
		ctx = context.WithValue(ctx, "pageSizeParam", pageSize)

		// Flow to the next middleware.
		fn(w, r.WithContext(ctx))
	}
}

func (h *Controller) AttachMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	// Attach our middleware handlers here. Please note that all our middleware
	// will start from the bottom and proceed upwards.
	// Ex: `URLProcessorMiddleware` will be executed first and
	//     `PaginationMiddleware` will be executed last.
	fn = h.ProtectedURLsMiddleware(fn)
	fn = h.IPAddressMiddleware(fn)
	fn = h.OSINProcessorMiddleware(fn)
	fn = h.PaginationMiddleware(fn)
	fn = h.URLProcessorMiddleware(fn)

	return func(w http.ResponseWriter, r *http.Request) {
		// Flow to the next middleware.
		fn(w, r)
	}
}
