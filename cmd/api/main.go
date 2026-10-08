package main

import (
	httpDelivery "ecommerce-api/internal/delivery/http"
	"ecommerce-api/internal/delivery/http/middleware"
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

	// Migrasi otomatis seluruh tabel relasional
	db.AutoMigrate(
		&domain.User{},
		&domain.Product{},
		&domain.Board{},
		&domain.List{},
		&domain.Card{},
		&domain.ChecklistItem{},
	)
	return db
}

func main() {
	// 1. Database Layer
	db := initDB()

	// 2. Product Modules
	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := httpDelivery.NewProductHandler(productUsecase)

	// 3. User Modules
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := httpDelivery.NewUserHandler(userUsecase)

	// 4. 📋 Kanban Board Modules (BARU!)
	boardRepo := repository.NewBoardRepository(db)
	boardUsecase := usecase.NewBoardUsecase(boardRepo)
	boardHandler := httpDelivery.NewBoardHandler(boardUsecase)

	// 5. Setup Router Gin
	r := gin.Default()

	// 6. Rute Publik
	r.GET("/products", productHandler.GetAll)
	r.GET("/products/search", productHandler.Search)
	r.GET("/products/:id", productHandler.GetByID)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	// 7. Rute Terkunci (Wajib Token JWT)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Produk
		protected.POST("/products", productHandler.Create)
		protected.PUT("/products/:id", productHandler.Update)
		protected.DELETE("/products/:id", productHandler.Delete)

		// 📋 Rute Kanban Board (Milik User)
		protected.GET("/boards", boardHandler.GetMyBoards)
		protected.GET("/boards/:id", boardHandler.GetByID)
		protected.POST("/boards", boardHandler.Create)
		protected.DELETE("/boards/:id", boardHandler.Delete)
	}

	// 8. Jalankan Server di port 8080
	r.Run(":8080")
}
