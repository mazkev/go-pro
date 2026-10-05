# E-Commerce REST API (Golang)

RESTful API sederhana untuk manajemen produk e-commerce menggunakan bahasa pemrograman Go, framework Gin, dan ORM GORM dengan database SQLite.

## 🚀 Fitur Saat Ini
### Manajemen Produk
- **GET /products**: Mengambil semua daftar produk.
- **GET /products/:id**: Mengambil detail satu produk berdasarkan ID.
- **GET /products/search?keyword=...**: Mencari produk berdasarkan kata kunci nama (SQL LIKE).
- **POST /products**: Menambahkan produk baru ke database.
- **PUT /products/:id**: Memperbarui nama, harga, atau stok produk.
- **DELETE /products/:id**: Menghapus produk dari database.

### Manajemen User & Autentikasi
- **POST /register**: Registrasi user baru dengan password hashing (Bcrypt).
- **POST /login**: Login user dengan verifikasi password hash.

## 🛠️ Tech Stack
- **Language**: Go (Golang) v1.26+
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **Security**: [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) (Password Hashing)
- **ORM**: [GORM](https://gorm.io)
- **Database**: SQLite (Pure Go via [glebarez/sqlite](https://github.com/glebarez/sqlite))
- **Testing**: REST Client (`rest.http`)

## 📁 Struktur Direktori (Clean Architecture)
```
├── cmd/
│   └── api/
│       └── main.go                 # Entry point aplikasi & Dependency Injection (Wiring)
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── product_handler.go  # Controller/Handler HTTP Produk
│   │       └── user_handler.go     # Controller/Handler HTTP User & Auth
│   ├── domain/
│   │   ├── product.go              # Model struct & Interface Produk
│   │   └── user.go                 # Model struct & Interface User
│   ├── repository/
│   │   ├── product_repository.go   # Data Access Layer Produk
│   │   └── user_repository.go      # Data Access Layer User
│   └── usecase/
│       ├── product_usecase.go      # Business Logic Produk
│       └── user_usecase.go         # Business Logic User & Keamanan Bcrypt
├── .gitignore
├── go.mod
├── go.sum
├── README.md
└── rest.http                       # File pengujian endpoint API
```

## 🏃 Cara Menjalankan
1. Clone repositori:
   ```bash
   git clone https://github.com/mazkev/go-pro.git
   cd go-pro
   ```
2. Jalankan aplikasi:
   ```bash
   go run cmd/api/main.go
   ```
3. Server akan aktif di `http://localhost:8080`.
4. Buka file `rest.http` untuk melakukan pengujian endpoint.
