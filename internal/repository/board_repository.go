package repository

import (
	"ecommerce-api/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type boardRepository struct {
	db *gorm.DB
}

// Constructor pembuat instance repository
func NewBoardRepository(db *gorm.DB) domain.BoardRepository {
	return &boardRepository{db: db}
}

// 1. Ambil semua board milik user atau di mana user terdaftar sebagai anggota
func (r *boardRepository) FindAllByUserID(userID uint) ([]domain.Board, error) {
	var boards []domain.Board
	// Nested Preload: Otomatis load Members, Lists, Cards, Checklist, Comments, & Attachments!
	err := r.db.Preload("Members").
		Preload("Lists.Cards.Checklist").
		Preload("Lists.Cards.Comments").
		Preload("Lists.Cards.Attachments").
		Where("user_id = ? OR id IN (SELECT board_id FROM board_members WHERE user_id = ?)", userID, userID).
		Find(&boards).Error
	return boards, err
}

// 2. Ambil 1 board detail berdasarkan ID jika user adalah pemilik atau anggota
func (r *boardRepository) FindByID(id uint, userID uint) (*domain.Board, error) {
	var board domain.Board
	err := r.db.Preload("Members").
		Preload("Lists.Cards.Checklist").
		Preload("Lists.Cards.Comments").
		Preload("Lists.Cards.Attachments").
		Where("id = ? AND (user_id = ? OR id IN (SELECT board_id FROM board_members WHERE user_id = ?))", id, userID, userID).
		First(&board).Error
	if err != nil {
		return nil, err
	}
	return &board, nil
}

// 3. Simpan Board baru (bisa sekaligus dengan default list/kolom jika ada)
func (r *boardRepository) Create(board *domain.Board) error {
	return r.db.Create(board).Error
}

// 4. Update data Board
func (r *boardRepository) Update(board *domain.Board) error {
	return r.db.Save(board).Error
}

// 5. Hapus Board (Hanya bisa dihapus oleh pemiliknya)
func (r *boardRepository) Delete(id uint, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Board{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("board tidak ditemukan atau bukan milik anda")
	}
	return nil
}
