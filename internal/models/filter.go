package models

import "github.com/AsmrS4/certificates-plugin/internal/enums"

type FilterParams struct {
	StudentID         *int64
	CertificateStatus *enums.CertificateStatus
	CertificateType   *enums.CertificateType
	Offset            int64
	Limit             int64
}
