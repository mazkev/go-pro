package domain

// 1. Model Data Produk
type Product struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

// 2. Kontrak untuk Repository (Database)
// Siapapun yang menjadi Repository produk, WAJIB punya 5 fungsi ini:
type ProductRepository interface {
	FindAll() ([]Product, error)
	FindByID(id int) (*Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(id int) error
}

// 3. Kontrak untuk Usecase (Logika Bisnis)
// Siapapun yang menjadi Usecase produk, WAJIB punya 5 fungsi ini:
type ProductUsecase interface {
	GetAll() ([]Product, error)
	GetByID(id int) (*Product, error)
	Create(product *Product) error
	Update(id int, product *Product) (*Product, error)
	Delete(id int) error
}
