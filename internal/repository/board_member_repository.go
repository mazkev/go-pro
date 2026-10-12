package repository

import (
	"ecommerce-api/internal/domain"

	"gorm.io/gorm"
)

type boardMemberRepository struct {
	db *gorm.DB
}

func NewBoardMemberRepository(db *gorm.DB) domain.BoardMemberRepository {
	return &boardMemberRepository{db: db}
}

func (r *boardMemberRepository) AddMember(member *domain.BoardMember) error {
	return r.db.Create(member).Error
}

func (r *boardMemberRepository) FindByBoardAndUser(boardID uint, userID uint) (*domain.BoardMember, error) {
	var member domain.BoardMember
	err := r.db.Where("board_id = ? AND user_id = ?", boardID, userID).First(&member).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *boardMemberRepository) FindMembersByBoardID(boardID uint) ([]domain.BoardMember, error) {
	var members []domain.BoardMember
	err := r.db.Where("board_id = ?", boardID).Order("created_at asc").Find(&members).Error
	return members, err
}

func (r *boardMemberRepository) UpdateRole(boardID uint, userID uint, role string) error {
	return r.db.Model(&domain.BoardMember{}).
		Where("board_id = ? AND user_id = ?", boardID, userID).
		Update("role", role).Error
}

func (r *boardMemberRepository) RemoveMember(boardID uint, userID uint) error {
	return r.db.Where("board_id = ? AND user_id = ?", boardID, userID).
		Delete(&domain.BoardMember{}).Error
}
