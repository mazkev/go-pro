package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"
)

type cardUsecase struct {
	repo domain.CardRepository
}

// Constructor instance usecase kartu
func NewCardUsecase(repo domain.CardRepository) domain.CardUsecase {
	return &cardUsecase{repo: repo}
}

// 1. Buat kartu baru di dalam kolom list tertentu
func (u *cardUsecase) CreateCard(listID uint, content string, label string, description string) (*domain.Card, error) {
	if content == "" {
		return nil, errors.New("konten/judul kartu wajib diisi")
	}

	if label == "" {
		label = "bg-blue-500" // Default label warna Trello
	}

	newCard := &domain.Card{
		ListID:      listID,
		Content:     content,
		Label:       label,
		Description: description,
		Position:    1,
		Checked:     false,
	}

	if err := u.repo.Create(newCard); err != nil {
		return nil, err
	}

	return newCard, nil
}

// 2. Update data kartu
func (u *cardUsecase) UpdateCard(id uint, content string, description string, label string, coverURL string, dueDate string, checked bool) (*domain.Card, error) {
	card, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("kartu tidak ditemukan")
	}

	if content != "" {
		card.Content = content
	}
	card.Description = description
	if label != "" {
		card.Label = label
	}
	card.CoverURL = coverURL
	card.DueDate = dueDate
	card.Checked = checked

	if err := u.repo.Update(card); err != nil {
		return nil, err
	}

	return card, nil
}

// 3. Geser / Drag and Drop Kartu (Pindah kolom dan atur urutan)
func (u *cardUsecase) MoveCard(id uint, targetListID uint, newPosition int) error {
	if targetListID == 0 {
		return errors.New("kolom tujuan tidak valid")
	}
	return u.repo.Move(id, targetListID, newPosition)
}

// 4. Hapus kartu
func (u *cardUsecase) DeleteCard(id uint) error {
	return u.repo.Delete(id)
}

// 5. Tugaskan kartu ke anggota (atau unassign jika nil)
func (u *cardUsecase) AssignCard(cardID uint, assigneeID *uint, assigneeName string) (*domain.Card, error) {
	card, err := u.repo.FindByID(cardID)
	if err != nil || card == nil {
		return nil, errors.New("kartu tidak ditemukan")
	}

	card.AssigneeID = assigneeID
	card.AssigneeName = assigneeName

	if err := u.repo.Update(card); err != nil {
		return nil, err
	}

	return card, nil
}
