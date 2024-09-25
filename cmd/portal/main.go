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

func main() {
	// Initialize the database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Migrate the models (optional, if you want to auto-create tables)
	err = db.AutoMigrate(
		&models.Tenant{},
		&models.UserTenantMapping{},
		&models.User{},
		&models.UserPassword{},
		&models.Token{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Initialize handlers
	tenantHandler := v1.NewTenantHandler(db)
	authHandler := v1.NewAuthHandler(db)
	userHandler := v1.NewUserHandler(db)

	// Tenant-related routes
	http.HandleFunc("/tenants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tenantHandler.GetAllTenants(w, r) // Get all tenants
		case http.MethodPost:
			tenantHandler.CreateTenant(w, r) // Create a tenant
		}
	})

	http.HandleFunc("/tenants/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			tenantHandler.UpdateTenant(w, r) // Update a tenant
		}
	})

	http.HandleFunc("/tenants/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tenantHandler.GetTenantsByUser(w, r) // Get tenants by user
		}
	})

	http.HandleFunc("/tenants/mapping", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			tenantHandler.MapUserToTenant(w, r) // Map a single user to a tenant
		case http.MethodDelete:
			tenantHandler.DeleteUserTenantMapping(w, r) // Delete user-tenant mapping
		}
	})

	http.HandleFunc("/tenants/mappings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			tenantHandler.MapUsersToTenant(w, r) // Map multiple users to a tenant
		}
	})

	// Auth-related routes (registration, login)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.Register(w, r) // Register a new user
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authHandler.Login(w, r) // Login a user
		}
	})

	// User-related routes
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetAllUsers(w, r) // Get all users
		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetUserByID(w, r) // Get user by ID
		case http.MethodPut:
			userHandler.UpdateUser(w, r) // Update a user
		case http.MethodDelete:
			userHandler.DeleteUser(w, r) // Delete a user
		}
	})

	// Profile route (for getting the authenticated user's profile)
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authenticatedHandler(w, r, db, userHandler.GetUserProfile) // Protected route for profile
		}
	})

	// Change password route
	http.HandleFunc("/password/change", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authenticatedHandler(w, r, db, userHandler.ChangePassword) // Protected route for password change
		}
	})

	// Start the server on port 8080
	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

// authenticatedHandler is a wrapper to ensure the user is authenticated before calling a handler function
func authenticatedHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB, handler func(http.ResponseWriter, *http.Request)) {
	// Initialize the token repository for validation
	tokenRepo := util.NewRepository[models.Token](db)

	// Validate the token
	userID, err := v1.ValidateToken(r, tokenRepo)
	if err != nil {
		util.HandleError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Set the user ID in the context
	r = r.WithContext(util.ContextWithUserID(r.Context(), userID))

	// Call the original handler
	handler(w, r)
}

// initDB initializes the GORM database connection (using PostgreSQL)
func initDB() (*gorm.DB, error) {
	// Replace with your actual PostgreSQL connection string
	dsn := "host=localhost user=postgres password=314#sg dbname=sg-portal port=5432 sslmode=disable TimeZone=Asia/Kolkata"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
