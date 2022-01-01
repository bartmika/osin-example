package controllers

import (
	// "encoding/json"
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/bartmika/osin-example/internal/idos"
	"github.com/bartmika/osin-example/internal/models"
	"github.com/bartmika/osin-example/internal/validators"
)

func (h *Controller) userGetEndpoint(w http.ResponseWriter, r *http.Request, idStr string) {
	defer r.Body.Close()

	//
	// Get the user based on the primary key.
	//

	// Extract the session details from our "Session" middleware.
	ctx := r.Context()
	userID := uint64(ctx.Value("user_tenant_id").(uint64))
	role_id := uint64(ctx.Value("user_role_id").(int8))

	// Permission handling - If use is not administrator then error.
	if role_id == models.UserTenantStudentRoleID {
		http.Error(w, "Forbidden - You are not an administrator", http.StatusForbidden)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m, err := h.UserRepo.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if m == nil {
		http.Error(w, "user d.n.e.", http.StatusNotFound)
		return
	}
	if userID != m.ID {
		http.Error(w, "user access forbidden", http.StatusForbidden)
		return
	}

	//
	// Serialize the data.
	//

	ido := idos.NewUserIDO(m)
	if err := json.NewEncoder(w).Encode(&ido); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Controller) userUpdateEndpoint(w http.ResponseWriter, r *http.Request, idStr string) {
	ctx := r.Context()
	ctxUserID := uint64(ctx.Value("user_tenant_id").(uint64))
	// userID := uint64(ctx.Value("user_id").(uint64))
	// userName := ctx.Value("user_name").(string)
	// ipAddress := ctx.Value("IPAddress").(string)

	// Parse the int.
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	m, err := h.UserRepo.GetByID(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if m == nil {
		http.Error(w, "user d.n.e.", http.StatusNotFound)
		return
	}
	if ctxUserID != m.ID {
		http.Error(w, "user update forbidden", http.StatusForbidden)
		return
	}

	// The following code will `unmarshal` the user request data or return error.
	var requestData *idos.UserUpdateRequestIDO
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform our validation and return validation error on any issues detected.
	isValid, errStr := validators.ValidateUserUpdateFromRequest(requestData)
	if isValid == false {
		http.Error(w, errStr, http.StatusBadRequest)
		return
	}

	// Update our record.
	m.Name = requestData.Name

	// Save to the database.
	err = h.UserRepo.UpdateByID(ctx, m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Call our 'get details' API endpoint and return it.
	h.userGetEndpoint(w, r, idStr)
}
