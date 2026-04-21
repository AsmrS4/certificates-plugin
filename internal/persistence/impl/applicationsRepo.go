package impl

import (
	"database/sql"

	"github.com/AsmrS4/certificates-plugin/internal/enums"
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

func (r *CertAppRepoImpl) FindByID(id int64) (*models.CertificateApplication, error) {
	row := r.db.QueryRow(
		`SELECT id, student_id, application_status, certificate_type, obtain_method, created_at FROM certificate_applications WHERE id = $1`, id)

	var found models.CertificateApplication

	err := scanItem(&found, row)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &found, nil
}

func (r *CertAppRepoImpl) FindAllActive(userID int64) ([]models.CertificateApplication, error) {
	query := `
        SELECT id, student_id, application_status, certificate_type, obtain_method, created_at
        FROM certificate_applications
        WHERE student_id = $1 AND application_status NOT IN ('Cancelled', 'Rejected')
        ORDER BY created_at DESC LIMIT 10
		`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanItems(rows)
}

func (r *CertAppRepoImpl) FindAllWithStatus(userID int64, st enums.CertificateStatus) ([]models.CertificateApplication, error) {
	query := `
        SELECT id, student_id, application_status, certificate_type, obtain_method, created_at 
        FROM certificate_applications
        WHERE student_id = $1 AND application_status = $2
        ORDER BY created_at DESC LIMIT 10
		`
	rows, err := r.db.Query(query, userID, st)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanItems(rows)
}

func (r *CertAppRepoImpl) Cancel(id int64) error {
	_, err := r.db.Exec(`UPDATE certificate_applications SET application_status = 'Cancelled' WHERE id = $1`, id)
	return err
}

func (r *CertAppRepoImpl) Done(id int64) error {
	_, err := r.db.Exec(`UPDATE certificate_applications SET application_status = 'Done' WHERE id = $1`, id)
	return err
}

func (r *CertAppRepoImpl) Prepare(id int64) error {
	_, err := r.db.Exec(`UPDATE certificate_applications SET application_status = 'Prepare' WHERE id = $1`, id)
	return err
}

func (r *CertAppRepoImpl) Reject(id int64) error {
	_, err := r.db.Exec(`UPDATE certificate_applications SET application_status = 'Reject' WHERE id = $1`, id)
	return err
}

func (c *CertAppRepoImpl) Update() {
	panic("unimplemented")
}

func scanItems(rows *sql.Rows) ([]models.CertificateApplication, error) {
	var items []models.CertificateApplication
	for rows.Next() {
		var item models.CertificateApplication
		if err := rows.Scan(&item.ID, &item.StudentID, &item.ApplicationStatus, &item.CertificateType, &item.ObtainMethod, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanItem(item *models.CertificateApplication, row *sql.Row) error {
	err := row.Scan(&item.ID, &item.StudentID, &item.ApplicationStatus, &item.CertificateType, &item.ObtainMethod, &item.CreatedAt)
	return err
}
