package controllers

import (
	"encoding/json"
	"fmt"
	"go-basics/models"
	"html/template"
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// ALl controller logic related to the users .
// Available endpoints - Signup, Profile, etc

// TempUser is a temperory struct to validate the user input.

type TempUser struct {
	Email           string `binding:"required"`
	FirstName       string `binding:"required"`
	LastName        string `binding:"required"`
	Password        string `binding:"required"`
	ConfirmPassword string `binding:"required"`
}

// Returns a custom user input
func ReturnParameterMissingError(c *gin.Context, parameter string) {
	var err = fmt.Sprintf("Required parameter %s missing.", parameter)
	c.JSON(http.StatusBadRequest, gin.H{"error": err})
}

// @Summary Signup endpoint is used for customer registeration. ( Supervisors/admin can be added only by admin. )
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/signup [post]
// @Tags auth
// @Accept json
// @Produce json
// @Param email formData string true "Email of the user"
// @Param first_name formData string true "First name of the user"
// @Param last_name formData string true "Last name of the user"
// @Param password formData string true "Password of the user"
// @Param confirm_password formData string true "Confirm password."

func Signup(c *gin.Context) {
	var tempUser TempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"email", "first_name", "last_name", "password", "confirm_password"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirm_password"))

	// Check if both passwords are same.
	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
	}

	// check if the password is strong and matches the password policy
	// length > 8, atleast 1 upper case, atleast 1 lower case, atleast 1 symbol
	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if ispasswordstrong == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "customer").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id,
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})
}

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Signin godoc
// @Summary Signin endpoint is used by the user to login.
// @Description API Endpoint to register the user with the role of customer.
// @Router /api/v1/signin [post]
// @Tags auth
// @Accept json
// @Produce json
// @Param login formData login true "Credentials of the user"

func Signin(c *gin.Context) (interface{}, error) {
	/* THIS FUNCTION IS USED BY THE AUTH MIDDLEWARE. */
	var loginVals login
	// var User User
	var users []models.User
	var count int64
	// var user models.User
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	// First check if the user exist or not...
	models.DB.Where("email = ?", email).Find(&users).Count(&count)
	if count == 0 {
		return nil, jwt.ErrFailedAuthentication
	}
	if CheckCredentials(loginVals.Email, loginVals.Password, models.DB) == true {
		return &models.User{
			Email: email,
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

// CreateSupervisorOrAdmin godoc
// @Summary CreateSupervisor endpoint is used by the admin role user to create a new admin or supervisor account.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/supervisor/create [post]
// @Tags supervisor
// @Accept json
// @Produce json
// @Param login formData TempUser true "Info of the user"

func CreateSupervisor(c *gin.Context) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	// Create a user with the role of supervisor.
	var tempUser TempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"email", "first_name", "last_name", "password", "confirm_password"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirm_password"))

	// Check if both passwords are same.
	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
	}

	// check if the password is strong and matches the password policy
	// length > 8, atleast 1 upper case, atleast 1 lower case, atleast 1 symbol
	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if ispasswordstrong == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "supervisor").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registeration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

	return
}

// CreateAdmin godoc
// @Summary CreateAdmin endpoint is used by the admin role user to create a new admin or supervisor account.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/admin/create [post]
// @Tags admin
// @Accept json
// @Produce json
// @Param login formData TempUser true "Info of the user"
func CreateAdmin(c *gin.Context) {

	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	// Create a user with the role of supervisor.
	var tempUser TempUser
	var Role models.UserRole

	c.Request.ParseForm()
	paramList := []string{"email", "first_name", "last_name", "password", "confirm_password"}

	for _, param := range paramList {
		if c.PostForm(param) == "" {
			ReturnParameterMissingError(c, param)
		}
	}

	tempUser.Email = template.HTMLEscapeString(c.PostForm("email"))
	tempUser.FirstName = template.HTMLEscapeString(c.PostForm("first_name"))
	tempUser.LastName = template.HTMLEscapeString(c.PostForm("last_name"))
	tempUser.Password = template.HTMLEscapeString(c.PostForm("password"))
	tempUser.ConfirmPassword = template.HTMLEscapeString(c.PostForm("confirm_password"))

	// Check if both passwords are same.
	if tempUser.Password != tempUser.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both passwords do not match."})
	}

	// check if the password is strong and matches the password policy
	// length > 8, atleast 1 upper case, atleast 1 lower case, atleast 1 symbol
	ispasswordstrong, _ := IsPasswordStrong(tempUser.Password)
	if ispasswordstrong == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is not strong."})
		return
	}

	// Check if the user already exists.
	if DoesUserExist(tempUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists."})
		return
	}

	encryptedPassword, error := HashPassword(tempUser.Password)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}

	err := models.DB.Where("role= ?", "admin").First(&Role).Error
	if err != nil {
		fmt.Println("err ", err.Error())
		return
	}

	SanitizedUser := models.User{
		FirstName:  tempUser.FirstName,
		LastName:   tempUser.LastName,
		Email:      tempUser.Email,
		Password:   encryptedPassword,
		UserRoleID: Role.Id, //This endpoint will be used only for customer registeration.
		CreatedAt:  time.Now(),
		IsActive:   true,
	}

	errs := models.DB.Create(&SanitizedUser).Error
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Some error occoured."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"msg": "User created successfully"})

	return
}
