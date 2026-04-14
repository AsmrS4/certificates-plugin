package certificate

import "github.com/AsmrS4/certificates-plugin/internal/models"

type CertificateServiceImpl struct {
}

func (сs CertificateServiceImpl) Create() (models.CertificateApplication, error) {
	return models.CertificateApplication{}, nil
}
