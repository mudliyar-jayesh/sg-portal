package main

import (
	"log"
	repository "sg-portal/pkg/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Define the User model
type User struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Email string `gorm:"unique;not null"`
	Age   int
}

// Initialize PostgreSQL database connection
func ConnectDB() (*gorm.DB, error) {
	// Update the connection string (DSN) as per your database credentials
	dsn := "host=localhost user=postgres password=314#sg dbname=sg-portal port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// Connect to the database
	db, err := ConnectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Auto-migrate the schema for User model
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to auto-migrate schema:", err)
	}
	log.Println("Database schema migrated successfully")

	// Create a repository instance for User model
	userRepo := repository.NewRepository[User](db)

	// Example: Create a new user using the repository helper
	newUser := User{Name: "John", Email: "John@example.com", Age: 25}
	if err := userRepo.Create(&newUser); err != nil {
		log.Fatal("Failed to create user:", err)
	}
	log.Println("New user created:", newUser)

	// Example: Fetch a user by email using the repository helper
	user, err := userRepo.GetByField("email", "John@example.com")
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}
	log.Println("User found:", user)

	// Example: Update user's age using the repository helper
	updateData := map[string]interface{}{"Age": 22}
	err = userRepo.UpdateOne("email", "John@example.com", updateData)
	if err != nil {
		log.Fatal("Failed to update user:", err)
	}
	log.Println("User updated successfully")

	// Example: Get all users whose age is greater than 20
	users, err := userRepo.GetAllByCondition("age > ?", 20)
	if err != nil {
		log.Fatal("Failed to get users:", err)
	}
	log.Println("Users found:", users)
}
