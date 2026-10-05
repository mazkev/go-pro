package repository

import (
	"ecommerce-api/internal/domain"
	"errors"

	"gorm.io/gorm"
)

// Struct ini yang memegang koneksi database GORM
type productRepository struct {
	db *gorm.DB
}

// Fungsi pembuat instance Repository baru
// Perhatikan: dia mengembalikan interface "domain.ProductRepository"
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

// 1. SELECT * FROM products
func (r *productRepository) FindAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Find(&products).Error
	return products, err
}

// 2. SELECT * FROM products WHERE id = ? LIMIT 1
func (r *productRepository) FindByID(id int) (*domain.Product, error) {
	var product domain.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// 3. INSERT INTO products (...) VALUES (...)
func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

// 4. UPDATE products SET ... WHERE id = ?
func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

// 5. DELETE FROM products WHERE id = ?
func (r *productRepository) Delete(id int) error {
	result := r.db.Delete(&domain.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}
// 6. Pencarian: SELECT * FROM products WHERE name LIKE '%keyword%'
func (r *productRepository) Search(keyword string) ([]domain.Product, error) {
	var products []domain.Product
	// Menggunakan parameter '?' untuk keamanan (Mencegah serangan hacker SQL Injection)
	err := r.db.Where("name LIKE ?", "%"+keyword+"%").Find(&products).Error
	return products, err
}
