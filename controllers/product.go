package controllers

import (
	"go-basics/models"
	"html/template"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// CreateProduct godoc
// @Summary CreateProduct endpoint is used by the supervisor role user to create a new product.
// @Description CreateProduct endpoint is used by the supervisor role user to create a new product
// @Router /api/v1/auth/product/create [post]
// @Tags product
// @Accept json
// @Produce json
// @Param name formData string true "name of the product"
// @Param category_id formData int true "id of the category"
func CreateProduct(c *gin.Context) {

	var existingProduct models.Product
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var category models.Category

	// Check if the current user had admin role.
	if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be added by supervisor user"})
		return
	}

	c.Request.ParseForm()

	if c.PostForm("name") == "" {
		ReturnParameterMissingError(c, "name")
		return
	}
	if c.PostForm("category_id") == "" {
		ReturnParameterMissingError(c, "category_id")
		return
	}

	product_title := template.HTMLEscapeString(c.PostForm("name"))
	category_id := template.HTMLEscapeString(c.PostForm("category_id"))

	// Check if the product already exists.
	err := models.DB.Where("title = ?", product_title).First(&existingProduct).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product already exists."})
		return
	}

	// Check if the category exists
	err = models.DB.First(&category, category_id).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not exists."})
		return
	}
	cat := models.Product{
		Title:      product_title,
		CategoryId: category.ID,
		CreatedBy:  User.ID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = models.DB.Create(&cat).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":   cat.ID,
		"name": cat.Title,
	})

}

// UpdateProduct godoc
// @Summary UpdateProduct endpoint is used by the supervisor role user to update a new product.
// @Description UpdateProduct endpoint is used by the supervisor role user to update a new product
// @Router /api/v1/auth/product/:id/ [PATCH]
// @Tags product
// @Accept json
// @Produce json
// @Param name formData string true "name of the product"
func UpdateProduct(c *gin.Context) {
	var existingProduct models.Product
	var updateProduct models.Product
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User

	// Check if the current user had admin role.
	if err := models.DB.Where("email = ? AND user_role_id=2", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product can only be updated by supervisor user"})
		return
	}

	// Check if the product already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingProduct).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product doesnot exists."})
		return
	}

	if err := c.ShouldBindJSON(&updateProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&existingProduct).Updates(updateProduct)

}

type ReturnedProduct struct {
	ID         int    `json:"id,string"`
	Title      string `json:"name"`
	CategoryId int    `json:"category_id"`
}

// GetProduct godoc
// @Summary GetProduct endpoint is used to get info of a product..
// @Description GetProduct endpoint is used to get info of a product.
// @Router /api/v1/auth/product/:id/ [get]
// @Tags product
// @Accept json
// @Produce json
// @Param name formData string true "name of the product"
func GetProduct(c *gin.Context) {
	var existingProduct models.Product

	// Check if the product already exists.
	err := models.DB.Where("id = ?", c.Param("id")).First(&existingProduct).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product doesnot exists."})
		return
	}

	// GET FROM CACHE FIRST
	c.JSON(http.StatusOK, gin.H{"product": existingProduct})
	return
}

// ListAllProduct godoc
// @Summary ListAllProduct endpoint is used to list all products.
// @Description API Endpoint to register the user with the role of Supervisor or Admin.
// @Router /api/v1/auth/product/ [get]
// @Tags product
// @Accept json
// @Produce json
func ListAllProduct(c *gin.Context) {

	// allProduct := []models.Product{}
	claims := jwt.ExtractClaims(c)
	user_email, _ := claims["email"]
	var User models.User
	var Product []models.Product
	var ExistingProduct []ReturnedProduct

	if err := models.DB.Where("email = ?", user_email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	models.DB.Model(Product).Find(&ExistingProduct)
	c.JSON(http.StatusOK, ExistingProduct)
	return
}
