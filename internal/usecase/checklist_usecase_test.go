package usecase_test

import (
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/usecase"
	"errors"
	"testing"
)

type mockChecklistRepo struct {
	domain.ChecklistRepository
	items map[uint]*domain.ChecklistItem
}

func (m *mockChecklistRepo) Create(item *domain.ChecklistItem) error {
	item.ID = uint(len(m.items) + 1)
	m.items[item.ID] = item
	return nil
}

func (m *mockChecklistRepo) FindByID(id uint) (*domain.ChecklistItem, error) {
	item, ok := m.items[id]
	if !ok {
		return nil, errors.New("item not found")
	}
	return item, nil
}

func (m *mockChecklistRepo) Update(item *domain.ChecklistItem) error {
	m.items[item.ID] = item
	return nil
}

func (m *mockChecklistRepo) Delete(id uint) error {
	delete(m.items, id)
	return nil
}

func TestChecklistUsecase_CreateItem(t *testing.T) {
	repo := &mockChecklistRepo{items: make(map[uint]*domain.ChecklistItem)}
	uc := usecase.NewChecklistUsecase(repo)

	tests := []struct {
		name       string
		cardID     uint
		title      string
		expectErr  bool
		errMessage string
	}{
		{
			name:      "Sukses - Judul Valid",
			cardID:    1,
			title:     "Buat komponen Button",
			expectErr: false,
		},
		{
			name:       "Gagal - Judul Kosong",
			cardID:     1,
			title:      "",
			expectErr:  true,
			errMessage: "judul sub-tugas wajib diisi",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			item, err := uc.CreateItem(tc.cardID, tc.title)
			if tc.expectErr {
				if err == nil {
					t.Errorf("[%s] Mengharapkan error, tapi dapat nil", tc.name)
				} else if err.Error() != tc.errMessage {
					t.Errorf("[%s] Pesan error salah! Dapat: %v, Ingin: %v", tc.name, err.Error(), tc.errMessage)
				}
			} else {
				if err != nil {
					t.Fatalf("[%s] Tidak mengharapkan error: %v", tc.name, err)
				}
				if item.Title != tc.title {
					t.Errorf("[%s] Judul item tidak sesuai", tc.name)
				}
				if item.Checked != false {
					t.Errorf("[%s] Default checked harus false", tc.name)
				}
			}
		})
	}
}

func TestChecklistUsecase_UpdateItem(t *testing.T) {
	repo := &mockChecklistRepo{
		items: map[uint]*domain.ChecklistItem{
			1: {ID: 1, CardID: 1, Title: "Tugas awal", Checked: false},
		},
	}
	uc := usecase.NewChecklistUsecase(repo)

	// 1. Sukses toggle checked
	updated, err := uc.UpdateItem(1, "Tugas diperbarui", true)
	if err != nil {
		t.Fatalf("Gagal update item: %v", err)
	}
	if updated.Title != "Tugas diperbarui" || !updated.Checked {
		t.Errorf("Hasil update tidak sesuai: %+v", updated)
	}

	// 2. Gagal jika item tidak ada
	_, err = uc.UpdateItem(99, "Tidak ada", true)
	if err == nil || err.Error() != "sub-tugas tidak ditemukan" {
		t.Errorf("Harusnya error 'sub-tugas tidak ditemukan', tapi dapat: %v", err)
	}
}

func TestChecklistUsecase_DeleteItem(t *testing.T) {
	repo := &mockChecklistRepo{
		items: map[uint]*domain.ChecklistItem{
			1: {ID: 1, CardID: 1, Title: "Hapus saya", Checked: false},
		},
	}
	uc := usecase.NewChecklistUsecase(repo)

	err := uc.DeleteItem(1)
	if err != nil {
		t.Fatalf("Gagal hapus item: %v", err)
	}
	if _, exists := repo.items[1]; exists {
		t.Errorf("Item harusnya terhapus dari repo")
	}
}
