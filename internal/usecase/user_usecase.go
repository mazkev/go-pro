package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	repo domain.UserRepository
}

// Fungsi inisialisasi Usecase User
func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: repo}
}

// 1. Logika Registrasi
func (u *userUsecase) Register(req *domain.RegisterRequest) (*domain.User, error) {
	// ATURAN BISNIS 1: Validasi input dasar
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, errors.New("nama, email, dan password wajib diisi")
	}
	if len(req.Password) < 6 {
		return nil, errors.New("password minimal harus 6 karakter")
	}

	// ATURAN BISNIS 2: Cek apakah email sudah terdaftar sebelumnya
	existingUser, _ := u.repo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar, silakan gunakan email lain")
	}

	// KEAMANAN: HASH PASSWORD DENGAN BCRYPT!
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal mengamankan password")
	}

	// Siapkan data user baru
	newUser := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword), // Yang disimpan adalah KODE HASH-nya!
		Role:     "customer",            // Default role pengguna biasa
	}

	// Simpan ke database via repository
	if err := u.repo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// 2. Logika Login
func (u *userUsecase) Login(req *domain.LoginRequest) (*domain.User, error) {
	// 1. Cari user berdasarkan email
	user, err := u.repo.FindByEmail(req.Email)
	if err != nil {
		// Pesan sengaja dibuat samar demi keamanan (agar hacker tidak bisa menebak email mana yang valid)
		return nil, errors.New("email atau password salah")
	}

	// 2. Cocokkan password yang diketik dengan HASH di database menggunakan Bcrypt!
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("email atau password salah")
	}

	// Login berhasil! Kembalikan data user
	return user, nil
}
