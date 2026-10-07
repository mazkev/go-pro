package usecase_test

import (
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/usecase"
	"testing"
)

// 1. MOCK REPOSITORY: Tiruan dari domain.ProductRepository
type mockProductRepo struct {
	domain.ProductRepository // Embed interface agar tidak perlu menulis semua method jika belum dipakai
}

// Kita hanya perlu mock method Create() yang sedang kita uji!
func (m *mockProductRepo) Create(product *domain.Product) error {
	return nil // Simulasi sukses simpan ke DB
}

// 2. UNIT TEST DENGAN TABLE-DRIVEN TESTS
func TestProductUsecase_Create(t *testing.T) {
	mockRepo := &mockProductRepo{}
	uc := usecase.NewProductUsecase(mockRepo)

	// Tabel kasus pengujian
	tests := []struct {
		name       string
		product    *domain.Product
		expectErr  bool
		errMessage string
	}{
		{
			name: "Sukses - Data Produk Valid",
			product: &domain.Product{
				Name:  "Keyboard Mechanical",
				Price: 350000,
				Stock: 10,
			},
			expectErr: false,
		},
		{
			name: "Gagal - Nama Produk Kosong",
			product: &domain.Product{
				Name:  "",
				Price: 350000,
				Stock: 10,
			},
			expectErr:  true,
			errMessage: "nama produk wajib diisi",
		},
		{
			name: "Gagal - Harga 0 atau Minus",
			product: &domain.Product{
				Name:  "Mouse Wireless",
				Price: -5000,
				Stock: 5,
			},
			expectErr:  true,
			errMessage: "harga produk harus lebih besar dari 0",
		},
		{
			name: "Gagal - Stok Minus",
			product: &domain.Product{
				Name:  "Mousepad XL",
				Price: 50000,
				Stock: -2,
			},
			expectErr:  true,
			errMessage: "stok tidak boleh negatif",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := uc.Create(tc.product)

			if tc.expectErr {
				if err == nil {
					t.Errorf("[%s] Mengharapkan error, tetapi hasilnya nil", tc.name)
				} else if err.Error() != tc.errMessage {
					t.Errorf("[%s] Pesan error salah! Dapat: %v, Ingin: %v", tc.name, err.Error(), tc.errMessage)
				}
			} else {
				if err != nil {
					t.Errorf("[%s] Mengharapkan sukses, tetapi dapat error: %v", tc.name, err)
				}
			}
		})
	}
}
