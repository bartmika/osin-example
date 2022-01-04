package controllers

import (
	// "encoding/json"
	"encoding/json"
	"log"
	"net/http"

	"strconv"

	null "gopkg.in/guregu/null.v4"

	"github.com/bartmika/osin-example/internal/idos"
	"github.com/bartmika/osin-example/internal/models"
	"github.com/bartmika/osin-example/internal/validators"
)

func (h *Controller) applicationsListEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := uint64(ctx.Value("user_tenant_id").(uint64))
	// userID := uint64(ctx.Value("user_id").(uint64))

	// Extract our parameters from the URL.
	offsetParamString := r.FormValue("offset")
	offsetParam, _ := strconv.ParseUint(offsetParamString, 10, 64)
	limitParamString := r.FormValue("limit")
	limitParam, _ := strconv.ParseUint(limitParamString, 10, 64)
	if limitParam == 0 || limitParam > 500 {
		limitParam = 100
	}
	searchString := r.FormValue("search")
	sortOrderString := r.FormValue("sort_order")
	if sortOrderString == "" {
		sortOrderString = "ASC"
	}
	sortFieldString := r.FormValue("sort_field")
	if sortFieldString == "" {
		sortFieldString = "name"
	}
	stateParamString := r.FormValue("state")
	stateParam, _ := strconv.ParseUint(stateParamString, 10, 64)

	// Start by defining our base listing filter and then append depending on
	// different cases.
	f := models.ApplicationLiteFilter{
		TenantID:  tenantID,
		SortField: sortFieldString,
		SortOrder: sortOrderString,
		Search:    null.NewString(searchString, searchString != ""),
		Offset:    offsetParam,
		Limit:     limitParam,
		States:    []int8{int8(stateParam)},
	}

	// // For debugging purposes only.
	// log.Println("TenantID", f.TenantID)
	// log.Println("Search", f.Search)
	// log.Println("Offset", f.Offset)
	// log.Println("Limit", f.Limit)
	// log.Println("SortOrder", f.SortOrder)
	// log.Println("SortField", f.SortField)

	arrCh := make(chan []*models.ApplicationLite)
	countCh := make(chan uint64)

	go func() {
		arr, err := h.ApplicationLiteRepo.ListByFilter(ctx, &f)
		if err != nil {
			log.Println("WARNING: applicationsListEndpoint|ListByFilter|err:", err.Error())
			arrCh <- nil
			return
		}
		arrCh <- arr[:]
	}()

	go func() {
		count, err := h.ApplicationLiteRepo.CountByFilter(ctx, &f)
		if err != nil {
			log.Println("WARNING: applicationsListEndpoint|CountByFilter|err:", err.Error())
			countCh <- 0
			return
		}
		countCh <- count
	}()

	arr, count := <-arrCh, <-countCh

	res := idos.NewApplicationLiteListResponseIDO(arr, count)

	if err := json.NewEncoder(w).Encode(&res); err != nil { // [2]
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Controller) applicationCreateEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// tenantID := uint64(ctx.Value("user_tenant_id").(uint64))
	// userID := uint64(ctx.Value("user_id").(uint64))
	// userName := ctx.Value("user_name").(string)
	// timezone := ctx.Value("user_timezone").(string)
	// ipAddress := ctx.Value("IPAddress").(string)

	// The following code will `unmarshal` the user request data or return error.
	var requestData *idos.ApplicationCreateRequestIDO
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform our validation and return validation error on any issues detected.
	isValid, errStr := validators.ValidateApplicationCreateFromRequest(requestData)
	if isValid == false {
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// Next we need to create a `User` account in our system before creating the application.
	m, err := idos.ApplicationCreateRequestUnmarshal(ctx, requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.ApplicationRepo.Insert(ctx, m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m, err = h.ApplicationRepo.GetByUUID(ctx, m.UUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call our 'get details' API endpoint and return it.
	idStr := strconv.FormatUint(m.ID, 10)
	h.applicationGetEndpoint(w, r, idStr)
}

func (h *Controller) applicationGetEndpoint(w http.ResponseWriter, r *http.Request, idStr string) {
	defer r.Body.Close()

	//
	// Get the application based on the primary key.
	//

	// Extract the session details from our "Session" middleware.
	ctx := r.Context()
	tenantID := uint64(ctx.Value("user_tenant_id").(uint64))

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := h.ApplicationRepo.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if m == nil {
		http.Error(w, "application d.n.e.", http.StatusNotFound)
		return
	}
	if tenantID != m.TenantID {
		http.Error(w, "application access forbidden", http.StatusForbidden)
		return
	}

	//
	// Serialize the data.
	//

	ido := idos.ApplicationResponseMarshal(m)
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Controller) applicationUpdateEndpoint(w http.ResponseWriter, r *http.Request, idStr string) {
	ctx := r.Context()
	tenantID := uint64(ctx.Value("user_tenant_id").(uint64))
	// userID := uint64(ctx.Value("user_id").(uint64))
	// userName := ctx.Value("user_name").(string)
	// ipAddress := ctx.Value("IPAddress").(string)

	// Parse the int.
	applicationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m, err := h.ApplicationRepo.GetByID(ctx, applicationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if m == nil {
		http.Error(w, "application d.n.e.", http.StatusNotFound)
		return
	}
	if tenantID != m.TenantID {
		http.Error(w, "application update forbidden", http.StatusForbidden)
		return
	}

	// The following code will `unmarshal` the user request data or return error.
	var requestData *idos.ApplicationUpdateRequestIDO
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform our validation and return validation error on any issues detected.
	isValid, errStr := validators.ValidateApplicationUpdateFromRequest(requestData)
	if isValid == false {
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// Update our record.
	m.Name = requestData.Name
	m.Description = requestData.Description
	m.WebsiteURL = requestData.WebsiteURL
	m.Scope = requestData.Scope
	m.RedirectURL = requestData.RedirectURL
	m.ImageURL = requestData.ImageURL

	// Save to the database.
	err = h.ApplicationRepo.UpdateByID(ctx, m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call our 'get details' API endpoint and return it.
	h.applicationGetEndpoint(w, r, idStr)
}

func (h *Controller) applicationDeleteEndpoint(w http.ResponseWriter, r *http.Request, idStr string) {
	ctx := r.Context()
	// tenantID := uint64(ctx.Value("user_tenant_id").(uint64))
	// userID := uint64(ctx.Value("user_id").(uint64))
	// userName := ctx.Value("user_name").(string)
	// ipAddress := ctx.Value("IPAddress").(string)

	// Parse the int.
	applicationID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.ApplicationRepo.DeleteByID(ctx, applicationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return successfull response.
	w.WriteHeader(http.StatusNoContent)
}
