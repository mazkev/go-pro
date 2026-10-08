package http

import (
	"ecommerce-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	usecase domain.CardUsecase
}

func NewCardHandler(usecase domain.CardUsecase) *CardHandler {
	return &CardHandler{usecase: usecase}
}

type CreateCardInput struct {
	Content     string `json:"content" binding:"required"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type UpdateCardInput struct {
	Content     string `json:"content"`
	Description string `json:"description"`
	Label       string `json:"label"`
	DueDate     string `json:"dueDate"`
	Checked     bool   `json:"checked"`
}

type MoveCardInput struct {
	TargetListID uint `json:"targetListId" binding:"required"`
	NewPosition  int  `json:"newPosition"`
}

// 1. POST /lists/:listId/cards (Tambah kartu ke kolom tertentu)
func (h *CardHandler) Create(c *gin.Context) {
	listIDParam, err := strconv.Atoi(c.Param("listId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID list tidak valid"})
		return
	}

	var input CreateCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Konten/judul kartu wajib diisi"})
		return
	}

	card, err := h.usecase.CreateCard(uint(listIDParam), input.Content, input.Label, input.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kartu berhasil ditambahkan ke kolom!",
		"data":    card,
	})
}

// 2. PUT /cards/:id (Update data kartu)
func (h *CardHandler) Update(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	var input UpdateCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.usecase.UpdateCard(uint(cardIDParam), input.Content, input.Description, input.Label, input.DueDate, input.Checked)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": card})
}

// 3. PUT /cards/:id/move (Drag & Drop: Pindah kolom / ganti urutan)
func (h *CardHandler) Move(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	var input MoveCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kolom tujuan wajib diisi"})
		return
	}

	if err := h.usecase.MoveCard(uint(cardIDParam), input.TargetListID, input.NewPosition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kartu berhasil dipindahkan!"})
}

// 4. DELETE /cards/:id (Hapus kartu)
func (h *CardHandler) Delete(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	if err := h.usecase.DeleteCard(uint(cardIDParam)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kartu berhasil dihapus"})
}
