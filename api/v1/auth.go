package v1

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"sg-portal/internal/models"
	"sg-portal/pkg/util"
	"time"

	"gorm.io/gorm"
)

type AuthHandler struct {
	UserRepo         *util.Repository[models.User]
	UserPasswordRepo *util.Repository[models.UserPassword]
	TokenRepo        *util.Repository[models.Token]
}

// NewAuthHandler initializes the auth handler with the repositories.
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		UserRepo:         util.NewRepository[models.User](db),
		UserPasswordRepo: util.NewRepository[models.UserPassword](db),
		TokenRepo:        util.NewRepository[models.Token](db),
	}
}

// Register handles user registration.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into the User and base64-encoded password.
	userData := struct {
		Email        string `json:"email"`
		Name         string `json:"name"`
		MobileNumber string `json:"mobile_number"`
		Password     string `json:"password"` // Base64 encoded
		Type         string `json:"type"`     // New field for user type
	}{}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid registration data")
		return
	}

	// Validate the user type
	if userData.Type != models.UserTypeClient && userData.Type != models.UserTypeSystem {
		util.HandleError(w, http.StatusBadRequest, "Invalid user type")
		return
	}

	// Decode the base64-encoded password
	passwordBytes, err := base64.StdEncoding.DecodeString(userData.Password)
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid password encoding")
		return
	}
	password := string(passwordBytes)

	// Create user entity
	user := &models.User{
		Email:        userData.Email,
		Name:         userData.Name,
		MobileNumber: userData.MobileNumber,
		Type:         userData.Type, // Set the user type
	}

	// Create user record
	if err := h.UserRepo.Create(user); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	// Generate salt and hash the password
	salt, err := models.GenerateSalt()

	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error generating salt")
		return
	}

	hashedPassword, err := models.HashPassword(password, salt)
	if err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// Store the password hash and salt
	userPassword := models.NewUserPassword(user.ID, hashedPassword, salt)
	if err := h.UserPasswordRepo.Create(userPassword); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error storing password")
		return
	}

	// add default demo subscription
	subscriptionRepo := util.NewRepository[models.Subscription](util.Db)
	demoSubscription, subErr := subscriptionRepo.GetByField("Code", "demo")
	if subErr != nil || demoSubscription == nil {
		util.HandleError(w, http.StatusInternalServerError, "Could not get Demo Subscription ")
		return
	}

	userSubscriptionMapping := &models.UserSubscriptionMapping{
		SubscriptionId: demoSubscription.ID,
		UserId:         user.ID,
	}

	userSubscriptionRepo := util.NewRepository[models.UserSubscriptionMapping](util.Db)
	if err := userSubscriptionRepo.Create(userSubscriptionMapping); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Could not get Demo Subscription ")
		return
	}

	// Respond with the newly created user (excluding password info)
	util.RespondJSON(w, http.StatusCreated, user)
}

// Login handles user login and token generation.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginData := struct {
		Credential string `json:"credential"` // Can be email or mobile number
		Password   string `json:"password"`   // Base64 encoded
	}{}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid login data")
		return
	}

	// Decode the base64-encoded password
	passwordBytes, err := base64.StdEncoding.DecodeString(loginData.Password)
	if err != nil {
		util.HandleError(w, http.StatusBadRequest, "Invalid password encoding")
		return
	}
	password := string(passwordBytes)

	// Determine if the credential is an email or a mobile number
	var user *models.User
	if util.IsValidEmail(loginData.Credential) {
		// Fetch user by email
		user, err = h.UserRepo.GetByField("email", loginData.Credential)
	} else if util.IsValidMobileNumber(loginData.Credential) {
		// Fetch user by mobile number
		user, err = h.UserRepo.GetByField("mobile_number", loginData.Credential)
	} else {
		util.HandleError(w, http.StatusUnauthorized, "Invalid email or mobile number")
		return
	}

	if err != nil {
		util.HandleError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Fetch stored password for user
	userPassword, err := h.UserPasswordRepo.GetByField("user_id", user.ID)
	if err != nil {
		util.HandleError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Validate password
	if err := models.ValidatePassword(password, userPassword.Salt, userPassword.Password); err != nil {
		util.HandleError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate token
	expiry := time.Now().Add(time.Hour * 72) // Token valid for 72 hours
	token := models.NewToken(user.ID, expiry)

	// Store the token in the database
	if err := h.TokenRepo.Create(token); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	// Respond with the generated token
	response := map[string]string{
		"token": token.Value.String(),
	}
	util.RespondJSON(w, http.StatusOK, &response)
}


// ValidateToken handles token validation and returns a JSON response
func (h *AuthHandler) ValidateToken (w http.ResponseWriter, r *http.Request) {
	// Extract the token from the "Token" header
	tokenHeader := r.Header.Get("Token")
	if tokenHeader == "" {
		// Respond with GenericResponseMessage
		response := models.GenericResponseMessage{
			Message: "Token is required",
			Result:  false,
		}
		util.RespondJSON(w, http.StatusUnauthorized, &response)
		return
	}

	// Validate the token
	token, err := h.TokenRepo.GetByField("value", tokenHeader)
	if err != nil {
		response := models.GenericResponseMessage{
			Message: "Invalid or expired token",
			Result:  false,
		}
		util.RespondJSON(w, http.StatusUnauthorized, &response)
		return
	}

	// Check if the token is expired
	if time.Now().After(token.Expiry) {
		response := models.GenericResponseMessage{
			Message: "Token has expired",
			Result:  false,
		}
		util.RespondJSON(w, http.StatusUnauthorized, &response)
		return
	}

	// Return success if the token is valid
	successResponse := struct {
		Message    string  `json:"message"`
		Result     bool    `json:"result"`
		UserID     uint64  `json:"user_id"`
		ExpiresIn  float64 `json:"expires_in"` // Time in seconds until expiry
	}{
		Message:   "Token is valid",
		Result:    true,
		UserID:    token.UserID,
		ExpiresIn: time.Until(token.Expiry).Seconds(),
	}

	util.RespondJSON(w, http.StatusOK, &successResponse)
}



// ValidateToken validates the token passed in the "Token" header and returns the user ID.
func ValidateToken(r *http.Request, tokenRepo *util.Repository[models.Token]) (uint64, error) {
	// Extract the token from the "Token" header
	tokenHeader := r.Header.Get("Token")
	if tokenHeader == "" {
		return 0, errors.New("missing token")
	}

	// Fetch the token from the database
	token, err := tokenRepo.GetByField("value", tokenHeader)
	if err != nil {
		return 0, errors.New("invalid or expired token")
	}

	// Check if the token is expired
	if time.Now().After(token.Expiry) {
		return 0, errors.New("expired token")
	}

	// Return the associated user ID
	return token.UserID, nil
}

// TokenValidationMiddleware checks for a valid token and attaches the user ID to the context.
func TokenValidationMiddleware(tokenRepo *util.Repository[models.Token]) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Validate the token and extract the user ID
			userID, err := ValidateToken(r, tokenRepo)
			if err != nil {
				util.HandleError(w, http.StatusUnauthorized, err.Error())
				return
			}

			// Set user ID in context for further use
			r = r.WithContext(util.ContextWithUserID(r.Context(), userID))

			next.ServeHTTP(w, r)
		})
	}
}
