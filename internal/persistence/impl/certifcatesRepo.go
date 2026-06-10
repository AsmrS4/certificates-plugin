package impl

import (
	"database/sql"

	repository "github.com/AsmrS4/certificates-plugin/internal/persistence"

	"github.com/AsmrS4/certificates-plugin/internal/models"
)

var _ repository.CertificateRepo = (*CertRepoImpl)(nil)

type CertRepoImpl struct {
	db *sql.DB
}

func (cr *CertRepoImpl) FindAll() ([]models.Certificate, error) {
	panic("unimplemented")
}

func (cr *CertRepoImpl) FindAllByReceiverID(receiverId int64) ([]models.Certificate, error) {
	panic("unimplemented")
}

func (cr *CertRepoImpl) Save(c *models.Certificate) (int64, error) {
	tx, err := cr.db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var id int64
	err = tx.QueryRow(
		`INSERT INTO certificates(author_id, order_id, file_name, storage_url)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		c.AuthorID, c.OrderID, c.FileName, c.StorageURL,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(`UPDATE certificate_applications SET application_status = 'Done' WHERE id = $1`, c.OrderID)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, err
}

func NewCertRepo(db *sql.DB) *CertRepoImpl {
	return &CertRepoImpl{db: db}
}
