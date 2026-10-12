package usecase_test

import (
	"ecommerce-api/internal/domain"
	"ecommerce-api/internal/usecase"
	"errors"
	"testing"
)

type mockAttachmentRepo struct {
	domain.AttachmentRepository
	attachments map[uint]*domain.Attachment
}

func (m *mockAttachmentRepo) Create(a *domain.Attachment) error {
	a.ID = uint(len(m.attachments) + 1)
	m.attachments[a.ID] = a
	return nil
}

func (m *mockAttachmentRepo) FindByID(id uint) (*domain.Attachment, error) {
	a, ok := m.attachments[id]
	if !ok {
		return nil, errors.New("attachment not found")
	}
	return a, nil
}

func (m *mockAttachmentRepo) FindByCardID(cardID uint) ([]domain.Attachment, error) {
	var list []domain.Attachment
	for _, a := range m.attachments {
		if a.CardID == cardID {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (m *mockAttachmentRepo) Delete(id uint) error {
	delete(m.attachments, id)
	return nil
}

type mockCardRepoForAttachment struct {
	domain.CardRepository
	cards map[uint]*domain.Card
}

func (m *mockCardRepoForAttachment) FindByID(id uint) (*domain.Card, error) {
	c, ok := m.cards[id]
	if !ok {
		return nil, errors.New("card not found")
	}
	return c, nil
}

func TestAttachmentUsecase_CreateAttachment(t *testing.T) {
	cardRepo := &mockCardRepoForAttachment{
		cards: map[uint]*domain.Card{
			1: {ID: 1, Content: "Card Valid"},
		},
	}
	attachmentRepo := &mockAttachmentRepo{attachments: make(map[uint]*domain.Attachment)}
	uc := usecase.NewAttachmentUsecase(attachmentRepo, cardRepo)

	// 1. Sukses simpan attachment
	validAttachment := &domain.Attachment{
		CardID:   1,
		FileName: "wireframe.png",
		FilePath: "/uploads/attachments/wireframe.png",
		FileSize: 1024,
		FileType: "image/png",
	}
	err := uc.CreateAttachment(validAttachment)
	if err != nil {
		t.Fatalf("Gagal simpan attachment valid: %v", err)
	}

	// 2. Gagal: Kartu tidak ditemukan
	invalidCardAttachment := &domain.Attachment{
		CardID:   999,
		FileName: "test.pdf",
		FilePath: "/uploads/test.pdf",
	}
	err = uc.CreateAttachment(invalidCardAttachment)
	if err == nil || err.Error() != "kartu tidak ditemukan" {
		t.Errorf("Harusnya error 'kartu tidak ditemukan', dapat: %v", err)
	}

	// 3. Gagal: Nama file kosong
	emptyNameAttachment := &domain.Attachment{
		CardID:   1,
		FileName: "   ",
		FilePath: "/uploads/test.pdf",
	}
	err = uc.CreateAttachment(emptyNameAttachment)
	if err == nil || err.Error() != "nama file tidak boleh kosong" {
		t.Errorf("Harusnya error 'nama file tidak boleh kosong', dapat: %v", err)
	}
}

func TestAttachmentUsecase_DeleteAttachment(t *testing.T) {
	attachmentRepo := &mockAttachmentRepo{
		attachments: map[uint]*domain.Attachment{
			1: {ID: 1, CardID: 1, FileName: "test.png", FilePath: "/uploads/attachments/dummy.png"},
		},
	}
	uc := usecase.NewAttachmentUsecase(attachmentRepo, nil)

	// 1. Sukses hapus
	err := uc.DeleteAttachment(1)
	if err != nil {
		t.Fatalf("Gagal hapus attachment: %v", err)
	}

	// 2. Gagal: Lampiran tidak ditemukan
	err = uc.DeleteAttachment(999)
	if err == nil || err.Error() != "lampiran tidak ditemukan" {
		t.Errorf("Harusnya error 'lampiran tidak ditemukan', dapat: %v", err)
	}
}
