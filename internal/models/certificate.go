package models

import (
	certificates "github.com/AsmrS4/certificates-plugin/internal/enums"
)

type (
	CertificateStatus = certificates.CertificateStatus
	CertificateType   = certificates.CertificateType
	ObtainMethod      = certificates.ObtainMethod
)

type CertificateApplication struct {
	ID                int64             `json:"id"`
	StudentID         int64             `json:"student_id"`
	ApplicationStatus CertificateStatus `json:"application_status"`
	CertificateType   CertificateType   `json:"certificate_type"`
	ObtainMethod      ObtainMethod      `json:"obtain_method"`
	RejectionReason   string            `json:"rejection_reason"`
	CreatedAt         string            `json:"created_at"`
}

type Certificate struct {
	ID         int64  `json:"id"`
	ReceiverID int64  `json:"receiver_id"`
	AuthorID   int64  `json:"author_id"`
	FileName   string `json:"file_name,omitempty"`
	StorageURL string `json:"storage_url,omitempty"`
	CreatedAt  string `json:"created_at"`
}
