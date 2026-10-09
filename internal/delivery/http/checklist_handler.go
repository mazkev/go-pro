package http

import (
	"ecommerce-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChecklistHandler struct {
	usecase domain.ChecklistUsecase
}

func NewChecklistHandler(usecase domain.ChecklistUsecase) *ChecklistHandler {
	return &ChecklistHandler{usecase: usecase}
}

type CreateChecklistInput struct {
	Title string `json:"title" binding:"required"`
}

type UpdateChecklistInput struct {
	Title   string `json:"title"`
	Checked bool   `json:"checked"`
}

// 1. POST /cards/:cardId/checklist (Tambah sub-tugas)
func (h *ChecklistHandler) Create(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	var input CreateChecklistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul sub-tugas wajib diisi"})
		return
	}

	item, err := h.usecase.CreateItem(uint(cardIDParam), input.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Sub-tugas berhasil ditambahkan!",
		"data":    item,
	})
}

// 2. PUT /checklist/:id (Centang selesai atau edit teks)
func (h *ChecklistHandler) Update(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sub-tugas tidak valid"})
		return
	}

	var input UpdateChecklistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	item, err := h.usecase.UpdateItem(uint(idParam), input.Title, input.Checked)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sub-tugas berhasil diperbarui",
		"data":    item,
	})
}

// 3. DELETE /checklist/:id (Hapus sub-tugas)
func (h *ChecklistHandler) Delete(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sub-tugas tidak valid"})
		return
	}

	if err := h.usecase.DeleteItem(uint(idParam)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sub-tugas berhasil dihapus"})
}
