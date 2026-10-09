package http

import (
	"ecommerce-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BoardHandler struct {
	usecase domain.BoardUsecase
}

func NewBoardHandler(usecase domain.BoardUsecase) *BoardHandler {
	return &BoardHandler{usecase: usecase}
}

// Helper aman untuk mengekstrak user_id dari JWT tanpa panic
func getUserID(c *gin.Context) uint {
	val, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	switch v := val.(type) {
	case int:
		return uint(v)
	case uint:
		return v
	case float64:
		return uint(v)
	default:
		return 0
	}
}

type CreateBoardInput struct {
	Title string `json:"title" binding:"required"`
	Color string `json:"color"`
}

// 1. GET /boards (Ambil semua board milik user yang sedang login)
func (h *BoardHandler) GetMyBoards(c *gin.Context) {
	userID := getUserID(c)

	boards, err := h.usecase.GetMyBoards(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": boards})
}

// 2. GET /boards/:id (Ambil detail 1 board beserta susunan kolom dan kartu)
func (h *BoardHandler) GetByID(c *gin.Context) {
	userID := getUserID(c)
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID board tidak valid"})
		return
	}

	board, err := h.usecase.GetBoardDetail(uint(idParam), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": board})
}

// 3. POST /boards (Buat board baru)
func (h *BoardHandler) Create(c *gin.Context) {
	userID := getUserID(c)

	var input CreateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul board wajib diisi"})
		return
	}

	board, err := h.usecase.CreateBoard(userID, input.Title, input.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Board berhasil dibuat lengkap dengan kolom bawaan!",
		"data":    board,
	})
}

// 4. DELETE /boards/:id (Hapus board)
func (h *BoardHandler) Delete(c *gin.Context) {
	userID := getUserID(c)
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID board tidak valid"})
		return
	}

	if err := h.usecase.DeleteBoard(uint(idParam), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Board berhasil dihapus beserta seluruh isinya"})
}

type UpdateBoardInput struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

// 5. PUT /boards/:id (Update Judul / Warna Board)
func (h *BoardHandler) Update(c *gin.Context) {
	userID := getUserID(c)
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID board tidak valid"})
		return
	}

	var input UpdateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	board, err := h.usecase.UpdateBoard(uint(idParam), userID, input.Title, input.Color)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Board berhasil diperbarui",
		"data":    board,
	})
}
