package impl

import (
	"database/sql"

	"github.com/AsmrS4/certificates-plugin/internal/models"
	repository "github.com/AsmrS4/certificates-plugin/internal/persistence"
)

// указание компилятору проверить, что реализация контракта есть
var _ repository.CertificateApplicationRepo = (*CertAppRepoImpl)(nil)

type CertAppRepoImpl struct {
	db *sql.DB
}

func (c *CertAppRepoImpl) Cancel(id int64) error {
	panic("unimplemented")
}

func (c *CertAppRepoImpl) Done(id int64) error {
	panic("unimplemented")
}

func (c *CertAppRepoImpl) FindAllActive() ([]models.CertificateApplication, error) {
	panic("unimplemented")
}

func (c *CertAppRepoImpl) FindByID(id int64) (*models.CertificateApplication, error) {
	panic("unimplemented")
}

func (c *CertAppRepoImpl) Prepare(id int64) error {
	panic("unimplemented")
}

func (c *CertAppRepoImpl) Reject(id int64) error {
	panic("unimplemented")
}

func (*CertAppRepoImpl) Save(c *models.CertificateApplication) (int, error) {
	panic("unimplemented")
}

func (c *CertAppRepoImpl) Update() {
	panic("unimplemented")
}

func New(db *sql.DB) *CertAppRepoImpl {
	return &CertAppRepoImpl{db: db}
}
