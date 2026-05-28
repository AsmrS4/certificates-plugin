package service

import (
	"github.com/AsmrS4/certificates-plugin/internal/models"
	"github.com/AsmrS4/certificates-plugin/internal/persistence"
)

type CertificateManagement struct {
	appRepo  persistence.CertificateApplicationRepo
	certRepo persistence.CertificateRepo
}

func NewManagementService(ar persistence.CertificateApplicationRepo, cr persistence.CertificateRepo) *CertificateManagement {
	return &CertificateManagement{appRepo: ar, certRepo: cr}
}

func (cm *CertificateManagement) ProcessRequest(id int64) (int64, int64, error) {
	exists, err := cm.appRepo.IsExists(id)
	if err != nil {
		return 0, 0, err
	}
	if !exists {
		return 0, 0, models.ErrOrderNotFound
	}

	pending, err := cm.appRepo.IsPending(id)
	if err != nil {
		return 0, 0, err
	}
	if !pending {
		return 0, 0, models.ErrOrderNotPending
	}

	orderID, studentID, err := cm.appRepo.Prepare(id)
	if err != nil {
		return 0, 0, err
	}

	return orderID, studentID, nil
}

func (cm *CertificateManagement) RejectCertificateRequest(id int64, reason string) (int64, int64, error) {
	exists, err := cm.appRepo.IsExists(id)
	if err != nil {
		return 0, 0, err
	}
	if !exists {
		return 0, 0, models.ErrOrderNotFound
	}

	rejected, err := cm.appRepo.IsRejected(id)
	if err != nil {
		return 0, 0, models.ErrOrderNotFound
	}
	if rejected {
		return 0, 0, models.ErrAlreadyRejected
	}

	pending, err := cm.appRepo.IsPending(id)
	if err != nil {
		return 0, 0, models.ErrOrderNotFound
	}
	if !pending {
		return 0, 0, models.ErrOrderNotPending
	}

	orderID, studentID, err := cm.appRepo.Reject(id, reason)
	if err != nil {
		return 0, 0, err
	}

	return orderID, studentID, nil
}

func (cm *CertificateManagement) UploadCertificateRequest(file models.CertificateData) error {
	exists, err := cm.appRepo.IsExists(file.OrderID)
	if err != nil {
		return err
	}
	if !exists {
		return models.ErrOrderNotFound
	}

	certificateDocument := models.Certificate{
		OrderID:    file.OrderID,
		AuthorID:   1,
		FileName:   file.Filename,
		StorageURL: file.StorageURL,
	}
	_, err = cm.certRepo.Save(&certificateDocument)
	if err != nil {
		return err
	}

	return nil
}

func (cm *CertificateManagement) FindAllRequests(params models.FilterParams) ([]models.CertificateApplication, int64, error) {
	orders, total, err := cm.appRepo.FindAllRequests(params)

	if err != nil {
		return nil, 0, err
	}

	if len(orders) == 0 {
		return make([]models.CertificateApplication, 0), total, nil
	}

	return orders, total, nil
}
