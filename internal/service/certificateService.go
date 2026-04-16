package service

import (
	"fmt"

	"github.com/AsmrS4/certificates-plugin/internal/enums"
	"github.com/AsmrS4/certificates-plugin/internal/models"
	"github.com/AsmrS4/certificates-plugin/internal/persistence"
	wasmplugin "github.com/StaZisS/SuperBotGo/sdk/go-plugin"
)

type CertificateService struct {
	appRepo  persistence.CertificateApplicationRepo
	certRepo persistence.CertificateRepo
}

func New(appRepo persistence.CertificateApplicationRepo, certRepo persistence.CertificateRepo) *CertificateService {
	return &CertificateService{appRepo: appRepo, certRepo: certRepo}
}

func (c *CertificateService) CreateCertificateOrder(ctx *wasmplugin.EventContext) (int64, error) {
	certificateType := ctx.Param("type")
	obtainMethod := ctx.Param("obtain_method")
	studentID := ctx.Messenger.UserID

	validatedType, err := validateRequiredType(certificateType)
	if err != nil {
		return 0, fmt.Errorf("%s", err.Error())
	}
	validatedMethod, err := validateObtainMethod(obtainMethod)
	if err != nil {
		return 0, fmt.Errorf("%s", err.Error())
	}

	newOrder := &models.CertificateApplication{
		StudentID:       studentID,
		CertificateType: validatedType,
		ObtainMethod:    validatedMethod,
	}

	id, err := c.appRepo.Save(newOrder)
	if err != nil {
		return 0, fmt.Errorf("%s", err.Error())
	}

	return id, err
}

func (c *CertificateService) CancelCertificateOrder(ctx *wasmplugin.EventContext, id int64) error {
	return nil
}

func validateRequiredType(requiredType string) (enums.CertificateType, error) {
	if requiredType == "" {
		return "", fmt.Errorf("Certificate type is required")
	}
	t, err := enums.ParseCertificateType(requiredType)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	return t, nil
}

func validateObtainMethod(requiredMethod string) (enums.ObtainMethod, error) {
	if requiredMethod == "" {
		return "", fmt.Errorf("Obtain method is required")
	}
	m, err := enums.ParseObtainMethod(requiredMethod)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	return m, nil
}
