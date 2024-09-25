package main

import (
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"sg-portal/api/v1"
	"sg-portal/internal/models"
)

func main() {
	// Initialize the database connection
	db, err := initDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// Migrate the models (optional, if you want to auto-create tables)
	err = db.AutoMigrate(&models.Tenant{}, &models.UserTenantMapping{}, &models.Feature{}, &models.UserFeatureMapping{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Initialize the TenantHandler
	tenantHandler := v1.NewTenantHandler(db)

	// Initialize the FeatureHandler
	featureHandler := v1.NewFeatureHandler(db)

	// Set up routes for the Tenant API endpoints
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

	// Set up routes for the Feature API endpoints
	http.HandleFunc("/features", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			featureHandler.GetAllFeatures(w, r) // Get all features
		case http.MethodPost:
			featureHandler.CreateFeature(w, r) // Create a feature
		}
	})

	http.HandleFunc("/features/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			featureHandler.UpdateFeature(w, r) // Update a feature
		}
	})

	http.HandleFunc("/features/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			featureHandler.GetFeaturesByUser(w, r) // Get features by user
		}
	})

	http.HandleFunc("/features/mapping", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			featureHandler.MapUserToFeature(w, r) // Map a single user to a feature
		case http.MethodDelete:
			featureHandler.DeleteUserFeatureMapping(w, r) // Delete user-feature mapping
		}
	})

	http.HandleFunc("/features/mappings", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			featureHandler.MapUsersToFeature(w, r) // Map multiple users to a feature
		}
	})

	// Start the server on port 8080
	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
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

