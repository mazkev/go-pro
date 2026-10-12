package repository

import (
	"ecommerce-api/internal/domain"

	"gorm.io/gorm"
)

type attachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) domain.AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(attachment *domain.Attachment) error {
	return r.db.Create(attachment).Error
}

func (r *attachmentRepository) FindByCardID(cardID uint) ([]domain.Attachment, error) {
	var attachments []domain.Attachment
	err := r.db.Where("card_id = ?", cardID).Order("created_at desc").Find(&attachments).Error
	return attachments, err
}

func (r *attachmentRepository) FindByID(id uint) (*domain.Attachment, error) {
	var attachment domain.Attachment
	err := r.db.First(&attachment, id).Error
	if err != nil {
		return nil, err
	}
	return &attachment, nil
}

func (r *attachmentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Attachment{}, id).Error
}
