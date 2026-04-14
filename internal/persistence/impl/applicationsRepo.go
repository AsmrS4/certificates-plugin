package impl

import (
	"github.com/AsmrS4/certificates-plugin/internal/models"
	repository "github.com/AsmrS4/certificates-plugin/internal/persistence"
	"github.com/jmoiron/sqlx"
)

// указание компилятору проверить, что реализация контракта есть
var _ repository.CertificateApplicationRepo = (*CertAppRepoImpl)(nil)

type CertAppRepoImpl struct {
	db *sqlx.DB
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

func New(db *sqlx.DB) *CertAppRepoImpl {
	return &CertAppRepoImpl{db: db}
}
