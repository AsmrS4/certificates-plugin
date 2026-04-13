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

type CertificateRequest struct {
	ID           int64             `json:"id"`
	StudentID    int64             `json:"student_id"`
	Status       CertificateStatus `json:"status"`
	Type         CertificateType   `json:"certificate_type"`
	ObtainMethod ObtainMethod      `json:"obtain_method"`
	CreatedAt    time.Time         `json:"created_at"`
}

type Certificate struct {
	ID         int64     `json:"id"`
	ReceiverID int64     `json:"receiver_id"`
	AuthorID   int64     `json:"author_id"`
	CreatedAt  time.Time `json:"created_at"`
	DocumentID int64     `json:"document_id,omitempty"`
}

type Document struct {
	ID         int64  `json:"id"`
	StorageURL string `json:"storage_url"`
	FileName   string `json:"file_name"`
	UploadedAt string `json:"uploaded_at"`
}
