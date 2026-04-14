package persistence

import "github.com/AsmrS4/certificates-plugin/internal/models"

type CertificateApplicationRepo interface {
	Save(c *models.CertificateApplication) (int, error)
	FindByID(id int64) (*models.CertificateApplication, error)
	FindAllActive() ([]models.CertificateApplication, error)
	Update()
	Cancel(id int64) error
	Reject(id int64) error
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
