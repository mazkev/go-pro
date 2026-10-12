package usecase_test

import (
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/usecase"
	"errors"
	"testing"
)

type mockCardRepo struct {
	domain.CardRepository
	cards   map[uint]*domain.Card
	movedTo struct {
		targetListID uint
		newPosition  int
	}
}

func (m *mockCardRepo) Create(card *domain.Card) error {
	card.ID = uint(len(m.cards) + 1)
	m.cards[card.ID] = card
	return nil
}

func (m *mockCardRepo) FindByID(id uint) (*domain.Card, error) {
	c, ok := m.cards[id]
	if !ok {
		return nil, errors.New("card not found")
	}
	return c, nil
}

func (m *mockCardRepo) Update(card *domain.Card) error {
	m.cards[card.ID] = card
	return nil
}

func (m *mockCardRepo) Move(id uint, targetListID uint, newPosition int) error {
	c, ok := m.cards[id]
	if !ok {
		return errors.New("card not found")
	}
	c.ListID = targetListID
	c.Position = newPosition
	m.movedTo.targetListID = targetListID
	m.movedTo.newPosition = newPosition
	return nil
}

func (m *mockCardRepo) Delete(id uint) error {
	delete(m.cards, id)
	return nil
}

func TestCardUsecase_CreateCard(t *testing.T) {
	repo := &mockCardRepo{cards: make(map[uint]*domain.Card)}
	uc := usecase.NewCardUsecase(repo)

	// 1. Sukses buat kartu
	card, err := uc.CreateCard(1, "Buat Navbar Baru", "", "Deskripsi detail")
	if err != nil {
		t.Fatalf("Gagal buat kartu: %v", err)
	}
	if card.Content != "Buat Navbar Baru" {
		t.Errorf("Konten kartu salah: %v", card.Content)
	}
	if card.Label != "bg-blue-500" { // Default label
		t.Errorf("Label default harus bg-blue-500, dapat: %v", card.Label)
	}

	// 2. Gagal jika konten kosong
	_, err = uc.CreateCard(1, "", "bg-red-500", "")
	if err == nil || err.Error() != "konten/judul kartu wajib diisi" {
		t.Errorf("Harusnya error 'konten/judul kartu wajib diisi', dapat: %v", err)
	}
}

func TestCardUsecase_MoveCard(t *testing.T) {
	repo := &mockCardRepo{
		cards: map[uint]*domain.Card{
			1: {ID: 1, ListID: 1, Position: 1, Content: "Pindahkan saya"},
		},
	}
	uc := usecase.NewCardUsecase(repo)

	// 1. Sukses geser kartu
	err := uc.MoveCard(1, 2, 3)
	if err != nil {
		t.Fatalf("Gagal memindahkan kartu: %v", err)
	}
	if repo.movedTo.targetListID != 2 || repo.movedTo.newPosition != 3 {
		t.Errorf("Data perpindahan kartu salah: %+v", repo.movedTo)
	}

	// 2. Gagal jika kolom tujuan tidak valid (0)
	err = uc.MoveCard(1, 0, 1)
	if err == nil || err.Error() != "kolom tujuan tidak valid" {
		t.Errorf("Harusnya error 'kolom tujuan tidak valid', dapat: %v", err)
	}
}

func TestCardUsecase_UpdateCard(t *testing.T) {
	repo := &mockCardRepo{
		cards: map[uint]*domain.Card{
			1: {ID: 1, Content: "Judul lama", Label: "bg-blue-500"},
		},
	}
	uc := usecase.NewCardUsecase(repo)

	// Sukses update kartu
	updated, err := uc.UpdateCard(1, "Judul baru", "Deskripsi baru", "bg-emerald-500", "http://image.png", "2026-10-30", true)
	if err != nil {
		t.Fatalf("Gagal update kartu: %v", err)
	}
	if updated.Content != "Judul baru" || updated.Label != "bg-emerald-500" || !updated.Checked {
		t.Errorf("Hasil update tidak sesuai: %+v", updated)
	}

	// Gagal jika kartu tidak ditemukan
	_, err = uc.UpdateCard(999, "Judul", "", "", "", "", false)
	if err == nil || err.Error() != "kartu tidak ditemukan" {
		t.Errorf("Harusnya error 'kartu tidak ditemukan', dapat: %v", err)
	}
}
