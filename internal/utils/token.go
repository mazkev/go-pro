package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Kunci rahasia untuk menandatangani token (Di dunia nyata disimpan di .env)
var SecretKey = []byte("kunci_rahasia_super_aman_12345")

// Data yang dititipkan di dalam token
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// 1. Fungsi Membuat Token Baru (Saat Login Berhasil)
func GenerateToken(userID int, email string, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			// Token berlaku selama 24 jam sejak dibuat
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Bungkus payload dengan algoritma HS256 dan kunci rahasia
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// 2. Fungsi Memeriksa Keaslian Token
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token tidak valid atau kadaluarsa")
	}

	return claims, nil
}
