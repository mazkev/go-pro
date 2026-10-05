package repository

import (
	"ecommerce-api/internal/domain"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// Fungsi inisialisasi Repository User
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// 1. Simpan user baru ke database: INSERT INTO users ...
func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// 2. Cari user berdasarkan Email: SELECT * FROM users WHERE email = ? LIMIT 1
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// 3. Cari user berdasarkan ID: SELECT * FROM users WHERE id = ? LIMIT 1
func (r *userRepository) FindByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
