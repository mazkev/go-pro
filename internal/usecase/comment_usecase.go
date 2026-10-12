package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"
	"strings"
)

type commentUsecase struct {
	commentRepo domain.CommentRepository
	cardRepo    domain.CardRepository
	userRepo    domain.UserRepository
}

func NewCommentUsecase(commentRepo domain.CommentRepository, cardRepo domain.CardRepository, userRepo domain.UserRepository) domain.CommentUsecase {
	return &commentUsecase{
		commentRepo: commentRepo,
		cardRepo:    cardRepo,
		userRepo:    userRepo,
	}
}

func (u *commentUsecase) CreateComment(cardID uint, userID uint, content string) (*domain.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, errors.New("isi komentar tidak boleh kosong")
	}

	// 1. Pastikan kartu ada
	card, err := u.cardRepo.FindByID(cardID)
	if err != nil || card == nil {
		return nil, errors.New("kartu tidak ditemukan")
	}

	// 2. Ambil nama pembuat komentar
	userName := "User"
	if u.userRepo != nil {
		user, err := u.userRepo.FindByID(int(userID))
		if err == nil && user != nil && user.Name != "" {
			userName = user.Name
		}
	}

	comment := &domain.Comment{
		CardID:   cardID,
		UserID:   userID,
		UserName: userName,
		Content:  content,
	}

	if err := u.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (u *commentUsecase) GetCommentsByCardID(cardID uint) ([]domain.Comment, error) {
	return u.commentRepo.FindByCardID(cardID)
}

func (u *commentUsecase) DeleteComment(id uint, userID uint) error {
	comment, err := u.commentRepo.FindByID(id)
	if err != nil || comment == nil {
		return errors.New("komentar tidak ditemukan")
	}

	// Hanya pemilik komentar yang dapat menghapusnya
	if comment.UserID != userID {
		return errors.New("anda tidak memiliki izin untuk menghapus komentar ini")
	}

	return u.commentRepo.Delete(id)
}
