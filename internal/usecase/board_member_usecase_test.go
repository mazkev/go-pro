package usecase_test

import (
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/usecase"
	"errors"
	"testing"
)

type mockBoardMemberRepo struct {
	domain.BoardMemberRepository
	members map[string]*domain.BoardMember // key: "boardID-userID"
}

func (m *mockBoardMemberRepo) AddMember(member *domain.BoardMember) error {
	key := memberKey(member.BoardID, member.UserID)
	m.members[key] = member
	return nil
}

func (m *mockBoardMemberRepo) FindByBoardAndUser(boardID uint, userID uint) (*domain.BoardMember, error) {
	key := memberKey(boardID, userID)
	mem, ok := m.members[key]
	if !ok {
		return nil, errors.New("member not found")
	}
	return mem, nil
}

func (m *mockBoardMemberRepo) FindMembersByBoardID(boardID uint) ([]domain.BoardMember, error) {
	var list []domain.BoardMember
	for _, mem := range m.members {
		if mem.BoardID == boardID {
			list = append(list, *mem)
		}
	}
	return list, nil
}

func (m *mockBoardMemberRepo) UpdateRole(boardID uint, userID uint, role string) error {
	key := memberKey(boardID, userID)
	mem, ok := m.members[key]
	if !ok {
		return errors.New("member not found")
	}
	mem.Role = role
	return nil
}

func (m *mockBoardMemberRepo) RemoveMember(boardID uint, userID uint) error {
	key := memberKey(boardID, userID)
	delete(m.members, key)
	return nil
}

func memberKey(boardID, userID uint) string {
	return string(rune(boardID)) + "-" + string(rune(userID))
}

type mockBoardRepoForMember struct {
	domain.BoardRepository
	boards map[uint]*domain.Board
}

func (m *mockBoardRepoForMember) FindByID(id uint, userID uint) (*domain.Board, error) {
	b, ok := m.boards[id]
	if !ok {
		return nil, errors.New("board not found")
	}
	return b, nil
}

type mockUserRepoForMember struct {
	domain.UserRepository
	usersByEmail map[string]*domain.User
	usersByID    map[int]*domain.User
}

func (m *mockUserRepoForMember) FindByEmail(email string) (*domain.User, error) {
	u, ok := m.usersByEmail[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (m *mockUserRepoForMember) FindByID(id int) (*domain.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func TestBoardMemberUsecase_InviteMember(t *testing.T) {
	boardRepo := &mockBoardRepoForMember{
		boards: map[uint]*domain.Board{
			1: {ID: 1, UserID: 10, Title: "Sprint Board"}, // Owner ID = 10
		},
	}
	userRepo := &mockUserRepoForMember{
		usersByEmail: map[string]*domain.User{
			"budi@example.com": {ID: 20, Name: "Budi", Email: "budi@example.com"},
		},
	}

	tests := []struct {
		name          string
		currentUserID uint
		email         string
		role          string
		expectErr     bool
		errMessage    string
	}{
		{
			name:          "Sukses - Owner mengundang member baru",
			currentUserID: 10,
			email:         "budi@example.com",
			role:          "member",
			expectErr:     false,
		},
		{
			name:          "Gagal - Non-admin mencoba mengundang",
			currentUserID: 99, // bukan admin/owner
			email:         "budi@example.com",
			role:          "viewer",
			expectErr:     true,
			errMessage:    "hanya admin yang memiliki izin untuk mengundang anggota",
		},
		{
			name:          "Gagal - Email belum terdaftar",
			currentUserID: 10,
			email:         "nobody@example.com",
			role:          "viewer",
			expectErr:     true,
			errMessage:    "pengguna dengan email tersebut belum terdaftar di aplikasi",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			memberRepo := &mockBoardMemberRepo{members: make(map[string]*domain.BoardMember)}
			uc := usecase.NewBoardMemberUsecase(memberRepo, boardRepo, userRepo)

			created, err := uc.InviteMember(1, tc.currentUserID, tc.email, tc.role)

			if tc.expectErr {
				if err == nil {
					t.Errorf("[%s] Mengharapkan error tapi dapat nil", tc.name)
				} else if err.Error() != tc.errMessage {
					t.Errorf("[%s] Pesan error salah! Dapat: %v, Ingin: %v", tc.name, err.Error(), tc.errMessage)
				}
			} else {
				if err != nil {
					t.Fatalf("[%s] Tidak mengharapkan error: %v", tc.name, err)
				}
				if created.Role != "member" {
					t.Errorf("[%s] Role salah: %v", tc.name, created.Role)
				}
			}
		})
	}
}

func TestBoardMemberUsecase_UpdateMemberRole(t *testing.T) {
	boardRepo := &mockBoardRepoForMember{
		boards: map[uint]*domain.Board{
			1: {ID: 1, UserID: 10, Title: "Sprint Board"}, // Owner = 10
		},
	}
	memberRepo := &mockBoardMemberRepo{
		members: map[string]*domain.BoardMember{
			memberKey(1, 20): {BoardID: 1, UserID: 20, Role: "member"},
		},
	}
	uc := usecase.NewBoardMemberUsecase(memberRepo, boardRepo, nil)

	// 1. Sukses update role oleh Admin
	err := uc.UpdateMemberRole(1, 10, 20, "admin")
	if err != nil {
		t.Fatalf("Gagal update role: %v", err)
	}

	// 2. Gagal: Non-admin mencoba ubah role
	err = uc.UpdateMemberRole(1, 99, 20, "viewer")
	if err == nil || err.Error() != "hanya admin yang dapat mengubah role anggota" {
		t.Errorf("Harusnya ditolak karena bukan admin, dapat: %v", err)
	}
}
