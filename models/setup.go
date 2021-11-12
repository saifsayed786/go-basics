package models

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Settings struct {
	DB_HOST     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     string
}

func InitializeSettings() Settings {
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")

	switch {
	case DB_HOST == "":
		fmt.Println("1 Environmet variable DB_HOST not set.")
		os.Exit(1)
	case DB_NAME == "":
		fmt.Println("Environmet variable DB_NAME not set.")
		os.Exit(1)
	case DB_USER == "":
		fmt.Println("Environmet variable DB_USER not set.")
		os.Exit(1)
	case DB_PASSWORD == "":
		fmt.Println("Environmet variable DB_PASSWORD not set.")
		os.Exit(1)
	}

	settings := Settings{
		DB_HOST:     DB_HOST,
		DB_NAME:     DB_NAME,
		DB_USER:     DB_USER,
		DB_PASSWORD: DB_PASSWORD,
		DB_PORT:     DB_PORT,
	}

	return settings
}

func CreateInitialData() {

	var UserRoles = []UserRole{
		{
			Id:   1,
			Role: "admin",
		},
		{
			Id:   2,
			Role: "supervisor",
		},
		{
			Id:   3,
			Role: "customer",
		},
	}

	err := DB.CreateInBatches(UserRoles, 3).Error
	if err != nil {
		fmt.Println("User roles created successfully.")
	}

	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("SuperPassword@123"), 8)
	// Create an initial admin user.
	err = DB.Create(&User{
		ID:         1,
		FirstName:  "admin",
		LastName:   "admin",
		Email:      "admin@localhost.com",
		UserRoleID: 1,
		Password:   string(encryptedPassword),
	}).Error
	if err != nil {
		fmt.Println("Admin User already exists.")
	}

}

func ConnectDataBase() {
	// THis file is used to initialize the defined models
	settings := InitializeSettings()
	dsn := "host=" + settings.DB_HOST + " user=" + settings.DB_USER + " password=" + settings.DB_PASSWORD + " dbname=" + settings.DB_NAME + " port=" + settings.DB_PORT + " sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	// Register the models here.
	db.AutoMigrate(&Category{}, &User{}, &UserRole{}, &Product{})
	DB = db
	CreateInitialData()
}
