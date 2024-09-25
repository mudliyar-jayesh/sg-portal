package v1

import (
	"encoding/base64"
	"net/http"
	"sg-portal/internal/models"
	"sg-portal/pkg/util"
	"time"

	"gorm.io/gorm"
)

type UserHandler struct {
	UserRepo         *util.Repository[models.User]
	UserPasswordRepo *util.Repository[models.UserPassword]
}

// NewUserHandler initializes the UserHandler with the user and user password repositories.
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		UserRepo:         util.NewRepository[models.User](db),
		UserPasswordRepo: util.NewRepository[models.UserPassword](db),
	}
}

// GetUserByID retrieves a user by their ID.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL query parameters (e.g., ?id=1)
	userID, err := util.ParseUintParam(r, "id")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Fetch the user by ID
	user, err := h.UserRepo.GetByField("id", userID)
	if err != nil {
		util.HandleError(w, http.StatusNotFound, "User not found")
		return
	}

	// Respond with the user data
	util.RespondJSON(w, http.StatusOK, user)
}

// GetAllUsers retrieves all users from the database.
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch all users
	users, err := h.UserRepo.GetAllByCondition("1 = 1")
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching users")
		return
	}

	// Respond with the list of users
	util.RespondJSON(w, http.StatusOK, &users)
}

// UpdateUser updates the details of an existing user.
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL query parameters (e.g., ?id=1)
	userID, err := util.ParseUintParam(r, "id")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Parse the JSON request body into a map for flexible updates
	userUpdates, err := util.ParseJSONBody[map[string]interface{}](w, r)
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid update data")
		return
	}

	// Ensure that some updates are provided
	if len(*userUpdates) == 0 {
		util.HandleError(w, http.StatusBadRequest, "No updates provided")
		return
	}

	// Apply the updates to the user by ID
	if err := h.UserRepo.UpdateOne("id", userID, *userUpdates); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

// DeleteUser deletes a user by their ID.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL query parameters (e.g., ?id=1)
	userID, err := util.ParseUintParam(r, "id")
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Delete the user by ID
	if err := h.UserRepo.Delete("id = ?", userID); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error deleting user")
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}

// GetUserProfile retrieves the profile of the currently authenticated user.
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context (set by the token middleware)
	userID, ok := util.UserIDFromContext(r.Context())
	if !ok {
		util.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Fetch the user by ID
	user, err := h.UserRepo.GetByField("id", userID)
	if err != nil {
		util.HandleError(w, http.StatusNotFound, "User not found")
		return
	}

	// Respond with the user's profile data
	util.RespondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get the user ID from the context (set by the token middleware)
	userID, ok := util.UserIDFromContext(r.Context())
	if !ok {
		util.HandleError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse the request body to get the old password and the new password
	passwordData, err := util.ParseJSONBody[struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"` // Base64 encoded
	}](w, r)
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid request data")
		return
	}

	// Decode the new base64-encoded password
	newPasswordBytes, err := base64.StdEncoding.DecodeString(passwordData.NewPassword)
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid password encoding")
		return
	}
	newPassword := string(newPasswordBytes)

	// Fetch the current user password details
	userPassword, err := h.UserPasswordRepo.GetByField("user_id", userID)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error fetching user password")
		return
	}

	// Validate the old password
	if err := models.ValidatePassword(passwordData.OldPassword, userPassword.Salt, userPassword.Password); err != nil {
		util.HandleError(w, http.StatusUnauthorized, "Incorrect old password")
		return
	}

	// Generate a new salt and hash the new password
	newSalt, err := models.GenerateSalt()
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error generating salt")
		return
	}
	newHashedPassword, err := models.HashPassword(newPassword, newSalt)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// Update the user's password and salt in the database
	if err := h.UserPasswordRepo.UpdateOne("user_id", userID, map[string]interface{}{
		"password": newHashedPassword,
		"salt":     newSalt,
		"updated_at": time.Now(),
	}); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error updating password")
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
}
