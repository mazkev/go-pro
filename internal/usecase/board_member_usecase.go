package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"
	"strings"
)

type boardMemberUsecase struct {
	memberRepo domain.BoardMemberRepository
	boardRepo  domain.BoardRepository
	userRepo   domain.UserRepository
}

func NewBoardMemberUsecase(
	memberRepo domain.BoardMemberRepository,
	boardRepo domain.BoardRepository,
	userRepo domain.UserRepository,
) domain.BoardMemberUsecase {
	return &boardMemberUsecase{
		memberRepo: memberRepo,
		boardRepo:  boardRepo,
		userRepo:   userRepo,
	}
}

// 1. Mengundang anggota baru berdasarkan email dan menentukan role-nya
func (u *boardMemberUsecase) InviteMember(boardID uint, currentUserID uint, email string, role string) (*domain.BoardMember, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return nil, errors.New("email anggota wajib diisi")
	}

	role = strings.ToLower(strings.TrimSpace(role))
	if role != "admin" && role != "member" && role != "viewer" {
		role = "member" // default role
	}

	// Cek apakah board ada dan user saat ini punya akses
	board, err := u.boardRepo.FindByID(boardID, currentUserID)
	if err != nil || board == nil {
		return nil, errors.New("papan tidak ditemukan atau anda tidak memiliki akses")
	}

	// Proteksi RBAC: Hanya Admin atau Pemilik board yang boleh mengundang anggota baru!
	currentUserRole, _ := u.GetUserRole(boardID, currentUserID)
	if currentUserRole != "admin" && board.UserID != currentUserID {
		return nil, errors.New("hanya admin yang memiliki izin untuk mengundang anggota")
	}

	// Cari user yang dituju berdasarkan email
	targetUser, err := u.userRepo.FindByEmail(email)
	if err != nil || targetUser == nil {
		return nil, errors.New("pengguna dengan email tersebut belum terdaftar di aplikasi")
	}

	// Cegah mengundang pemilik utama papan
	if uint(targetUser.ID) == board.UserID {
		return nil, errors.New("pengguna tersebut adalah pemilik utama papan ini")
	}

	// Cek apakah user sudah terdaftar di board_members
	existing, _ := u.memberRepo.FindByBoardAndUser(boardID, uint(targetUser.ID))
	if existing != nil {
		return nil, errors.New("pengguna tersebut sudah menjadi anggota di papan ini")
	}

	newMember := &domain.BoardMember{
		BoardID:   boardID,
		UserID:    uint(targetUser.ID),
		UserName:  targetUser.Name,
		UserEmail: targetUser.Email,
		Role:      role,
	}

	if err := u.memberRepo.AddMember(newMember); err != nil {
		return nil, err
	}

	return newMember, nil
}

// 2. Ambil semua anggota board
func (u *boardMemberUsecase) GetBoardMembers(boardID uint, currentUserID uint) ([]domain.BoardMember, error) {
	board, err := u.boardRepo.FindByID(boardID, currentUserID)
	if err != nil || board == nil {
		return nil, errors.New("papan tidak ditemukan atau anda tidak memiliki akses")
	}

	members, err := u.memberRepo.FindMembersByBoardID(boardID)
	if err != nil {
		return nil, err
	}

	// Tambahkan pemilik utama ke daftar anggota paling atas (sebagai Owner/Admin)
	ownerUser, _ := u.userRepo.FindByID(int(board.UserID))
	ownerName := "Owner"
	ownerEmail := ""
	if ownerUser != nil {
		ownerName = ownerUser.Name
		ownerEmail = ownerUser.Email
	}

	allMembers := []domain.BoardMember{
		{
			ID:        0,
			BoardID:   boardID,
			UserID:    board.UserID,
			UserName:  ownerName,
			UserEmail: ownerEmail,
			Role:      "admin",
		},
	}
	allMembers = append(allMembers, members...)

	return allMembers, nil
}

// 3. Mengubah role anggota
func (u *boardMemberUsecase) UpdateMemberRole(boardID uint, currentUserID uint, targetUserID uint, newRole string) error {
	newRole = strings.ToLower(strings.TrimSpace(newRole))
	if newRole != "admin" && newRole != "member" && newRole != "viewer" {
		return errors.New("role tidak valid (pilih: admin, member, atau viewer)")
	}

	board, err := u.boardRepo.FindByID(boardID, currentUserID)
	if err != nil || board == nil {
		return errors.New("papan tidak ditemukan")
	}

	currentUserRole, _ := u.GetUserRole(boardID, currentUserID)
	if currentUserRole != "admin" && board.UserID != currentUserID {
		return errors.New("hanya admin yang dapat mengubah role anggota")
	}

	if targetUserID == board.UserID {
		return errors.New("role pemilik utama papan tidak dapat diubah")
	}

	return u.memberRepo.UpdateRole(boardID, targetUserID, newRole)
}

// 4. Mengeluarkan anggota dari board
func (u *boardMemberUsecase) RemoveMember(boardID uint, currentUserID uint, targetUserID uint) error {
	board, err := u.boardRepo.FindByID(boardID, currentUserID)
	if err != nil || board == nil {
		return errors.New("papan tidak ditemukan")
	}

	if targetUserID == board.UserID {
		return errors.New("pemilik papan tidak dapat dikeluarkan dari papan")
	}

	currentUserRole, _ := u.GetUserRole(boardID, currentUserID)
	// Boleh jika yang menghapus adalah admin ATAU anggota mengeluarkan dirinya sendiri (Leave)
	if currentUserRole != "admin" && board.UserID != currentUserID && currentUserID != targetUserID {
		return errors.New("hanya admin yang dapat mengeluarkan anggota lain")
	}

	return u.memberRepo.RemoveMember(boardID, targetUserID)
}

// 5. Cek role user di board
func (u *boardMemberUsecase) GetUserRole(boardID uint, userID uint) (string, error) {
	board, err := u.boardRepo.FindByID(boardID, userID)
	if err == nil && board != nil && board.UserID == userID {
		return "admin", nil
	}

	member, err := u.memberRepo.FindByBoardAndUser(boardID, userID)
	if err != nil || member == nil {
		return "", errors.New("bukan anggota board")
	}

	return member.Role, nil
}
