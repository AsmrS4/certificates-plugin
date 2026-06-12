package persistence

import (
	"github.com/AsmrS4/certificates-plugin/internal/enums"
	"github.com/AsmrS4/certificates-plugin/internal/models"
)

type CertificateApplicationRepo interface {
	Save(c *models.CertificateApplication) (int64, error)
	FindByID(id int64) (*models.CertificateApplication, error)
	FindAllActive(userID int64) ([]models.CertificateApplication, error)
	FindAllWithStatus(userID int64, st enums.CertificateStatus) ([]models.CertificateApplication, error)
	FindAllRequests(params models.FilterParams) ([]models.CertificateApplication, int64, error)
	Update()
	IsPaper(id int64) (bool, error)
	IsRejected(id int64) (bool, error)
	IsPending(id int64) (bool, error)
	IsExists(id int64) (bool, error)
	IsProcessing(id int64) (bool, error)
	Cancel(id int64) error
	Reject(id int64, reason string) (int64, int64, error)
	Prepare(id int64) (int64, int64, error)
	Done(id int64) (int64, int64, error)
}

type CertificateRepo interface {
	Save(c *models.Certificate) (int64, int64, error)
	FindAll() ([]models.Certificate, error)
	FindAllByReceiverID(receiverId int64) ([]models.Certificate, error)
	FindCertificateByOrderID(orderID int64) (*models.CertificateShort, error)
}
