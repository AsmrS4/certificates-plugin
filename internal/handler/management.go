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
	var orders []models.CertificateApplication
	var err error

	filters, err := cmh.validateRequestParams(ctx.HTTP)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": err.Error()})
		return nil
	}

	orders, total, err := cmh.cmService.FindAllRequests(*filters)

	if err != nil {
		ctx.JSON(500, map[string]string{"error": "Internal server error"})
		ctx.LogError("http error: " + err.Error())
		return nil
	}

	body := map[string]interface{}{
		"orders": orders,
		"limit":  filters.Limit,
		"offset": filters.Offset,
		"total":  total,
	}

	ctx.JSON(200, body)

	return nil
}

func (cmh *CertificateManagementHandler) FindRequestDetails(ctx *wasmplugin.EventContext) error {
	return nil
}

func (cmh *CertificateManagementHandler) validateRequestParams(params *wasmplugin.HTTPEventData) (*models.FilterParams, error) {
	status := params.Query["status"]
	certificateType := params.Query["type"]
	userID := params.Query["user_id"]
	rOffset := params.Query["offset"]
	rLimit := params.Query["limit"]

	var filters models.FilterParams = models.FilterParams{
		Limit:  10,
		Offset: 0,
	}

	if status == "" && certificateType == "" && userID == "" && rOffset == "" && rLimit == "" {
		return &filters, nil
	}

	if status != "" {
		val, err := enums.ParseCertificateStatus(status)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		filters.CertificateStatus = &val
	}

	if certificateType != "" {
		val, err := enums.ParseCertificateType(certificateType)
		if err != nil {
			return nil, fmt.Errorf("%s", err.Error())
		}
		filters.CertificateType = &val
	}

	if userID != "" {
		val, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Incorrect id format. Required integer")
		}
		filters.StudentID = &val
	}

	if rLimit != "" {
		limit, err := strconv.ParseInt(rLimit, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("The limit value must be positive.")
		}
		if limit == 0 {
			return nil, fmt.Errorf("The limit value must be greater than zero.")
		}
		filters.Limit = limit
	}

	if rOffset != "" {
		offset, err := strconv.ParseInt(rOffset, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("The offset value must be positive.")
		}
		filters.Offset = offset
	}

	return &filters, nil
}
