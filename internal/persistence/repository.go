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
	IsRejected(id int64) (bool, error)
	IsPending(id int64) (bool, error)
	IsExists(id int64) (bool, error)
	Cancel(id int64) error
	Reject(id int64, reason string) error
	Prepare(id int64) error
	Done(id int64) error
}

type CertificateRepo interface {
	Save(c *models.Certificate) (int, error)
	FindByID(id int64) (*models.Certificate, error)
	FindAll() ([]models.Certificate, error)
	FindByReceiverIDAndID(receiverId int64, id int64) (*models.Certificate, error)
	FindAllByReceiverID(receiverId int64) ([]models.Certificate, error)
}
