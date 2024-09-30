package v1

import (
	"net/http"
	"sg-portal/internal/models"
	"sg-portal/pkg/util"

	"gorm.io/gorm"
)

type CompanyHandler struct {
	UserRepo          *util.Repository[models.User]
	TenantRepo        *util.Repository[models.Tenant]
	TenantMappingRepo *util.Repository[models.UserTenantMapping]
}

func NewCompanyHandler(db *gorm.DB) *CompanyHandler {
	return &CompanyHandler{
		UserRepo:          util.NewRepository[models.User](db),
		TenantRepo:        util.NewRepository[models.Tenant](db),
		TenantMappingRepo: util.NewRepository[models.UserTenantMapping](db),
	}
}

func (h *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {

	userId, paramErr := util.ParseUintParam(r, "id")
	if paramErr != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid User Id Provided")
		return
	}

	mappings, mappingErr := h.TenantMappingRepo.GetAllByCondition("user_id = ?", userId)
	if mappingErr != nil || len(mappings) < 1 {
		util.HandleError(w, http.StatusNoContent, "No companies mapped to the user")
		return
	}

	var tenantIds []uint64
	for _, mapping := range mappings {
		tenantIds = append(tenantIds, mapping.TenantId)
	}

	tenants, err := h.TenantRepo.GetAllByCondition("id IN ?", tenantIds)
	if err != nil || len(tenants) < 1 {
		util.HandleError(w, http.StatusNoContent, "No companies found")
		return
	}
	util.RespondJSON(w, http.StatusOK, &tenants)
}
