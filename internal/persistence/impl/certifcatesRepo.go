package impl

import (
	"database/sql"

	"github.com/AsmrS4/certificates-plugin/internal/models"
	repository "github.com/AsmrS4/certificates-plugin/internal/persistence"
)

// указание компилятору проверить, что реализация контракта есть
var _ repository.CertificateRepo = (*CertRepoImpl)(nil)

type CertRepoImpl struct {
	db *sql.DB
}

func (c *CertRepoImpl) FindAll() ([]models.Certificate, error) {
	panic("unimplemented")
}

func (c *CertRepoImpl) FindAllByReceiverID(receiverId int64) ([]models.Certificate, error) {
	panic("unimplemented")
}

func (c *CertRepoImpl) FindByID(id int64) (*models.Certificate, error) {
	panic("unimplemented")
}

func (c *CertRepoImpl) FindByReceiverIDAndID(receiverId int64, id int64) (*models.Certificate, error) {
	panic("unimplemented")
}

func (*CertRepoImpl) Save(c *models.Certificate) (int, error) {
	panic("unimplemented")
}
