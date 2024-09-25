package v1

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"sg-portal/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestRegister tests the user registration endpoint with the new type field
func TestRegister(t *testing.T) {
    db := SetupTestDB(t)  // Set up the in-memory test database
    authHandler := NewAuthHandler(db)  // Initialize the handler

    // Prepare the registration payload
    payload := map[string]string{
        "email":        "test@example.com",
        "name":         "Test User",
        "mobile_number": "1234567890",
        "password":     base64.StdEncoding.EncodeToString([]byte("password123")),
        "type":         models.UserTypeClient,  // Register as a client user
    }
    body, _ := json.Marshal(payload)

    // Create a request to the register endpoint
    req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    // Execute the request
    rr := executeRequest(req, authHandler.Register)

    // Check the response code
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("Expected status code %d, got %d", http.StatusCreated, status)
    }

    // Check the response body contains the registered email
    var response map[string]interface{}
    json.Unmarshal(rr.Body.Bytes(), &response)

    if response["email"] != "test@example.com" {
        t.Errorf("Expected email to be 'test@example.com', got %v", response["email"])
    }
}

// SetupTestDB initializes an in-memory SQLite database for testing.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper() // Marks the function as a test helper, hiding it from test reports

	// Open an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the models (like User, UserPassword, Token, etc.)
	err = db.AutoMigrate(
		&models.User{},
		&models.UserPassword{},
		&models.Token{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}


// executeRequest helps simulate an HTTP request and records the response.
func executeRequest(req *http.Request, handlerFunc http.HandlerFunc) *httptest.ResponseRecorder {
	// Create a new ResponseRecorder to capture the response
	rr := httptest.NewRecorder()

	// Create a handler from the provided handler function
	handler := http.HandlerFunc(handlerFunc)

	// Serve the request using the handler
	handler.ServeHTTP(rr, req)

	// Return the ResponseRecorder, which contains the response details
	return rr
}
