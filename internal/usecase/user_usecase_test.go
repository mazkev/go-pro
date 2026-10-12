package usecase_test

import (
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/usecase"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	domain.UserRepository
	usersByEmail map[string]*domain.User
	usersByID    map[int]*domain.User
}

func (m *mockUserRepo) Create(user *domain.User) error {
	user.ID = len(m.usersByID) + 1
	m.usersByID[user.ID] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func (m *mockUserRepo) FindByEmail(email string) (*domain.User, error) {
	u, ok := m.usersByEmail[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByID(id int) (*domain.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepo) Update(user *domain.User) error {
	m.usersByID[user.ID] = user
	m.usersByEmail[user.Email] = user
	return nil
}

func TestUserUsecase_Register(t *testing.T) {
	repo := &mockUserRepo{
		usersByEmail: make(map[string]*domain.User),
		usersByID:    make(map[int]*domain.User),
	}
	uc := usecase.NewUserUsecase(repo)

	// 1. Sukses registrasi
	req := &domain.RegisterRequest{
		Name:     "Kevin",
		Email:    "kevin@example.com",
		Password: "password123",
	}
	user, err := uc.Register(req)
	if err != nil {
		t.Fatalf("Gagal registrasi: %v", err)
	}
	if user.Email != "kevin@example.com" {
		t.Errorf("Email salah: %v", user.Email)
	}
	// Password harus di-hash bcrypt, bukan plaintext
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123")); err != nil {
		t.Errorf("Password harus di-hash bcrypt dengan benar")
	}

	// 2. Gagal: Email duplikat
	_, err = uc.Register(req)
	if err == nil || err.Error() != "email sudah terdaftar, silakan gunakan email lain" {
		t.Errorf("Harusnya error email sudah terdaftar, dapat: %v", err)
	}

	// 3. Gagal: Password kurang dari 6 karakter
	shortReq := &domain.RegisterRequest{
		Name:     "Budi",
		Email:    "budi@example.com",
		Password: "123",
	}
	_, err = uc.Register(shortReq)
	if err == nil || err.Error() != "password minimal harus 6 karakter" {
		t.Errorf("Harusnya error password minimal 6 karakter, dapat: %v", err)
	}
}

func TestUserUsecase_UpdateProfile(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("rahasia123"), bcrypt.DefaultCost)
	repo := &mockUserRepo{
		usersByEmail: make(map[string]*domain.User),
		usersByID: map[int]*domain.User{
			1: {ID: 1, Name: "Kevin Lama", Email: "kevin@example.com", Password: string(hashedPassword)},
		},
	}
	uc := usecase.NewUserUsecase(repo)

	// Sukses update nama
	updated, err := uc.UpdateProfile(1, &domain.UpdateProfileRequest{Name: "Kevin Pratama Baru"})
	if err != nil {
		t.Fatalf("Gagal update profil: %v", err)
	}
	if updated.Name != "Kevin Pratama Baru" {
		t.Errorf("Nama tidak terupdate: %v", updated.Name)
	}

	// Gagal: User tidak ditemukan
	_, err = uc.UpdateProfile(999, &domain.UpdateProfileRequest{Name: "Nobody"})
	if err == nil || err.Error() != "user tidak ditemukan" {
		t.Errorf("Harusnya error user tidak ditemukan, dapat: %v", err)
	}
}
