package usecase

import (
	"ecommerce-api/internal/domain"
	"errors"
	"os"
	"strings"
)

type attachmentUsecase struct {
	attachmentRepo domain.AttachmentRepository
	cardRepo       domain.CardRepository
}

func NewAttachmentUsecase(attachmentRepo domain.AttachmentRepository, cardRepo domain.CardRepository) domain.AttachmentUsecase {
	return &attachmentUsecase{
		attachmentRepo: attachmentRepo,
		cardRepo:       cardRepo,
	}
}

func (u *attachmentUsecase) CreateAttachment(attachment *domain.Attachment) error {
	if attachment.CardID == 0 {
		return errors.New("id kartu tidak valid")
	}
	if strings.TrimSpace(attachment.FileName) == "" {
		return errors.New("nama file tidak boleh kosong")
	}

	// 1. Pastikan kartu ada di database
	card, err := u.cardRepo.FindByID(attachment.CardID)
	if err != nil || card == nil {
		return errors.New("kartu tidak ditemukan")
	}

	// 2. Simpan rekaman lampiran ke database
	return u.attachmentRepo.Create(attachment)
}

func (u *attachmentUsecase) GetAttachmentsByCardID(cardID uint) ([]domain.Attachment, error) {
	return u.attachmentRepo.FindByCardID(cardID)
}

func (u *attachmentUsecase) DeleteAttachment(id uint) error {
	// 1. Cari data lampiran terlebih dahulu
	attachment, err := u.attachmentRepo.FindByID(id)
	if err != nil || attachment == nil {
		return errors.New("lampiran tidak ditemukan")
	}

	// 2. Hapus berkas fisik dari disk penyimpanan server (jika ada)
	filePathOnDisk := "." + attachment.FilePath
	_ = os.Remove(filePathOnDisk) // Bersihkan file tanpa memblokir jika file sudah terhapus

	// 3. Hapus data dari database
	return u.attachmentRepo.Delete(id)
}
