package main

import "errors"

// 1. KONTRAK INTERFACE
type OrderRepositoryInterface interface {
	SimpanPesanan(id int, total int) error
}

// 2. USECASE (Logika Bisnis)
type OrderUsecase struct {
	repo OrderRepositoryInterface // Bergantung pada Interface, bukan DB nyata!
}

func (u *OrderUsecase) Checkout(id int, total int) (string, error) {
	// Aturan Bisnis: Total belanja tidak boleh 0 atau negatif
	if total <= 0 {
		return "", errors.New("total belanja tidak valid")
	}

	// Simpan ke repository
	err := u.repo.SimpanPesanan(id, total)
	if err != nil {
		return "", errors.New("gagal menyimpan pesanan")
	}

	return "Pesanan berhasil dibuat!", nil
}
