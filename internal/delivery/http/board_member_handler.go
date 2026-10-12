package http

import (
	"ecommerce-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BoardMemberHandler struct {
	memberUsecase domain.BoardMemberUsecase
	cardUsecase   domain.CardUsecase
}

func NewBoardMemberHandler(memberUsecase domain.BoardMemberUsecase, cardUsecase domain.CardUsecase) *BoardMemberHandler {
	return &BoardMemberHandler{
		memberUsecase: memberUsecase,
		cardUsecase:   cardUsecase,
	}
}

type InviteMemberInput struct {
	Email string `json:"email" binding:"required"`
	Role  string `json:"role"` // "admin", "member", "viewer"
}

type UpdateRoleInput struct {
	Role string `json:"role" binding:"required"`
}

type AssignCardInput struct {
	AssigneeID   *uint  `json:"assigneeId"`
	AssigneeName string `json:"assigneeName"`
}

// 1. POST /boards/:boardId/members (Undang anggota baru ke board)
func (h *BoardMemberHandler) Invite(c *gin.Context) {
	boardIDParam, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID papan tidak valid"})
		return
	}

	currentUserID := getUserID(c)
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input InviteMemberInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email anggota wajib diisi"})
		return
	}

	member, err := h.memberUsecase.InviteMember(uint(boardIDParam), currentUserID, input.Email, input.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Anggota berhasil diundang",
		"data":    member,
	})
}

// 2. GET /boards/:boardId/members (Ambil daftar anggota board)
func (h *BoardMemberHandler) GetMembers(c *gin.Context) {
	boardIDParam, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID papan tidak valid"})
		return
	}

	currentUserID := getUserID(c)
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	members, err := h.memberUsecase.GetBoardMembers(uint(boardIDParam), currentUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": members,
	})
}

// 3. PUT /boards/:boardId/members/:userId (Ubah role anggota)
func (h *BoardMemberHandler) UpdateRole(c *gin.Context) {
	boardIDParam, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID papan tidak valid"})
		return
	}

	targetUserIDParam, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	currentUserID := getUserID(c)
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input UpdateRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role baru wajib diisi"})
		return
	}

	if err := h.memberUsecase.UpdateMemberRole(uint(boardIDParam), currentUserID, uint(targetUserIDParam), input.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Peran anggota berhasil diperbarui",
	})
}

// 4. DELETE /boards/:boardId/members/:userId (Keluarkan anggota dari board)
func (h *BoardMemberHandler) RemoveMember(c *gin.Context) {
	boardIDParam, err := strconv.Atoi(c.Param("boardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID papan tidak valid"})
		return
	}

	targetUserIDParam, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengguna tidak valid"})
		return
	}

	currentUserID := getUserID(c)
	if currentUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.memberUsecase.RemoveMember(uint(boardIDParam), currentUserID, uint(targetUserIDParam)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Anggota berhasil dikeluarkan dari papan",
	})
}

// 5. PUT /cards/:cardId/assign (Tugaskan anggota ke kartu tugas)
func (h *BoardMemberHandler) AssignCard(c *gin.Context) {
	cardIDParam, err := strconv.Atoi(c.Param("cardId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID kartu tidak valid"})
		return
	}

	var input AssignCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data penugasan tidak valid"})
		return
	}

	card, err := h.cardUsecase.AssignCard(uint(cardIDParam), input.AssigneeID, input.AssigneeName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Penugasan anggota kartu berhasil diperbarui",
		"data":    card,
	})
}
