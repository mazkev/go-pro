package repository

import (
	"ecommerce-api/internal/domain"

	"gorm.io/gorm"
)

type listRepository struct {
	db *gorm.DB
}

func NewListRepository(db *gorm.DB) domain.ListRepository {
	return &listRepository{db: db}
}

// 1. Ambil list beserta kartu dan checklist di dalamnya
func (r *listRepository) FindByID(id uint) (*domain.List, error) {
	var list domain.List
	err := r.db.Preload("Cards.Checklist").First(&list, id).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

// 2. Simpan list baru
func (r *listRepository) Create(list *domain.List) error {
	return r.db.Create(list).Error
}

// 3. Update data list (misal: rename)
func (r *listRepository) Update(list *domain.List) error {
	return r.db.Save(list).Error
}

// 4. Hapus list (GORM Cascade akan otomatis hapus kartu-kartu di dalamnya)
func (r *listRepository) Delete(id uint) error {
	return r.db.Delete(&domain.List{}, id).Error
}

// 5. Geser posisi list (Drag and Drop kolom)
func (r *listRepository) Move(listID uint, newPosition int) error {
	return r.db.Model(&domain.List{}).Where("id = ?", listID).Update("position", newPosition).Error
}
