package v1

import (
	"net/http"

	"sg-portal/internal/models"
	"sg-portal/pkg/util"

	"gorm.io/gorm"
)

type TenantHandler struct {
	TenantRepo     *util.Repository[models.Tenant]
	UserTenantRepo *util.Repository[models.UserTenantMapping]
}

// NewTenantHandler initializes the TenantHandler with the repositories
func NewTenantHandler(db *gorm.DB) *TenantHandler {
	return &TenantHandler{
		TenantRepo:     util.NewRepository[models.Tenant](db),
		UserTenantRepo: util.NewRepository[models.UserTenantMapping](db),
	}
}

// CreateTenant creates a new tenant
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	tenant, err := util.ParseJSONBody[models.Tenant](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.TenantRepo.Create(tenant); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error creating tenant")
		return
	}
	util.RespondJSON(w, http.StatusCreated, tenant)
}

// MapUserToTenant maps a single user to a tenant
func (h *TenantHandler) MapUserToTenant(w http.ResponseWriter, r *http.Request) {
	mapping, err := util.ParseJSONBody[models.UserTenantMapping](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.UserTenantRepo.Create(mapping); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error mapping user to tenant")
		return
	}
	util.RespondJSON(w, http.StatusCreated, mapping)
}

// MapUsersToTenant maps multiple users to a tenant
func (h *TenantHandler) MapUsersToTenant(w http.ResponseWriter, r *http.Request) {
	mappings, err := util.ParseJSONBody[[]models.UserTenantMapping](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.UserTenantRepo.CreateMultiple(mappings); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error mapping users to tenant")
		return
	}
	util.RespondJSON(w, http.StatusCreated, mappings)
}

// GetAllTenants returns all tenants
func (h *TenantHandler) GetAllTenants(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.TenantRepo.GetAllByCondition("1 = 1")
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching tenants")
		return
	}
	util.RespondJSON(w, http.StatusOK, &tenants)
}

// GetTenantsByUser returns all tenants mapped to a specific user
func (h *TenantHandler) GetTenantsByUser(w http.ResponseWriter, r *http.Request) {
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	mappings, err := h.UserTenantRepo.GetAllByCondition("user_id = ?", userId)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching user-tenant mappings")
		return
	}
	var tenantIds []uint64
	for _, mapping := range mappings {
		tenantIds = append(tenantIds, mapping.TenantId)
	}
	tenants, err := h.TenantRepo.GetAllByCondition("id IN ?", tenantIds)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching tenants")
		return
	}
	util.RespondJSON(w, http.StatusOK, &tenants)
}

// UpdateTenant updates an existing tenant
func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	// Extract tenantId from query parameters
	tenantId, err := util.ParseUintParam(r, "tenantId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse the JSON request body into a map for flexible updates
	tenantUpdates, err := util.ParseJSONBody[map[string]interface{}](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}

	// Validate that the update map is not empty
	if len(*tenantUpdates) == 0 {
		util.HandleError(w, http.StatusBadRequest, "No updates provided")
		return
	}

	// Apply the updates to the tenant by ID
	if err := h.TenantRepo.UpdateOne("id", tenantId, *tenantUpdates); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error updating tenant")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}

// DeleteUserTenantMapping deletes a user-tenant mapping (hard delete)
func (h *TenantHandler) DeleteUserTenantMapping(w http.ResponseWriter, r *http.Request) {
	// Extract userId and tenantId from query parameters
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	tenantId, err := util.ParseUintParam(r, "tenantId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use the repository's Delete method to delete the UserTenantMapping where user_id and tenant_id match
	condition := "user_id = ? AND tenant_id = ?"
	if err := h.UserTenantRepo.Delete(condition, userId, tenantId); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error deleting user-tenant mapping")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}
