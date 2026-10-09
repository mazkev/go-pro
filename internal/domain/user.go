package domain

import "time"

// 1. Model Tabel User di Database
// Perhatikan tag gorm:"unique;not null"
type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-"` // json:"-" agar password HASH tidak ikut keluar di respon JSON
	Role      string    `json:"role"` // "customer" atau "admin"
	CreatedAt time.Time `json:"created_at"`
}

// 2. DTO untuk input registrasi & login dari client
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// 3. Kontrak Repository (Akses Database)
type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id int) (*User, error)
	Update(user *User) error
}

// 4. Kontrak Usecase (Logika Bisnis & Keamanan)
type UserUsecase interface {
	Register(req *RegisterRequest) (*User, error)
	Login(req *LoginRequest) (*AuthResponse, error)
	GetProfile(userID int) (*User, error)
	UpdateProfile(userID int, req *UpdateProfileRequest) (*User, error)
}

// Respon saat login berhasil (membawa Tiket Token & Data User)
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
