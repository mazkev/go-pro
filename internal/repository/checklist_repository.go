package repository

import (
	"ecommerce-api/internal/domain"

	"gorm.io/gorm"
)

type checklistRepository struct {
	db *gorm.DB
}

func NewChecklistRepository(db *gorm.DB) domain.ChecklistRepository {
	return &checklistRepository{db: db}
}

func (r *checklistRepository) Create(item *domain.ChecklistItem) error {
	return r.db.Create(item).Error
}

func (r *checklistRepository) FindByID(id uint) (*domain.ChecklistItem, error) {
	var item domain.ChecklistItem
	err := r.db.First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *checklistRepository) Update(item *domain.ChecklistItem) error {
	return r.db.Save(item).Error
}

func (r *checklistRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ChecklistItem{}, id).Error
}
