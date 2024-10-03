package v1

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sg-portal/internal/models"
	"sg-portal/pkg/util"
	"time"

	"gorm.io/gorm"
)

type AuthHandler struct {
	UserRepo          *util.Repository[models.User]
	UserPasswordRepo  *util.Repository[models.UserPassword]
	TokenRepo         *util.Repository[models.Token]
	TenantRepo        *util.Repository[models.Tenant]
	TenantMappingRepo *util.Repository[models.UserTenantMapping]
}

// NewAuthHandler initializes the auth handler with the repositories.
func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		UserRepo:          util.NewRepository[models.User](db),
		UserPasswordRepo:  util.NewRepository[models.UserPassword](db),
		TokenRepo:         util.NewRepository[models.Token](db),
		TenantRepo:        util.NewRepository[models.Tenant](db),
		TenantMappingRepo: util.NewRepository[models.UserTenantMapping](db),
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

	// map default tenant
	tenantsRepo := util.NewRepository[models.Tenant](util.Db)
	defaultTenant, tenantErr := tenantsRepo.GetByField("company_name", "default")
	if tenantErr != nil || defaultTenant == nil {
		util.HandleError(w, http.StatusInternalServerError, "Could not map to demo server")
	}

	tenantMapping := models.UserTenantMapping{
		UserId:   user.ID,
		TenantId: defaultTenant.ID,
	}

	tenantsMappingRepo := util.NewRepository[models.UserTenantMapping](util.Db)
	if err := tenantsMappingRepo.Create(&tenantMapping); err != nil {
		util.HandleError(w, http.StatusInternalServerError, "Could not map to demo server")
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

	// Respond with the
	response := struct {
		User  *models.User `json:"user_info"`
		Token string       `json:"token"`
	}{
		Token: token.Value.String(),
		User:  user,
	}
	util.RespondJSON(w, http.StatusOK, &response)
}

// validate token and resolve tenant

func (h *AuthHandler) ResolveTenant(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	companyId := r.Header.Get("companyid")

	tokenRepo := util.NewRepository[models.Token](util.Db)
	tokenInfo, err := tokenRepo.GetByField("value", token)
	if err != nil {
		// Respond with GenericResponseMessage
		response := models.TokenTenantInfo{
			Message: "Invalid Token Provided",
			Success: false,
		}
		util.RespondJSON(w, http.StatusUnauthorized, &response)
		return
	}

	tenantRepo := util.NewRepository[models.Tenant](util.Db)
	tenantInfo, err := tenantRepo.GetByField("company_guid", companyId)
	if err != nil {
		// Respond with GenericResponseMessage
		response := models.TokenTenantInfo{
			Message: "Non Registered Company Requested",
			Success: false,
		}
		util.RespondJSON(w, http.StatusUnauthorized, &response)
		return
	}
	tenantMappingRepo := util.NewRepository[models.UserTenantMapping](util.Db)
	tenantMapping, err := tenantMappingRepo.GetAllByCondition("user_id = ? and tenant_id = ?", tokenInfo.UserID, tenantInfo.ID)
	if err != nil || len(tenantMapping) < 1 {
		// Respond with GenericResponseMessage
		response := models.TokenTenantInfo{
			Message: "No Tenants Configured for the user",
			Success: false,
		}
		util.RespondJSON(w, http.StatusUnauthorized, &response)
		return
	}

	response := models.TokenTenantInfo{
		TenantInfo: tenantInfo,
		UserId:     &tokenInfo.UserID,
		Message:    "Token Valid",
		Success:    true,
	}
	util.RespondJSON(w, http.StatusOK, &response)
}
