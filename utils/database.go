package utils

import (
	"com.say.more.server/config"
	"com.say.more.server/internal/app/models/db"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.MySQL.User,
		config.AppConfig.MySQL.Password,
		config.AppConfig.MySQL.Host,
		config.AppConfig.MySQL.Port,
		config.AppConfig.MySQL.DBName,
	)

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable detailed logging
	}

	// Connect to the database
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established successfully")

	// Auto migrate database tables
	err = migrateDB()
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed initial data
	err = seedDB()
	if err != nil {
		log.Fatal("Failed to seed database:", err)
	}
}

// migrateDB auto migrates database tables
func migrateDB() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&db.User{},
		// Add other models here
	)

	if err != nil {
		return fmt.Errorf("auto migration failed: %v", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// seedDB seeds initial data
func seedDB() error {
	log.Println("Checking if database needs seeding...")

	// Check if admin user already exists
	var count int64
	DB.Model(&db.User{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping...")
		return nil
	}

	// Create default admin user
	adminUser := db.User{
		Username: "admin",
		Password: "admin123", // Note: Use encrypted password in production
		Email:    "admin@example.com",
		Nickname: "Administrator",
		Status:   1,
	}

	// Create test user
	testUser := db.User{
		Username: "test_user",
		Password: "test123", // Note: Use encrypted password in production
		Email:    "test@example.com",
		Nickname: "Test User",
		Status:   1,
	}

	// Create initial users in a transaction
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&adminUser).Error; err != nil {
			return err
		}
		if err := tx.Create(&testUser).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("database seeding failed: %v", err)
	}

	log.Println("Database seeded successfully")
	return nil
}
