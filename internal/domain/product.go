package domain

// Product adalah model data produk kita
type Product struct {
	ID    int    `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}
