package repository

import (
	"ecommerce-api/internal/domain"

	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) domain.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByCardID(cardID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.Where("card_id = ?", cardID).Order("created_at asc").Find(&comments).Error
	return comments, err
}

func (r *commentRepository) FindByID(id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Comment{}, id).Error
}
