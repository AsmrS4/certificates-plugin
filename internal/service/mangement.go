package service

import (
	"github.com/AsmrS4/certificates-plugin/internal/enums"
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

func (cm *CertificateManagement) ProcessRequest(id int64) error {
	return nil
}

func (cm *CertificateManagement) RejectCertificateRequest(id int64) error {
	return nil
}

func (cm *CertificateManagement) UploadCertificateRequest(id int64) error {
	return nil
}

func (cm *CertificateManagement) FindRequestByID(id int64) error {
	return nil
}

func (cm *CertificateManagement) FindAllRequests(st enums.CertificateStatus, offset int, limit int) ([]models.CertificateApplication, int64, error) {
	var orders []models.CertificateApplication
	var total int64
	var err error

	if st == "" {
		orders, total, err = cm.appRepo.FindAllRequests(offset, limit)
	} else {
		orders, total, err = cm.appRepo.FindAllRequestsWithStatus(st, offset, limit)
	}

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (cmh *CertificateManagement) changeStatus(st *enums.CertificateStatus) error {
	return nil
}
