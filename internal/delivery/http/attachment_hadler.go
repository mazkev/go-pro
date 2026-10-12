package http

import (
	"ecommerce-api/internal/domain"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	usecase domain.AttachmentUsecase
}

func NewAttachmentHandler(usecase domain.AttachmentUsecase) *AttachmentHandler {
	return &AttachmentHandler{usecase: usecase}
}

// 1. POST /cards/:cardId/attachments (Upload file lampiran ke kartu)
func (h *AttachmentHandler) Upload(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	// Baca file dari multipart form input bernama "file"
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File wajib diunggah (form field: 'file')"})
		return
	}

	// Batas ukuran file: Maksimal 10 MB (10 * 1024 * 1024 bytes)
	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar (maksimal 10 MB)"})
		return
	}

	// Buat direktori ./uploads/attachments jika belum ada
	uploadDir := filepath.Join(".", "uploads", "attachments")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyiapkan folder penyimpanan"})
		return
	}

	// Beri nama file unik agar tidak saling menimpa
	cleanFileName := filepath.Base(file.Filename)
	uniqueFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), cleanFileName)
	destination := filepath.Join(uploadDir, uniqueFileName)

	// Simpan file fisik ke disk
	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file ke disk"})
		return
	}

	// Ambil tipe MIME file (e.g. "image/png", "application/pdf")
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		fileType = "application/octet-stream"
	}

	// Siapkan objek Attachment untuk disimpan ke database
	attachment := &domain.Attachment{
		CardID:   uint(cardIDParam),
		FileName: cleanFileName,
		FilePath: "/uploads/attachments/" + uniqueFileName, // URL path yang bisa diakses publik
		FileSize: file.Size,
		FileType: fileType,
	}

	if err := h.usecase.CreateAttachment(attachment); err != nil {
		// Jika gagal simpan DB, hapus kembali file fisiknya
		_ = os.Remove(destination)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "File berhasil diunggah",
		"data":    attachment,
	})
}

// 2. GET /cards/:cardId/attachments (Ambil daftar lampiran milik kartu)
func (h *AttachmentHandler) GetByCardID(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	attachments, err := h.usecase.GetAttachmentsByCardID(uint(cardIDParam))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": attachments,
	})
}

// 3. DELETE /attachments/:id (Hapus lampiran dari kartu & server)
func (h *AttachmentHandler) Delete(c *gin.Context) {
	idParam, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID lampiran tidak valid"})
		return
	}

	if err := h.usecase.DeleteAttachment(uint(idParam)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lampiran berhasil dihapus",
	})
}
