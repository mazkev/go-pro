package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"
)

type boardUsecase struct {
	repo domain.BoardRepository
}

// Constructor instance usecase
func NewBoardUsecase(repo domain.BoardRepository) domain.BoardUsecase {
	return &boardUsecase{repo: repo}
}

// 1. Ambil semua board milik user yang sedang login
func (u *boardUsecase) GetMyBoards(userID uint) ([]domain.Board, error) {
	return u.repo.FindAllByUserID(userID)
}

// 2. Ambil detail 1 board beserta kolom dan kartu tugasnya
func (u *boardUsecase) GetBoardDetail(id uint, userID uint) (*domain.Board, error) {
	return u.repo.FindByID(id, userID)
}

// 3. Buat Board baru (Otomatis menyertakan 3 kolom Kanban: To Do, In Progress, Done!)
func (u *boardUsecase) CreateBoard(userID uint, title string, color string) (*domain.Board, error) {
	// Aturan Bisnis 1: Judul tidak boleh kosong
	if title == "" {
		return nil, errors.New("judul board wajib diisi")
	}

	// Aturan Bisnis 2: Warna default jika tidak dipilih
	if color == "" {
		color = "from-blue-500 to-indigo-500"
	}

	// Buat Board baru beserta 3 kolom bawaan standar Trello!
	newBoard := &domain.Board{
		UserID: userID,
		Title:  title,
		Color:  color,
		Lists: []domain.List{
			{Title: "To Do", Position: 1},
			{Title: "In Progress", Position: 2},
			{Title: "Done", Position: 3},
		},
	}

	if err := u.repo.Create(newBoard); err != nil {
		return nil, err
	}

	return newBoard, nil
}

// 4. Hapus Board
func (u *boardUsecase) DeleteBoard(id uint, userID uint) error {
	return u.repo.Delete(id, userID)
}
