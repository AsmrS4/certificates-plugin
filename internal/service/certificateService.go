package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	comment := ctx.Param("comment")
	studentID := ctx.Messenger.UserID

	validatedType, err := validateRequiredType(certificateType)
	if err != nil {
		ctx.LogError(fmt.Sprintf("validate certificate type failed: %s", fmt.Sprint(err.Error())))
		return 0, fmt.Errorf("%s", err.Error())
	}
	validatedMethod, err := validateObtainMethod(obtainMethod)
	if err != nil {
		ctx.LogError(fmt.Sprintf("validate obtain method failed: %s", fmt.Sprint(err.Error())))
		return 0, fmt.Errorf("%s", err.Error())
	}

	trimmedComment := strings.TrimSpace(comment)
	if trimmedComment != "" {
		if len(trimmedComment) > 255 {
			ctx.LogError("comment is too long")
			return 0, fmt.Errorf("comment_too_long")
		}
	}
	if trimmedComment == "" && validatedType == enums.Common {
		ctx.LogError("comment is required for common certificate type")
		return 0, fmt.Errorf("comment_req")
	}

	userInfo, err := ctx.GetUserInfo(studentID)
	var fullName string
	if err != nil {
		ctx.LogError(fmt.Sprintf("get user info failed: %s", fmt.Sprint(err.Error())))
	} else {
		fullName = userInfo.FullName
	}

	newOrder := &models.CertificateApplication{
		StudentID:       studentID,
		CertificateType: validatedType,
		ObtainMethod:    validatedMethod,
		FullName:        fullName,
		Comment:         trimmedComment,
	}

	id, err := c.appRepo.Save(newOrder)
	if err != nil {
		ctx.LogError(fmt.Sprintf("save certificate order failed: %s", fmt.Sprint(err.Error())))
		return 0, fmt.Errorf("%s", "order_creation_error")
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
		return "", fmt.Errorf("certificate_type_req")
	}
	t, err := enums.ParseCertificateType(requiredType)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	return t, nil
}

func validateObtainMethod(requiredMethod string) (enums.ObtainMethod, error) {
	if requiredMethod == "" {
		return "", fmt.Errorf("obtain_method_req")
	}
	m, err := enums.ParseObtainMethod(requiredMethod)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	return m, nil
}

func validateCertificateApplicationStatus(application_status string) (enums.CertificateStatus, error) {
	if application_status == "" {
		return "", fmt.Errorf("order_status_req")
	}
	st, err := enums.ParseCertificateStatus(application_status)
	if err != nil {
		return "", fmt.Errorf("%s", err.Error())
	}
	return st, nil
}
