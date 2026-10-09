package http

import (
	"ecommerce-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ListHandler struct {
	usecase domain.ListUsecase
}

func NewListHandler(usecase domain.ListUsecase) *ListHandler {
	return &ListHandler{usecase: usecase}
}

type CreateListInput struct {
	Title string `json:"title" binding:"required"`
}

type UpdateListInput struct {
	Title string `json:"title" binding:"required"`
}

type MoveListInput struct {
	Position int `json:"position"`
}

// 1. POST /boards/:boardId/lists (Tambah kolom ke board)
func (h *ListHandler) Create(c *gin.Context) {
	boardIDParam, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID board tidak valid"})
		return
	}

	var input CreateListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul kolom wajib diisi"})
		return
	}

	list, err := h.usecase.CreateList(uint(boardIDParam), input.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kolom berhasil ditambahkan!",
		"data":    list,
	})
}

// 2. PUT /lists/:id (Rename kolom)
func (h *ListHandler) Update(c *gin.Context) {
	listIDParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID list tidak valid"})
		return
	}

	var input UpdateListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Judul kolom tidak boleh kosong"})
		return
	}

	list, err := h.usecase.UpdateList(uint(listIDParam), input.Title)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kolom berhasil diperbarui",
		"data":    list,
	})
}

// 3. PUT /lists/:id/move (Geser urutan posisi kolom)
func (h *ListHandler) Move(c *gin.Context) {
	listIDParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID list tidak valid"})
		return
	}

	var input MoveListInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	if err := h.usecase.MoveList(uint(listIDParam), input.Position); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Posisi kolom berhasil diubah"})
}

// 4. DELETE /lists/:id (Hapus kolom beserta seluruh kartu di dalamnya)
func (h *ListHandler) Delete(c *gin.Context) {
	listIDParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID list tidak valid"})
		return
	}

	if err := h.usecase.DeleteList(uint(listIDParam)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kolom berhasil dihapus"})
}
