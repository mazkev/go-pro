package http

import (
	"ecommerce-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductHandler menampung koneksi database agar bisa dipakai semua fungsi di bawah
type ProductHandler struct {
	db *gorm.DB
}

// Fungsi pembantu untuk membuat ProductHandler baru
func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// 1. GET /products
func (h *ProductHandler) GetAll(c *gin.Context) {
	var products []domain.Product
	h.db.Find(&products)
	c.JSON(http.StatusOK, products)
}

// 2. GET /products/:id
func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus angka"})
		return
	}

	var product domain.Product
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// 3. POST /products
func (h *ProductHandler) Create(c *gin.Context) {
	var newProduct domain.Product
	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	h.db.Create(&newProduct)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk berhasil disimpan ke database!",
		"data":    newProduct,
	})
}

// 4. PUT /products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus angka"})
		return
	}

	var product domain.Product
	if err := h.db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	h.db.Model(&product).Updates(domain.Product{
		Name:  input.Name,
		Price: input.Price,
		Stock: input.Stock,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil diperbarui!",
		"data":    product,
	})
}

// 5. DELETE /products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus angka"})
		return
	}

	result := h.db.Delete(&domain.Product{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan untuk dihapus"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil dihapus dari database!",
	})
}
