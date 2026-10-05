package main

import (
	httpDelivery "ecommerce-api/internal/delivery/http"
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/repository"
	"ecommerce-api/internal/usecase"

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
	db.AutoMigrate(&domain.Product{},&domain.User{})
	return db
}

func main() {
	// 1. Database Layer (Paling Bawah)
	db := initDB()

	// 2. Repository Layer (Pegang Database)
	productRepo := repository.NewProductRepository(db)

	// 3. Usecase Layer (Pegang Repository)
	productUsecase := usecase.NewProductUsecase(productRepo)

	// 4. Delivery/Handler Layer (Pegang Usecase)
	productHandler := httpDelivery.NewProductHandler(productUsecase)
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := httpDelivery.NewUserHandler(userUsecase)

	// 5. Setup Router Gin
	r := gin.Default()

	// 6. Daftarkan Rute
	r.GET("/products", productHandler.GetAll)
	r.GET("/products/search", productHandler.Search)
	r.GET("/products/:id", productHandler.GetByID)
	r.POST("/products", productHandler.Create)
	r.PUT("/products/:id", productHandler.Update)
	r.DELETE("/products/:id", productHandler.Delete)

	// Rute untuk User
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// 7. Jalankan Server di port 8080
	r.Run(":8080")
}
