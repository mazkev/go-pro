package usecase_test

import (
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/usecase"
	"errors"
	"testing"
)

// Mock Repositories
type mockCommentRepo struct {
	domain.CommentRepository
	comments   map[uint]*domain.Comment
	createErr  error
	findByCard []domain.Comment
}

func (m *mockCommentRepo) Create(comment *domain.Comment) error {
	if m.createErr != nil {
		return m.createErr
	}
	comment.ID = uint(len(m.comments) + 1)
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockCommentRepo) FindByID(id uint) (*domain.Comment, error) {
	c, ok := m.comments[id]
	if !ok {
		return nil, errors.New("comment not found")
	}
	return c, nil
}

func (m *mockCommentRepo) FindByCardID(cardID uint) ([]domain.Comment, error) {
	return m.findByCard, nil
}

func (m *mockCommentRepo) Delete(id uint) error {
	delete(m.comments, id)
	return nil
}

type mockCardRepoForComment struct {
	domain.CardRepository
	cards map[uint]*domain.Card
}

func (m *mockCardRepoForComment) FindByID(id uint) (*domain.Card, error) {
	c, ok := m.cards[id]
	if !ok {
		return nil, errors.New("card not found")
	}
	return c, nil
}

type mockUserRepoForComment struct {
	domain.UserRepository
	users map[int]*domain.User
}

func (m *mockUserRepoForComment) FindByID(id int) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func TestCommentUsecase_CreateComment(t *testing.T) {
	cardRepo := &mockCardRepoForComment{
		cards: map[uint]*domain.Card{
			1: {ID: 1, Content: "Card 1"},
		},
	}
	userRepo := &mockUserRepoForComment{
		users: map[int]*domain.User{
			10: {ID: 10, Name: "Kevin Pratama"},
		},
	}

	tests := []struct {
		name       string
		cardID     uint
		userID     uint
		content    string
		expectErr  bool
		errMessage string
	}{
		{
			name:      "Sukses - Tambah Komentar Valid",
			cardID:    1,
			userID:    10,
			content:   "Desain wireframe selesai",
			expectErr: false,
		},
		{
			name:       "Gagal - Konten Kosong",
			cardID:     1,
			userID:     10,
			content:    "   ",
			expectErr:  true,
			errMessage: "isi komentar tidak boleh kosong",
		},
		{
			name:       "Gagal - Kartu Tidak Ditemukan",
			cardID:     999,
			userID:     10,
			content:    "Halo",
			expectErr:  true,
			errMessage: "kartu tidak ditemukan",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			commentRepo := &mockCommentRepo{comments: make(map[uint]*domain.Comment)}
			uc := usecase.NewCommentUsecase(commentRepo, cardRepo, userRepo)

			created, err := uc.CreateComment(tc.cardID, tc.userID, tc.content)

			if tc.expectErr {
				if err == nil {
					t.Errorf("[%s] Mengharapkan error, tapi dapat nil", tc.name)
				} else if err.Error() != tc.errMessage {
					t.Errorf("[%s] Pesan error salah! Dapat: %v, Ingin: %v", tc.name, err.Error(), tc.errMessage)
				}
			} else {
				if err != nil {
					t.Fatalf("[%s] Tidak mengharapkan error, tapi dapat: %v", tc.name, err)
				}
				if created.UserName != "Kevin Pratama" {
					t.Errorf("[%s] Nama user salah: %v", tc.name, created.UserName)
				}
				if created.Content != tc.content {
					t.Errorf("[%s] Konten salah: %v", tc.name, created.Content)
				}
			}
		})
	}
}

func TestCommentUsecase_DeleteComment(t *testing.T) {
	tests := []struct {
		name       string
		commentID  uint
		userID     uint
		expectErr  bool
		errMessage string
	}{
		{
			name:      "Sukses - Hapus Komentar Milik Sendiri",
			commentID: 1,
			userID:    10,
			expectErr: false,
		},
		{
			name:       "Gagal - Hapus Komentar Orang Lain (Unauthorized)",
			commentID:  1,
			userID:     99, // bukan pemilik ID 10
			expectErr:  true,
			errMessage: "anda tidak memiliki izin untuk menghapus komentar ini",
		},
		{
			name:       "Gagal - Komentar Tidak Ada",
			commentID:  404,
			userID:     10,
			expectErr:  true,
			errMessage: "komentar tidak ditemukan",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			commentRepo := &mockCommentRepo{
				comments: map[uint]*domain.Comment{
					1: {ID: 1, CardID: 1, UserID: 10, UserName: "Kevin", Content: "Komentar tes"},
				},
			}
			uc := usecase.NewCommentUsecase(commentRepo, nil, nil)

			err := uc.DeleteComment(tc.commentID, tc.userID)

			if tc.expectErr {
				if err == nil {
					t.Errorf("[%s] Mengharapkan error, tapi dapat nil", tc.name)
				} else if err.Error() != tc.errMessage {
					t.Errorf("[%s] Pesan error salah! Dapat: %v, Ingin: %v", tc.name, err.Error(), tc.errMessage)
				}
			} else {
				if err != nil {
					t.Errorf("[%s] Tidak mengharapkan error, tapi dapat: %v", tc.name, err)
				}
			}
		})
	}
}
