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

func (cr *CertRepoImpl) FindCertificateByOrderID(orderID int64) (*models.CertificateShort, error) {
	row := cr.db.QueryRow(
		`SELECT id, COALESCE(file_name, ''), COALESCE(storage_url, '') FROM certificates WHERE order_id = $1`, orderID)

	var found models.CertificateShort
	if err := row.Scan(&found.ID, &found.FileName, &found.StorageURL); err != nil {
		return nil, err
	}
	return &found, nil
}

func (cr *CertRepoImpl) FindAll() ([]models.Certificate, error) {
	panic("unimplemented")
}

func (cr *CertRepoImpl) FindAllByReceiverID(receiverId int64) ([]models.Certificate, error) {
	panic("unimplemented")
}

func (cr *CertRepoImpl) Save(c *models.Certificate) (int64, int64, error) {
	tx, err := cr.db.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var id int64
	err = tx.QueryRow(
		`INSERT INTO certificates(author_id, order_id, file_name, file_id, storage_url)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		c.AuthorID, c.OrderID, c.FileName, c.FileID, c.StorageURL,
	).Scan(&id)
	if err != nil {
		return 0, 0, err
	}
	row := cr.db.QueryRow(
		`SELECT student_id FROM certificate_applications WHERE id = $1`, id)
	var studentID int64
	if err := row.Scan(&studentID); err != nil {
		return 0, 0, err
	}
	_, err = tx.Exec(`UPDATE certificate_applications SET application_status = 'Done' WHERE id = $1`, c.OrderID)
	if err != nil {
		return 0, 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, 0, err
	}

	return c.OrderID, studentID, nil
}

func NewCertRepo(db *sql.DB) *CertRepoImpl {
	return &CertRepoImpl{db: db}
}
