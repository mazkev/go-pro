package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"
)

// Struct ini yang memegang Repository
type productUsecase struct {
	repo domain.ProductRepository
}

// Fungsi pembuat instance Usecase baru
// Perhatikan: dia mengembalikan interface "domain.ProductUsecase"
func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{repo: repo}
}

// 1. Ambil semua produk
func (u *productUsecase) GetAll() ([]domain.Product, error) {
	return u.repo.FindAll()
}

// 2. Ambil 1 produk by ID
func (u *productUsecase) GetByID(id int) (*domain.Product, error) {
	return u.repo.FindByID(id)
}

// 3. Tambah produk baru (dengan aturan bisnis!)
func (u *productUsecase) Create(product *domain.Product) error {
	// ATURAN BISNIS 1: Nama tidak boleh kosong
	if product.Name == "" {
		return errors.New("nama produk wajib diisi")
	}

	// ATURAN BISNIS 2: Harga tidak boleh 0 atau minus
	if product.Price <= 0 {
		return errors.New("harga produk harus lebih besar dari 0")
	}

	// ATURAN BISNIS 3: Stok tidak boleh minus
	if product.Stock < 0 {
		return errors.New("stok tidak boleh negatif")
	}

	// Jika semua aturan lolos, suruh repository simpan ke database
	return u.repo.Create(product)
}

// 4. Update produk (dengan aturan bisnis!)
func (u *productUsecase) Update(id int, input *domain.Product) (*domain.Product, error) {
	// Cari dulu barangnya lewat repository
	existingProduct, err := u.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	// Terapkan data baru
	if input.Name != "" {
		existingProduct.Name = input.Name
	}
	if input.Price < 0 {
		return nil, errors.New("harga produk tidak boleh minus")
	}
	if input.Price > 0 {
		existingProduct.Price = input.Price
	}
	// 3. Validasi & Update Stok
	if input.Stock < 0 {
		return nil, errors.New("stok produk tidak boleh minus")
	}
	if input.Stock >= 0 {
		existingProduct.Stock = input.Stock
	}

	// Suruh repository simpan perubahan
	if err := u.repo.Update(existingProduct); err != nil {
		return nil, err
	}

	return existingProduct, nil
}

// 5. Hapus produk
func (u *productUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
