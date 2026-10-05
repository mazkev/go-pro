package main

import (
	httpDelivery "ecommerce-api/internal/delivery/http"
	"ecommerce-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	if err != nil {
		panic("Gagal menghubungkan ke database: " + err.Error())
	}

	// Migrasi tabel menggunakan struct dari domain
	db.AutoMigrate(&domain.Product{})
	return db
}

func main() {
	// 1. Nyalakan Database
	db := initDB()

	// 2. Buat Handler Produk (dengan menyuntikkan koneksi DB ke dalamnya)
	productHandler := httpDelivery.NewProductHandler(db)

	// 3. Router Gin
	r := gin.Default()

	// 4. Daftarkan Rute (Sangat rapi & mudah dibaca!)
	r.GET("/products", productHandler.GetAll)
	r.GET("/products/:id", productHandler.GetByID)
	r.POST("/products", productHandler.Create)
	r.PUT("/products/:id", productHandler.Update)
	r.DELETE("/products/:id", productHandler.Delete)

	// 5. Buka Server di port 8080
	r.Run(":8080")
}
