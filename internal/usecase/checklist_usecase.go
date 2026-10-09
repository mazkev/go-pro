package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"
)

type checklistUsecase struct {
	repo domain.ChecklistRepository
}

func NewChecklistUsecase(repo domain.ChecklistRepository) domain.ChecklistUsecase {
	return &checklistUsecase{repo: repo}
}

func (u *checklistUsecase) CreateItem(cardID uint, title string) (*domain.ChecklistItem, error) {
	if title == "" {
		return nil, errors.New("judul sub-tugas wajib diisi")
	}

	item := &domain.ChecklistItem{
		CardID:  cardID,
		Title:   title,
		Checked: false,
	}

	if err := u.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (u *checklistUsecase) UpdateItem(id uint, title string, checked bool) (*domain.ChecklistItem, error) {
	item, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("sub-tugas tidak ditemukan")
	}

	if title != "" {
		item.Title = title
	}
	item.Checked = checked

	if err := u.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

func (u *checklistUsecase) DeleteItem(id uint) error {
	return u.repo.Delete(id)
}
