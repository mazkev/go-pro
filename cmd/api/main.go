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
		&domain.Comment{},
		&domain.Attachment{},
		&domain.BoardMember{},
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

	
	// 5. 🗂️ Kanban Card Modules (BARU!)
	cardRepo := repository.NewCardRepository(db)
	cardUsecase := usecase.NewCardUsecase(cardRepo)
	cardHandler := httpDelivery.NewCardHandler(cardUsecase)

	// 5. Setup Router Gin
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	// 5. 📑 Kanban List Modules (BARU!)
	listRepo := repository.NewListRepository(db)
	listUsecase := usecase.NewListUsecase(listRepo)
	listHandler := httpDelivery.NewListHandler(listUsecase)

	// 6. ✅ Kanban Checklist Modules
	checklistRepo := repository.NewChecklistRepository(db)
	checklistUsecase := usecase.NewChecklistUsecase(checklistRepo)
	checklistHandler := httpDelivery.NewChecklistHandler(checklistUsecase)

	// 7. 💬 Kanban Card Comment Modules
	commentRepo := repository.NewCommentRepository(db)
	commentUsecase := usecase.NewCommentUsecase(commentRepo, cardRepo, userRepo)
	commentHandler := httpDelivery.NewCommentHandler(commentUsecase)

	// 8. 📎 Kanban Card Attachment Modules (BARU!)
	attachmentRepo := repository.NewAttachmentRepository(db)
	attachmentUsecase := usecase.NewAttachmentUsecase(attachmentRepo, cardRepo)
	attachmentHandler := httpDelivery.NewAttachmentHandler(attachmentUsecase)

	// 9. 👥 Kanban Board Members & RBAC Modules (BARU!)
	memberRepo := repository.NewBoardMemberRepository(db)
	memberUsecase := usecase.NewBoardMemberUsecase(memberRepo, boardRepo, userRepo)
	memberHandler := httpDelivery.NewBoardMemberHandler(memberUsecase, cardUsecase)

	// 6. Rute Publik
	r.GET("/products", productHandler.GetAll)
	r.GET("/products/search", productHandler.Search)
	r.GET("/products/:id", productHandler.GetByID)
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)
	r.Static("/uploads", "./uploads")

	// 7. Rute Terkunci (Wajib Token JWT)
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{	
		protected.GET("/me", userHandler.GetMe)
		protected.PUT("/me", userHandler.UpdateProfile)
		protected.GET("/users/me", userHandler.GetMe)
		protected.PUT("/users/profile", userHandler.UpdateProfile)
		// Produk
		protected.POST("/products", productHandler.Create)
		protected.PUT("/products/:id", productHandler.Update)
		protected.DELETE("/products/:id", productHandler.Delete)

		// 📋 Rute Kanban Board (Milik User)
		protected.GET("/boards", boardHandler.GetMyBoards)
		protected.GET("/boards/:id", boardHandler.GetByID)
		protected.POST("/boards", boardHandler.Create)
		protected.PUT("/boards/:id", boardHandler.Update)
		protected.DELETE("/boards/:id", boardHandler.Delete)

		// 📑 Rute Kolom/List (BARU!)
		protected.POST("/boards/:boardId/lists", listHandler.Create)
		protected.PUT("/lists/:id", listHandler.Update)
		protected.PUT("/lists/:id/move", listHandler.Move)
		protected.DELETE("/lists/:id", listHandler.Delete)

		// 🗂️ Rute Kartu Tugas & Drag and Drop (BARU!)
		protected.POST("/lists/:listId/cards", cardHandler.Create)
		protected.PUT("/cards/:id", cardHandler.Update)
		protected.PUT("/cards/:id/move", cardHandler.Move)
		protected.DELETE("/cards/:id", cardHandler.Delete)

		// ✅ Rute Sub-tugas Checklist
		protected.POST("/cards/:cardId/checklist", checklistHandler.Create)
		protected.PUT("/checklist/:id", checklistHandler.Update)
		protected.DELETE("/checklist/:id", checklistHandler.Delete)

		// 💬 Rute Komentar & Aktivitas Kartu
		protected.POST("/cards/:cardId/comments", commentHandler.Create)
		protected.GET("/cards/:cardId/comments", commentHandler.GetByCardID)
		protected.DELETE("/comments/:id", commentHandler.Delete)

		// 📎 Rute Lampiran File Kartu (BARU!)
		protected.POST("/cards/:cardId/attachments", attachmentHandler.Upload)
		protected.GET("/cards/:cardId/attachments", attachmentHandler.GetByCardID)
		protected.DELETE("/attachments/:id", attachmentHandler.Delete)

		// 👥 Rute Anggota Board & Penugasan Kartu (RBAC - BARU!)
		protected.POST("/boards/:boardId/members", memberHandler.Invite)
		protected.GET("/boards/:boardId/members", memberHandler.GetMembers)
		protected.PUT("/boards/:boardId/members/:userId", memberHandler.UpdateRole)
		protected.DELETE("/boards/:boardId/members/:userId", memberHandler.RemoveMember)
		protected.PUT("/cards/:cardId/assign", memberHandler.AssignCard)
	}
	

	// 8. Jalankan Server di port 8080
	r.Run(":8080")
}
