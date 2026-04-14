package impl

import (
	"github.com/AsmrS4/certificates-plugin/internal/models"
	repository "github.com/AsmrS4/certificates-plugin/internal/persistence"
	"github.com/jmoiron/sqlx"
)

// указание компилятору проверить, что реализация контракта есть
var _ repository.DocumentRepo = (*DocumentRepoImpl)(nil)

type DocumentRepoImpl struct {
	db *sqlx.DB
}

// Save implements [persistence.DocumentRepo].
func (*DocumentRepoImpl) Save(d *models.Document) (int, error) {
	panic("unimplemented")
}
