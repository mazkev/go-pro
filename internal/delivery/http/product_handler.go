package http

import (
	"ecommerce-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProductHandler sekarang HANYA memegang Usecase (tidak kenal database GORM sama sekali)
type ProductHandler struct {
	usecase domain.ProductUsecase
}

// Inisialisasi Handler baru dengan menyuntikkan Usecase
func NewProductHandler(usecase domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: usecase}
}

// 1. GET /products
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// 2. GET /products/:id
func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
		return
	}

	product, err := h.usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

	// Serahkan validasi dan penyimpanan ke Usecase!
	if err := h.usecase.Create(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk berhasil dibuat!",
		"data":    newProduct,
	})
}

// 4. PUT /products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
		return
	}

	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	updatedProduct, err := h.usecase.Update(id, &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil diperbarui!",
		"data":    updatedProduct,
	})
}

// 5. DELETE /products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID harus berupa angka"})
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil dihapus!",
	})
}
// 6. GET /products/search?keyword=...
func (h *ProductHandler) Search(c *gin.Context) {
	// Ambil kata kunci dari URL Query Parameter (?keyword=...)
	keyword := c.Query("keyword")

	// Panggil usecase
	products, err := h.usecase.Search(keyword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Kembalikan hasil pencarian ke user
	c.JSON(http.StatusOK, products)
}
