package v1

import (
	"net/http"

	"gorm.io/gorm"
	"sg-portal/internal/models"
	"sg-portal/pkg/util"
)

type FeatureHandler struct {
	FeatureRepo     *util.Repository[models.Feature]
	UserFeatureRepo *util.Repository[models.UserFeatureMapping]
}

// NewFeatureHandler initializes the FeatureHandler with the repositories
func NewFeatureHandler(db *gorm.DB) *FeatureHandler {
	return &FeatureHandler{
		FeatureRepo:     util.NewRepository[models.Feature](db),
		UserFeatureRepo: util.NewRepository[models.UserFeatureMapping](db),
	}
}

// CreateFeature creates a new feature
func (h *FeatureHandler) CreateFeature(w http.ResponseWriter, r *http.Request) {
	feature, err := util.ParseJSONBody[models.Feature](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.FeatureRepo.Create(feature); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error creating feature")
		return
	}
	util.RespondJSON(w, http.StatusCreated, feature)
}

// MapUserToFeature maps a single user to a feature
func (h *FeatureHandler) MapUserToFeature(w http.ResponseWriter, r *http.Request) {
	mapping, err := util.ParseJSONBody[models.UserFeatureMapping](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.UserFeatureRepo.Create(mapping); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error mapping user to feature")
		return
	}
	util.RespondJSON(w, http.StatusCreated, mapping)
}

// MapUsersToFeature maps multiple users to a feature
func (h *FeatureHandler) MapUsersToFeature(w http.ResponseWriter, r *http.Request) {
	mappings, err := util.ParseJSONBody[[]models.UserFeatureMapping](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.UserFeatureRepo.CreateMultiple(mappings); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error mapping users to feature")
		return
	}
	util.RespondJSON(w, http.StatusCreated, mappings)
}

// GetAllFeatures returns all features
func (h *FeatureHandler) GetAllFeatures(w http.ResponseWriter, r *http.Request) {
	features, err := h.FeatureRepo.GetAllByCondition("1 = 1")
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching features")
		return
	}
	util.RespondJSON(w, http.StatusOK, &features)
}

// GetFeaturesByUser returns all features mapped to a specific user
func (h *FeatureHandler) GetFeaturesByUser(w http.ResponseWriter, r *http.Request) {
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	mappings, err := h.UserFeatureRepo.GetAllByCondition("user_id = ?", userId)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching user-feature mappings")
		return
	}
	var featureIds []uint32
	for _, mapping := range mappings {
		featureIds = append(featureIds, mapping.FeatureId)
	}
	features, err := h.FeatureRepo.GetAllByCondition("id IN ?", featureIds)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching features")
		return
	}
	util.RespondJSON(w, http.StatusOK, &features)
}

// UpdateFeature updates an existing feature
func (h *FeatureHandler) UpdateFeature(w http.ResponseWriter, r *http.Request) {
	// Extract featureId from query parameters
	featureId, err := util.ParseUintParam(r, "featureId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse the JSON request body into a map for flexible updates
	featureUpdates, err := util.ParseJSONBody[map[string]interface{}](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}

	// Validate that the update map is not empty
	if len(*featureUpdates) == 0 {
		util.HandleError(w, http.StatusBadRequest, "No updates provided")
		return
	}

	// Apply the updates to the feature by ID
	if err := h.FeatureRepo.UpdateOne("id", featureId, *featureUpdates); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error updating feature")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}

// DeleteFeature deletes a feature by its ID
func (h *FeatureHandler) DeleteFeature(w http.ResponseWriter, r *http.Request) {
	// Extract featureId from query parameters
	featureId, err := util.ParseUintParam(r, "featureId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use the repository's Delete method to delete the feature by ID
	if err := h.FeatureRepo.Delete("id = ?", featureId); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error deleting feature")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}

// DeleteUserFeatureMapping deletes a user-feature mapping (hard delete)
func (h *FeatureHandler) DeleteUserFeatureMapping(w http.ResponseWriter, r *http.Request) {
	// Extract userId and featureId from query parameters
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	featureId, err := util.ParseUintParam(r, "featureId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use the repository's Delete method to delete the UserFeatureMapping where user_id and feature_id match
	condition := "user_id = ? AND feature_id = ?"
	if err := h.UserFeatureRepo.Delete(condition, userId, featureId); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error deleting user-feature mapping")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}
