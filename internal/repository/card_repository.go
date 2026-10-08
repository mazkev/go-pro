package repository

import (
	"ecommerce-api/internal/domain"

	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}

// Constructor pembuat instance repository
func NewCardRepository(db *gorm.DB) domain.CardRepository {
	return &cardRepository{db: db}
}

// 1. Ambil 1 kartu detail beserta checklist sub-tugasnya
func (r *cardRepository) FindByID(id uint) (*domain.Card, error) {
	var card domain.Card
	err := r.db.Preload("Checklist").First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// 2. Simpan kartu baru
func (r *cardRepository) Create(card *domain.Card) error {
	return r.db.Create(card).Error
}

// 3. Update data kartu (judul, deskripsi, tanggal, dll)
func (r *cardRepository) Update(card *domain.Card) error {
	return r.db.Save(card).Error
}

// 4. Hapus kartu
func (r *cardRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Card{}, id).Error
}

// 5. Geser Kartu / Drag and Drop (Pindah kolom dan/atau ganti urutan)
func (r *cardRepository) Move(cardID uint, targetListID uint, newPosition int) error {
	return r.db.Model(&domain.Card{}).Where("id = ?", cardID).Updates(map[string]interface{}{
		"list_id":  targetListID,
		"position": newPosition,
	}).Error
}
