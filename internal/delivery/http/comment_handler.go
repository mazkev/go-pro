package http

import (
	"ecommerce-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	usecase domain.CommentUsecase
}

func NewCommentHandler(usecase domain.CommentUsecase) *CommentHandler {
	return &CommentHandler{usecase: usecase}
}

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
}

// 1. POST /cards/:cardId/comments (Tambah komentar baru ke kartu)
func (h *CommentHandler) Create(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Isi komentar wajib diisi"})
		return
	}

	comment, err := h.usecase.CreateComment(uint(cardIDParam), userID, input.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Komentar berhasil ditambahkan!",
		"data":    comment,
	})
}

// 2. GET /cards/:cardId/comments (Ambil semua komentar pada suatu kartu)
func (h *CommentHandler) GetByCardID(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	comments, err := h.usecase.GetCommentsByCardID(uint(cardIDParam))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": comments,
	})
}

// 3. DELETE /comments/:id (Hapus komentar milik sendiri)
func (h *CommentHandler) Delete(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID komentar tidak valid"})
		return
	}

	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.usecase.DeleteComment(uint(idParam), userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Komentar berhasil dihapus!",
	})
}
