package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"
)

type listUsecase struct {
	repo domain.ListRepository
}

func NewListUsecase(repo domain.ListRepository) domain.ListUsecase {
	return &listUsecase{repo: repo}
}

// 1. Logika Tambah Kolom Baru
func (u *listUsecase) CreateList(boardID uint, title string) (*domain.List, error) {
	if title == "" {
		return nil, errors.New("judul kolom wajib diisi")
	}

	newList := &domain.List{
		BoardID:  boardID,
		Title:    title,
		Position: 99, // Posisi default di ujung kanan
		Cards:    []domain.Card{},
	}

	if err := u.repo.Create(newList); err != nil {
		return nil, err
	}
	return newList, nil
}

// 2. Logika Rename Kolom
func (u *listUsecase) UpdateList(id uint, title string) (*domain.List, error) {
	list, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("kolom tidak ditemukan")
	}

	if title == "" {
		return nil, errors.New("judul kolom tidak boleh kosong")
	}

	list.Title = title
	if err := u.repo.Update(list); err != nil {
		return nil, err
	}
	return list, nil
}

// 3. Logika Hapus Kolom
func (u *listUsecase) DeleteList(id uint) error {
	return u.repo.Delete(id)
}

// 4. Logika Geser Urutan Kolom
func (u *listUsecase) MoveList(id uint, newPosition int) error {
	return u.repo.Move(id, newPosition)
}
