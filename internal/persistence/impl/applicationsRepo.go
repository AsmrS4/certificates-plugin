package impl

import (
	"database/sql"

	"github.com/AsmrS4/certificates-plugin/internal/models"
	repository "github.com/AsmrS4/certificates-plugin/internal/persistence"
)

var _ repository.CertificateApplicationRepo = (*CertAppRepoImpl)(nil)

type CertAppRepoImpl struct {
	db *sql.DB
}

func NewApplicationRepo(db *sql.DB) *CertAppRepoImpl {
	return &CertAppRepoImpl{db: db}
}

func (r *CertAppRepoImpl) Save(c *models.CertificateApplication) (int64, error) {
	var id int64

	err := r.db.QueryRow(
		`INSERT INTO certificate_applications(student_id, certificate_type, obtain_method)
		 VALUES ($1, $2, $3) RETURNING id`,
		c.StudentID, c.CertificateType, c.ObtainMethod,
	).Scan(&id)

	return id, err
}

func (r *CertAppRepoImpl) Cancel(id int64) error {
	panic("unimplemented")
}

func (r *CertAppRepoImpl) Done(id int64) error {
	panic("unimplemented")
}

func (r *CertAppRepoImpl) FindAllActive() ([]models.CertificateApplication, error) {
	panic("unimplemented")
}

func (r *CertAppRepoImpl) FindByID(id int64) (*models.CertificateApplication, error) {
	panic("unimplemented")
}

func (r *CertAppRepoImpl) Prepare(id int64) error {
	panic("unimplemented")
}

func (r *CertAppRepoImpl) Reject(id int64) error {
	panic("unimplemented")
}

func (c *CertAppRepoImpl) Update() {
	panic("unimplemented")
}
