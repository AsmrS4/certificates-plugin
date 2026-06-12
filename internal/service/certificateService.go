package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/AsmrS4/certificates-plugin/internal/enums"
	"github.com/AsmrS4/certificates-plugin/internal/models"
	"github.com/AsmrS4/certificates-plugin/internal/persistence"
	wasmplugin "github.com/SuperBotForge/sdk/go-sdk"
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
	userInfo, err := ctx.GetUserInfo(studentID)
	if err != nil {
		return 0, fmt.Errorf("%s", err.Error())
	}

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
		FullName:        userInfo.FullName,
	}

	id, err := c.appRepo.Save(newOrder)
	if err != nil {
		return 0, fmt.Errorf("%s", err.Error())
	}

	return id, err
}

func (c *CertificateService) FindByID(id int64) (*models.CertificateApplication, error) {
	found, err := c.appRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrOrderNotFound
		}
		return nil, err
	}

	return found, nil
}

func (c *CertificateService) GetCertificateFileByOrderID(orderID int64) (*models.CertificateShort, error) {
	cert, err := c.certRepo.FindCertificateByOrderID(orderID)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func (c *CertificateService) CancelCertificateOrder(id int64) (bool, error) {
	err := c.appRepo.Cancel(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, models.ErrOrderNotFound
		}
		return false, err
	}

	return true, nil
}

func (c *CertificateService) FindAllActive(userID int64) ([]models.CertificateApplication, error) {
	orders, err := c.appRepo.FindAllActive(userID)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return orders, nil
}

func (c *CertificateService) FindAllWithStatus(userID int64, cs string) ([]models.CertificateApplication, error) {
	st, err := validateCertificateApplicationStatus(cs)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	orders, err := c.appRepo.FindAllWithStatus(userID, st)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	return orders, nil
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

func validateCertificateApplicationStatus(application_status string) (enums.CertificateStatus, error) {
	if application_status == "" {
		return "", fmt.Errorf("Order status is required")
	}
	st, err := enums.ParseCertificateStatus(application_status)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	return st, nil
}
