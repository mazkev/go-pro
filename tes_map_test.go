package main

import "testing"

// MOCK REPOSITORY (Gudang Tiruan / Palsu untuk Test)
type MockOrderRepository struct {
	simulasiGagal bool
}

// Implementasikan method sesuai kontrak Interface!
func (m *MockOrderRepository) SimpanPesanan(id int, total int) error {
	if m.simulasiGagal {
		return nil
	}
	return nil // Pura-puranya sukses tersimpan
}

func TestCheckout(t *testing.T) {
	// Skenario 1: Checkout Sukses dengan Mock Repo
	mockRepoSukses := &MockOrderRepository{simulasiGagal: false}
	usecaseSukses := &OrderUsecase{repo: mockRepoSukses}

	pesan, err := usecaseSukses.Checkout(101, 50000)
	if err != nil || pesan != "Pesanan berhasil dibuat!" {
		t.Errorf("Harusnya sukses, tapi gagal: %v", err)
	}

	// Skenario 2: Aturan Bisnis Menolak Total <= 0 (Tanpa manggil repo)
	pesan2, err2 := usecaseSukses.Checkout(102, -5000)
	if err2 == nil {
		t.Errorf("Harusnya ditolak karena total minus, tapi malah lolos: %s", pesan2)
	}
}
