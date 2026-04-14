package enums

type CertificateStatus string
type CertificateType string
type ObtainMethod string

const (
	Pending   CertificateStatus = "Pending"
	Rejected  CertificateStatus = "Rejected"
	Prepare   CertificateStatus = "Prepare"
	Done      CertificateStatus = "Done"
	Cancelled CertificateStatus = "Cancelled"
)

const (
	StudyPeriod    CertificateType = "StudyPeriod"
	Academic       CertificateType = "Academic"
	Recommendation CertificateType = "Recommendation"
	Common         CertificateType = "Common"
)

const (
	Electronic ObtainMethod = "Electronic"
	Paper      ObtainMethod = "Paper"
)
