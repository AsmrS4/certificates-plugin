package external

import (
	enums "github.com/AsmrS4/certificates-plugin/internal/enums/certificate"
	"github.com/AsmrS4/certificates-plugin/internal/models"
)

type (
	CertificateStatus = enums.CertificateStatus
	CertificateType   = enums.CertificateType
	ObtainMethod      = enums.ObtainMethod
)

type CertificateOrders struct {
	Certificates []models.CertificateApplication `json:"certificates"`
}

type CertificateCreateRequest struct {
	StudentID    int64           `json:"student_id"`
	Type         CertificateType `json:"certificate_type"`
	ObtainMethod ObtainMethod    `json:"obtain_method"`
}
