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

// 1. Ambil semua board milik user tertentu (beserta seluruh list, kartu, dan checklist-nya)
func (r *boardRepository) FindAllByUserID(userID uint) ([]domain.Board, error) {
	var boards []domain.Board
	// Nested Preload: Otomatis load Lists -> Cards -> Checklist!
	err := r.db.Preload("Lists.Cards.Checklist").Where("user_id = ?", userID).Find(&boards).Error
	return boards, err
}

// 2. Ambil 1 board detail berdasarkan ID dan pemiliknya
func (r *boardRepository) FindByID(id uint, userID uint) (*domain.Board, error) {
	var board domain.Board
	err := r.db.Preload("Lists.Cards.Checklist").Where("id = ? AND user_id = ?", id, userID).First(&board).Error
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
