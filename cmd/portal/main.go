package main

import (
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	v1 "sg-portal/api/v1"
	"sg-portal/internal/models"
	"sg-portal/pkg/util"
)

// corsMiddleware handles the CORS settings for incoming requests
func corsMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Priority, companyid, token")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Printf("[#] Starting Server...\n")
	// Initialize the database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	util.Db = db

	// Migrate the models
	err = db.AutoMigrate(
		&models.Tenant{}, &models.UserTenantMapping{},
		&models.User{}, &models.UserPassword{}, &models.Token{},
		&models.Feature{}, &models.UserFeatureMapping{},
		&models.Subscription{}, &models.UserSubscriptionMapping{},
		&models.FeatureSubscriptionMapping{}, &models.UserSubscriptionHistory{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Initialize handlers
	authHandler := v1.NewAuthHandler(db)
	userHandler := v1.NewUserHandler(db)
	tenantHandler := v1.NewTenantHandler(db)
	featureHandler := v1.NewFeatureHandler(db)
	subscriptionHandler := v1.NewSubscriptionHandler(db)
	userSubscriptionHistoryHandler := v1.NewUserSubscriptionHistoryHandler(db)
	companyHandler := v1.NewCompanyHandler(db)

	// Define a new ServeMux to register routes
	mux := http.NewServeMux()

	// company-related routes
	mux.HandleFunc("/companies", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			companyHandler.GetCompanies(w, r)
		}
	})

	// Tenant-related routes
	mux.HandleFunc("/tenants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tenantHandler.GetAllTenants(w, r)
		case http.MethodPost:
			tenantHandler.CreateTenant(w, r)
		}
	})

	mux.HandleFunc("/tenants/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			tenantHandler.UpdateTenant(w, r)
		}
	})

	mux.HandleFunc("/tenants/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tenantHandler.GetTenantsByUser(w, r)
		}
	})

	// Auth-related routes
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.Register(w, r)
		}
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.Login(w, r)
		}
	})

	mux.HandleFunc("/token/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authHandler.ValidateToken(w, r)
		}
	})

	// User-related routes
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			userHandler.GetAllUsers(w, r)
		}
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUserByID(w, r)
		case http.MethodPut:
			userHandler.UpdateUser(w, r)
		case http.MethodDelete:
			userHandler.DeleteUser(w, r)
		}
	})

	// Profile route (for getting the authenticated user's profile)
	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authenticatedHandler(w, r, db, userHandler.GetUserProfile)
		}
	})

	// Change password route
	mux.HandleFunc("/password/change", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authenticatedHandler(w, r, db, userHandler.ChangePassword)
		}
	})

	// Set up routes for the Feature API
	mux.HandleFunc("/features", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			featureHandler.GetAllFeatures(w, r)
		case http.MethodPost:
			featureHandler.CreateFeature(w, r)
		}
	})

	mux.HandleFunc("/features/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			featureHandler.UpdateFeature(w, r)
		}
	})

	mux.HandleFunc("/features/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			featureHandler.GetFeaturesByUser(w, r)
		}
	})

	// Set up routes for the Subscription API
	mux.HandleFunc("/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			subscriptionHandler.GetAllSubscriptions(w, r)
		case http.MethodPost:
			subscriptionHandler.CreateSubscription(w, r)
		}
	})

	mux.HandleFunc("/subscriptions/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			subscriptionHandler.UpdateSubscription(w, r)
		}
	})

	mux.HandleFunc("/subscriptions/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			subscriptionHandler.GetSubscriptionsByUser(w, r)
		}
	})

	// Set up routes for UserSubscriptionHistory API
	mux.HandleFunc("/subscriptions/history", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Query().Has("userId") {
				userSubscriptionHistoryHandler.GetUserSubscriptionHistory(w, r)
			} else {
				userSubscriptionHistoryHandler.GetAllUserSubscriptionHistories(w, r)
			}
		case http.MethodPost:
			userSubscriptionHistoryHandler.CreateUserSubscriptionHistory(w, r)
		}
	})

	// Start the server on port 8080 with CORS-enabled middleware
	corsMux := corsMiddleware(mux)
	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", corsMux); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// authenticatedHandler ensures the user is authenticated before calling a handler function
func authenticatedHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB, handler func(http.ResponseWriter, *http.Request)) {
	tokenRepo := util.NewRepository[models.Token](db)

	userID, err := v1.ValidateToken(r, tokenRepo)
	if err != nil {
		util.HandleError(w, http.StatusUnauthorized, err.Error())
		return
	}

	r = r.WithContext(util.ContextWithUserID(r.Context(), userID))
	handler(w, r)
}

// initDB initializes the GORM database connection using PostgreSQL
func initDB() (*gorm.DB, error) {
	log.Printf("[+] Connecting to Postgres...\n")
	dsn := "host=192.168.1.36 user=postgres password=314#sg dbname=sg-portal port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Printf("[!] Connected.\n")
	return db, nil
}
