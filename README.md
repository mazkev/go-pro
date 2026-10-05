# E-Commerce REST API (Golang)

RESTful API sederhana untuk manajemen produk e-commerce menggunakan bahasa pemrograman Go, framework Gin, dan ORM GORM dengan database SQLite.

## 🚀 Fitur Saat Ini
- **GET /products**: Mengambil semua daftar produk.
- **GET /products/:id**: Mengambil detail satu produk berdasarkan ID.
- **POST /products**: Menambahkan produk baru ke database.
- **PUT /products/:id**: Memperbarui nama, harga, atau stok produk.
- **DELETE /products/:id**: Menghapus produk dari database.

## 🛠️ Tech Stack
- **Language**: Go (Golang) v1.26+
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io)
- **Database**: SQLite (Pure Go via [glebarez/sqlite](https://github.com/glebarez/sqlite))
- **Testing**: REST Client (`rest.http`)

## 📁 Struktur Direktori
```
├── cmd/
│   └── api/
│       └── main.go           # Entry point aplikasi & registrasi router
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       └── product_handle.go # Controller/Handler HTTP untuk Produk
│   └── domain/
│       └── product.go        # Model struct Produk
├── .gitignore
├── go.mod
├── go.sum
├── README.md
└── rest.http                 # File pengujian endpoint API
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
