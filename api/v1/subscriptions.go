package v1

import (
	"net/http"

	"gorm.io/gorm"
	"sg-portal/internal/models"
	"sg-portal/pkg/util"
)

type SubscriptionHandler struct {
	SubscriptionRepo        *util.Repository[models.Subscription]
	UserSubscriptionRepo    *util.Repository[models.UserSubscriptionMapping]
	FeatureSubscriptionRepo *util.Repository[models.FeatureSubscriptionMapping]
}

// NewSubscriptionHandler initializes the SubscriptionHandler with the repositories
func NewSubscriptionHandler(db *gorm.DB) *SubscriptionHandler {
	return &SubscriptionHandler{
		SubscriptionRepo:        util.NewRepository[models.Subscription](db),
		UserSubscriptionRepo:    util.NewRepository[models.UserSubscriptionMapping](db),
		FeatureSubscriptionRepo: util.NewRepository[models.FeatureSubscriptionMapping](db),
	}
}

// CreateSubscription creates a new subscription
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	subscription, err := util.ParseJSONBody[models.Subscription](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.SubscriptionRepo.Create(subscription); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error creating subscription")
		return
	}
	util.RespondJSON(w, http.StatusCreated, subscription)
}

// GetAllSubscriptions returns all subscriptions
func (h *SubscriptionHandler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions, err := h.SubscriptionRepo.GetAllByCondition("1 = 1")
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching subscriptions")
		return
	}
	util.RespondJSON(w, http.StatusOK, &subscriptions)
}

// UpdateSubscription updates an existing subscription
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	// Extract subscriptionId from query parameters
	subscriptionId, err := util.ParseUintParam(r, "subscriptionId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse the JSON request body into a map for flexible updates
	subscriptionUpdates, err := util.ParseJSONBody[map[string]interface{}](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}

	// Validate that the update map is not empty
	if len(*subscriptionUpdates) == 0 {
		util.HandleError(w, http.StatusBadRequest, "No updates provided")
		return
	}

	// Apply the updates to the subscription by ID
	if err := h.SubscriptionRepo.UpdateOne("id", subscriptionId, *subscriptionUpdates); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error updating subscription")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}

// DeleteSubscription deletes a subscription by its ID
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	// Extract subscriptionId from query parameters
	subscriptionId, err := util.ParseUintParam(r, "subscriptionId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use the repository's Delete method to delete the subscription by ID
	if err := h.SubscriptionRepo.Delete("id = ?", subscriptionId); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error deleting subscription")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}

// MapUserToSubscription maps a single user to a subscription
func (h *SubscriptionHandler) MapUserToSubscription(w http.ResponseWriter, r *http.Request) {
	mapping, err := util.ParseJSONBody[models.UserSubscriptionMapping](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.UserSubscriptionRepo.Create(mapping); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error mapping user to subscription")
		return
	}
	util.RespondJSON(w, http.StatusCreated, mapping)
}

// GetSubscriptionsByUser returns all subscriptions mapped to a specific user
func (h *SubscriptionHandler) GetSubscriptionsByUser(w http.ResponseWriter, r *http.Request) {
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	mappings, err := h.UserSubscriptionRepo.GetAllByCondition("user_id = ?", userId)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching user-subscription mappings")
		return
	}
	var subscriptionIds []uint32
	for _, mapping := range mappings {
		subscriptionIds = append(subscriptionIds, mapping.SubscriptionId)
	}
	subscriptions, err := h.SubscriptionRepo.GetAllByCondition("id IN ?", subscriptionIds)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching subscriptions")
		return
	}
	util.RespondJSON(w, http.StatusOK, &subscriptions)
}

// DeleteUserSubscriptionMapping deletes a user-subscription mapping (hard delete)
func (h *SubscriptionHandler) DeleteUserSubscriptionMapping(w http.ResponseWriter, r *http.Request) {
	// Extract userId and subscriptionId from query parameters
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	subscriptionId, err := util.ParseUintParam(r, "subscriptionId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use the repository's Delete method to delete the UserSubscriptionMapping where user_id and subscription_id match
	condition := "user_id = ? AND subscription_id = ?"
	if err := h.UserSubscriptionRepo.Delete(condition, userId, subscriptionId); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error deleting user-subscription mapping")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}

// MapFeatureToSubscription maps a feature to a subscription
func (h *SubscriptionHandler) MapFeatureToSubscription(w http.ResponseWriter, r *http.Request) {
	mapping, err := util.ParseJSONBody[models.FeatureSubscriptionMapping](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}
	if err := h.FeatureSubscriptionRepo.Create(mapping); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error mapping feature to subscription")
		return
	}
	util.RespondJSON(w, http.StatusCreated, mapping)
}

// GetFeaturesBySubscription returns all features mapped to a specific subscription
func (h *SubscriptionHandler) GetFeaturesBySubscription(w http.ResponseWriter, r *http.Request) {
	subscriptionId, err := util.ParseUintParam(r, "subscriptionId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	mappings, err := h.FeatureSubscriptionRepo.GetAllByCondition("subscription_id = ?", subscriptionId)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching feature-subscription mappings")
		return
	}
	var featureIds []uint32
	for _, mapping := range mappings {
		featureIds = append(featureIds, mapping.FeatureId)
	}
	features, err := h.FeatureSubscriptionRepo.GetAllByCondition("id IN ?", featureIds)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching features")
		return
	}
	util.RespondJSON(w, http.StatusOK, &features)
}

// DeleteFeatureSubscriptionMapping deletes a feature-subscription mapping (hard delete)
func (h *SubscriptionHandler) DeleteFeatureSubscriptionMapping(w http.ResponseWriter, r *http.Request) {
	// Extract featureId and subscriptionId from query parameters
	featureId, err := util.ParseUintParam(r, "featureId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	subscriptionId, err := util.ParseUintParam(r, "subscriptionId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use the repository's Delete method to delete the FeatureSubscriptionMapping where feature_id and subscription_id match
	condition := "feature_id = ? AND subscription_id = ?"
	if err := h.FeatureSubscriptionRepo.Delete(condition, featureId, subscriptionId); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error deleting feature-subscription mapping")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}
