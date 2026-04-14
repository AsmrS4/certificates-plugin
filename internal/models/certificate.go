package models

import (
	"time"

	enums "github.com/AsmrS4/certificates-plugin/internal/enums/certificate"
)

type (
	CertificateStatus = enums.CertificateStatus
	CertificateType   = enums.CertificateType
	ObtainMethod      = enums.ObtainMethod
)

type CertificateApplication struct {
	ID                int64             `json:"id"`
	StudentID         int64             `json:"student_id"`
	ApplicationStatus CertificateStatus `json:"application_status"`
	CertificateType   CertificateType   `json:"certificate_type"`
	ObtainMethod      ObtainMethod      `json:"obtain_method"`
	CreatedAt         time.Time         `json:"created_at"`
}

type Certificate struct {
	ID         int64     `json:"id"`
	ReceiverID int64     `json:"receiver_id"`
	AuthorID   int64     `json:"author_id"`
	FileName   string    `json:"file_name,omitempty"`
	StorageURL string    `json:"storage_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
