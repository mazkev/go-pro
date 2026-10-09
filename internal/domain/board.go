package domain

import "time"

// 1. Model Board (Papan Utama)
type Board struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"index"` // Terhubung ke User.ID (Pemilik)
	Title     string    `json:"title"`
	Color     string    `json:"color"`               // e.g. "from-blue-500 to-indigo-500"
	Lists     []List    `json:"lists" gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 2. Model List (Kolom: To Do, In Progress, Done)
type List struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	BoardID   uint      `json:"boardId" gorm:"index"` // Terhubung ke Board.ID
	Title     string    `json:"title"`
	Position  int       `json:"position"`             // Urutan kolom
	Cards     []Card    `json:"cards" gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 3. Model Card (Kartu Tugas)
type Card struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	ListID      uint            `json:"listId" gorm:"index"` // Terhubung ke List.ID
	Content     string          `json:"content"`             // Judul tugas
	Description string          `json:"description"`
	Label       string          `json:"label"`               // e.g. "bg-blue-500"
	CoverURL    string          `json:"coverUrl"`
	DueDate     string          `json:"dueDate"`
	Checked     bool            `json:"checked"`
	Position    int             `json:"position"`            // Urutan kartu di dalam kolom
	Checklist   []ChecklistItem `json:"checklist" gorm:"foreignKey:CardID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// 4. Model ChecklistItem (Sub-tugas di dalam Kartu)
type ChecklistItem struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	CardID  uint   `json:"cardId" gorm:"index"` // Terhubung ke Card.ID
	Title   string `json:"title"`
	Checked bool   `json:"checked"`
}

// 5. Kontrak Repository Board
type BoardRepository interface {
	FindAllByUserID(userID uint) ([]Board, error)
	FindByID(id uint, userID uint) (*Board, error)
	Create(board *Board) error
	Update(board *Board) error
	Delete(id uint, userID uint) error
}

// 6. Kontrak Usecase Board
type BoardUsecase interface {
	GetMyBoards(userID uint) ([]Board, error)
	GetBoardDetail(id uint, userID uint) (*Board, error)
	CreateBoard(userID uint, title string, color string) (*Board, error)
	UpdateBoard(id uint, userID uint, title string, color string) (*Board, error)
	DeleteBoard(id uint, userID uint) error
}

// 7. Kontrak Repository Card
type CardRepository interface {
	FindByID(id uint) (*Card, error)
	Create(card *Card) error
	Update(card *Card) error
	Delete(id uint) error
	Move(cardID uint, targetListID uint, newPosition int) error
}

// 8. Kontrak Usecase Card
type CardUsecase interface {
	CreateCard(listID uint, content string, label string, description string) (*Card, error)
	UpdateCard(id uint, content string, description string, label string, dueDate string, checked bool) (*Card, error)
	MoveCard(id uint, targetListID uint, newPosition int) error
	DeleteCard(id uint) error
}

// 9. Kontrak Repository List
type ListRepository interface {
	FindByID(id uint) (*List, error)
	Create(list *List) error
	Update(list *List) error
	Delete(id uint) error
	Move(listID uint, newPosition int) error
}

// 10. Kontrak Usecase List
type ListUsecase interface {
	CreateList(boardID uint, title string) (*List, error)
	UpdateList(id uint, title string) (*List, error)
	DeleteList(id uint) error
	MoveList(id uint, newPosition int) error
}

// 11. Kontrak Repository Checklist
type ChecklistRepository interface {
	Create(item *ChecklistItem) error
	FindByID(id uint) (*ChecklistItem, error)
	Update(item *ChecklistItem) error
	Delete(id uint) error
}

// 12. Kontrak Usecase Checklist
type ChecklistUsecase interface {
	CreateItem(cardID uint, title string) (*ChecklistItem, error)
	UpdateItem(id uint, title string, checked bool) (*ChecklistItem, error)
	DeleteItem(id uint) error
}
