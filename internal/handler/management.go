package handler

import (
	"fmt"
	"strconv"

	"github.com/AsmrS4/certificates-plugin/internal/enums"
	"github.com/AsmrS4/certificates-plugin/internal/models"
	"github.com/AsmrS4/certificates-plugin/internal/service"
	wasmplugin "github.com/StaZisS/SuperBotGo/sdk/go-plugin"
)

type CertificateManagementHandler struct {
	cService  *service.CertificateService
	cmService *service.CertificateManagement
	cat       *wasmplugin.Catalog
}

func NewManagementHandler(cs *service.CertificateService, cms *service.CertificateManagement, cat *wasmplugin.Catalog) *CertificateManagementHandler {
	return &CertificateManagementHandler{cService: cs, cmService: cms, cat: cat}
}

func (cmh *CertificateManagementHandler) ProcessRequest(ctx *wasmplugin.EventContext) error {
	return nil
}

func (cmh *CertificateManagementHandler) UploadCertificate(ctx *wasmplugin.EventContext) error {
	return nil
}

func (cmh *CertificateManagementHandler) RejectCertificateRequest(ctx *wasmplugin.EventContext) error {
	return nil
}

func (cmh *CertificateManagementHandler) FindRequests(ctx *wasmplugin.EventContext) error {

	status, offset, limit, vErr := cmh.validateRequestParams(ctx.HTTP)
	if vErr != nil {
		ctx.JSON(400, map[string]string{"error": vErr.Error()})
	}

	var orders []models.CertificateApplication
	var err error

	orders, total, err := cmh.cmService.FindAllRequests(status, offset, limit)

	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Internal server error"})
		ctx.LogError("http error: " + err.Error())
	}

	body := map[string]interface{}{
		"orders": orders,
		"limit":  limit,
		"offset": offset,
		"total":  total,
	}

	ctx.JSON(200, body)

	return nil
}

func (cmh *CertificateManagementHandler) FindRequestDetails(ctx *wasmplugin.EventContext) error {
	return nil
}

func (cmh *CertificateManagementHandler) validateRequestParams(params *wasmplugin.HTTPEventData) (enums.CertificateStatus, int, int, error) {
	status := params.Query["status"]
	rOffset := params.Query["offset"]
	rLimit := params.Query["limit"]
	var st enums.CertificateStatus
	var err error

	if status == "" && rOffset == "" && rLimit == "" {
		return "", 0, 10, nil
	}

	if status != "" {
		st, err = enums.ParseCertificateStatus(status)
		if err != nil {
			return "", 0, 0, fmt.Errorf("%s", err.Error())
		}
	}

	limit, err := strconv.ParseInt(rLimit, 10, 32)
	if err != nil {
		return "", 0, 0, fmt.Errorf("The limit value must be positive.")
	}

	if limit == 0 {
		return "", 0, 0, fmt.Errorf("The limit value must be greater than zero.")
	}

	offset, err := strconv.ParseInt(rOffset, 10, 32)
	if err != nil {
		return "", 0, 0, fmt.Errorf("The offset value must be positive.")
	}

	return st, int(offset), int(limit), nil
}
