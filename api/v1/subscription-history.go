package v1

import (
	"net/http"
	"time"

	"gorm.io/gorm"
	"sg-portal/internal/models"
	"sg-portal/pkg/util"
)

type UserSubscriptionHistoryHandler struct {
	UserSubscriptionHistoryRepo *util.Repository[models.UserSubscriptionHistory]
}

// NewUserSubscriptionHistoryHandler initializes the handler with the repository
func NewUserSubscriptionHistoryHandler(db *gorm.DB) *UserSubscriptionHistoryHandler {
	return &UserSubscriptionHistoryHandler{
		UserSubscriptionHistoryRepo: util.NewRepository[models.UserSubscriptionHistory](db),
	}
}

// CreateUserSubscriptionHistory creates a new user subscription history entry
func (h *UserSubscriptionHistoryHandler) CreateUserSubscriptionHistory(w http.ResponseWriter, r *http.Request) {
	history, err := util.ParseJSONBody[models.UserSubscriptionHistory](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}

	history.StartDate = time.Now() // Set the start date to the current date

	if err := h.UserSubscriptionHistoryRepo.Create(history); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error creating user subscription history")
		return
	}
	util.RespondJSON(w, http.StatusCreated, history)
}

// GetAllUserSubscriptionHistories returns all subscription histories for users
func (h *UserSubscriptionHistoryHandler) GetAllUserSubscriptionHistories(w http.ResponseWriter, r *http.Request) {
	histories, err := h.UserSubscriptionHistoryRepo.GetAllByCondition("1 = 1")
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching subscription histories")
		return
	}
	util.RespondJSON(w, http.StatusOK, &histories)
}

// GetUserSubscriptionHistory returns a specific user's subscription history
func (h *UserSubscriptionHistoryHandler) GetUserSubscriptionHistory(w http.ResponseWriter, r *http.Request) {
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}
	history, err := h.UserSubscriptionHistoryRepo.GetByField("user_id", userId)
	if err != nil {
		util.HandleError(w, http.StatusNotFound, "User subscription history not found")
		return
	}
	util.RespondJSON(w, http.StatusOK, history)
}

// UpdateUserSubscriptionHistory updates an existing subscription history entry
func (h *UserSubscriptionHistoryHandler) UpdateUserSubscriptionHistory(w http.ResponseWriter, r *http.Request) {
	// Extract userId from query parameters
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse the JSON request body into a map for flexible updates
	historyUpdates, err := util.ParseJSONBody[map[string]interface{}](w, r)
	if err != nil {
		return // Error already handled by ParseJSONBody
	}

	// Apply the updates to the user subscription history by userId
	if err := h.UserSubscriptionHistoryRepo.UpdateOne("user_id", userId, *historyUpdates); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error updating user subscription history")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}

// DeleteUserSubscriptionHistory deletes a subscription history by userId
func (h *UserSubscriptionHistoryHandler) DeleteUserSubscriptionHistory(w http.ResponseWriter, r *http.Request) {
	userId, err := util.ParseUintParam(r, "userId")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Use the repository's Delete method to delete the subscription history by userId
	if err := h.UserSubscriptionHistoryRepo.Delete("user_id = ?", userId); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error deleting user subscription history")
		return
	}

	// Respond with success status
	w.WriteHeader(http.StatusOK)
}
