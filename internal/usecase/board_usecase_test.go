package usecase_test

import (
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/usecase"
	"errors"
	"testing"
)

type mockBoardRepo struct {
	domain.BoardRepository
	boards map[uint]*domain.Board
}

func (m *mockBoardRepo) Create(board *domain.Board) error {
	board.ID = uint(len(m.boards) + 1)
	m.boards[board.ID] = board
	return nil
}

func (m *mockBoardRepo) FindByID(id uint, userID uint) (*domain.Board, error) {
	b, ok := m.boards[id]
	if !ok || b.UserID != userID {
		return nil, errors.New("board not found")
	}
	return b, nil
}

func (m *mockBoardRepo) FindAllByUserID(userID uint) ([]domain.Board, error) {
	var res []domain.Board
	for _, b := range m.boards {
		if b.UserID == userID {
			res = append(res, *b)
		}
	}
	return res, nil
}

func (m *mockBoardRepo) Update(board *domain.Board) error {
	m.boards[board.ID] = board
	return nil
}

func (m *mockBoardRepo) Delete(id uint, userID uint) error {
	b, ok := m.boards[id]
	if !ok || b.UserID != userID {
		return errors.New("board not found")
	}
	delete(m.boards, id)
	return nil
}

func TestBoardUsecase_CreateBoard(t *testing.T) {
	repo := &mockBoardRepo{boards: make(map[uint]*domain.Board)}
	uc := usecase.NewBoardUsecase(repo)

	// 1. Sukses buat board (otomatis generate 3 kolom Trello standar)
	board, err := uc.CreateBoard(1, "Proyek Sprint Q4", "")
	if err != nil {
		t.Fatalf("Gagal buat board: %v", err)
	}
	if board.Title != "Proyek Sprint Q4" {
		t.Errorf("Judul salah: %v", board.Title)
	}
	if board.Color != "from-blue-500 to-indigo-500" {
		t.Errorf("Warna default salah: %v", board.Color)
	}
	if len(board.Lists) != 3 {
		t.Fatalf("Board harus memiliki 3 kolom bawaan, dapat: %d", len(board.Lists))
	}
	if board.Lists[0].Title != "To Do" || board.Lists[1].Title != "In Progress" || board.Lists[2].Title != "Done" {
		t.Errorf("Kolom bawaan salah: %+v", board.Lists)
	}

	// 2. Gagal jika judul kosong
	_, err = uc.CreateBoard(1, "", "from-emerald-500 to-teal-600")
	if err == nil || err.Error() != "judul board wajib diisi" {
		t.Errorf("Harusnya error 'judul board wajib diisi', dapat: %v", err)
	}
}

func TestBoardUsecase_UpdateBoard(t *testing.T) {
	repo := &mockBoardRepo{
		boards: map[uint]*domain.Board{
			1: {ID: 1, UserID: 10, Title: "Judul Lama", Color: "blue"},
		},
	}
	uc := usecase.NewBoardUsecase(repo)

	// 1. Sukses update judul dan warna
	updated, err := uc.UpdateBoard(1, 10, "Judul Baru", "purple")
	if err != nil {
		t.Fatalf("Gagal update board: %v", err)
	}
	if updated.Title != "Judul Baru" || updated.Color != "purple" {
		t.Errorf("Hasil update tidak sesuai: %+v", updated)
	}

	// 2. Gagal jika user bukan pemilik (akses ditolak)
	_, err = uc.UpdateBoard(1, 999, "Judul Hacking", "black")
	if err == nil || err.Error() != "board tidak ditemukan atau bukan milik anda" {
		t.Errorf("Harusnya ditolak karena bukan pemilik board, dapat: %v", err)
	}
}

func TestBoardUsecase_DeleteBoard(t *testing.T) {
	repo := &mockBoardRepo{
		boards: map[uint]*domain.Board{
			1: {ID: 1, UserID: 10, Title: "Board Hapus"},
		},
	}
	uc := usecase.NewBoardUsecase(repo)

	// Sukses hapus
	err := uc.DeleteBoard(1, 10)
	if err != nil {
		t.Fatalf("Gagal hapus board: %v", err)
	}
	if len(repo.boards) != 0 {
		t.Errorf("Board harusnya sudah kosong")
	}
}
