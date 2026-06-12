package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/AsmrS4/certificates-plugin/internal/enums"
	"github.com/AsmrS4/certificates-plugin/internal/models"
	"github.com/AsmrS4/certificates-plugin/internal/service"

	wasmplugin "github.com/SuperBotForge/sdk/go-sdk"
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
	id := ctx.HTTP.Query["id"]
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Incorrect id format. Int or long value is required."})
		return nil
	}

	orderID, studentID, err := cmh.cmService.ProcessRequest(id64)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(404, map[string]string{"error": "Order not found"})
			return nil
		}
		if errors.Is(err, models.ErrOrderNotFound) {
			ctx.JSON(404, map[string]string{"error": err.Error()})
			return nil
		}
		if errors.Is(err, models.ErrOrderNotPending) {
			ctx.JSON(400, map[string]string{"error": err.Error()})
			return nil
		}

		ctx.JSON(500, map[string]string{"error": "Internal server error"})
		ctx.LogError(fmt.Sprintf("unexpected rejection order error: %s", err.Error()))
		return nil
	}

	var event = models.OrderEvent{
		UserID:      studentID,
		OrderID:     orderID,
		OrderStatus: string(enums.Prepare),
	}

	err = wasmplugin.PublishEvent("certificate_order.updated", event)
	if err != nil {
		ctx.LogError(fmt.Sprintf("failed send notification after prepare: %s", err.Error()))
	}
	ctx.JSON(200, true)
	return nil
}

func (cmh *CertificateManagementHandler) UploadCertificate(ctx *wasmplugin.EventContext) error {
	orderID := ctx.HTTP.Query["id"]
	id64, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Incorrect id format. Int or long value is required."})
		return nil
	}

	rawPayload := ctx.HTTP.Body
	if rawPayload == "" {
		ctx.JSON(400, map[string]string{"error": "Payload data is required."})
		return nil
	}

	var payload struct {
		FileName string `json:"file_name,omitempty"`
		FileID   string `json:"file_id,omitempty"`
	}

	err = json.Unmarshal([]byte(ctx.HTTP.Body), &payload)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Incorrect payload. Payload must contain \"file_name\" and \"file_id\" fields."})
		return nil
	}
	storageUrl, err := ctx.FileURL(payload.FileID)
	if err != nil {
		ctx.LogError(fmt.Sprintf("error getting file url for file id %s: %s", payload.FileID, err.Error()))
		ctx.JSON(400, map[string]string{"error": "Incorrect file_id. File with provided id doesn't exist."})
		return nil
	}

	_, studentID, err := cmh.cmService.UploadCertificateRequest(models.CertificateData{
		OrderID:    id64,
		Filename:   payload.FileName,
		FileID:     payload.FileID,
		StorageURL: storageUrl,
	})
	if err != nil {
		ctx.LogError(fmt.Sprintf("file store in plugin error: %s", err.Error()))
		ctx.JSON(500, map[string]string{"error": "Internal server error."})
		return nil
	}
	ctx.JSON(201, map[string]string{"message": "file uploaded and saved successfully."})

	var event = models.OrderEvent{
		UserID:      studentID,
		OrderID:     id64,
		OrderStatus: string(enums.Done),
	}

	err = wasmplugin.PublishEvent("certificate_order.updated", event)
	if err != nil {
		ctx.LogError(fmt.Sprintf("failed send notification after complete: %s", err.Error()))
	}

	return nil
}

func (cmh *CertificateManagementHandler) FinishProcessingPaperCertificate(ctx *wasmplugin.EventContext) error {
	orderID := ctx.HTTP.Query["id"]
	id64, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Incorrect id format. Int or long value is required."})
		return nil
	}

	exists, err := cmh.cmService.ExistsByID(id64)
	if err != nil {
		if !exists {
			ctx.JSON(404, map[string]string{"error": "Order not found."})
			return nil
		}
		ctx.JSON(500, map[string]string{"error": "Internal server error"})
		ctx.LogError(fmt.Sprintf("unexpected error while finish processing paper order: %s", err.Error()))
		return nil
	}

	orderID64, studentID, err := cmh.cmService.FinishProcessingOrder(id64)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(404, map[string]string{"error": "Order not found"})
			return nil
		}
		if errors.Is(err, models.ErrOrderNotFound) {
			ctx.JSON(404, map[string]string{"error": err.Error()})
			return nil
		}
		if errors.Is(err, models.ErrOrderNotInPrepare) {
			ctx.JSON(400, map[string]string{"error": err.Error()})
			return nil
		}

		ctx.JSON(500, map[string]string{"error": "Internal server error"})
		ctx.LogError(fmt.Sprintf("unexpected error while finish processing paper order: %s", err.Error()))
		return nil
	}

	var event = models.OrderEvent{
		UserID:      studentID,
		OrderID:     orderID64,
		OrderStatus: string(enums.Done),
	}

	err = wasmplugin.PublishEvent("certificate_order.updated", event)
	if err != nil {
		ctx.LogError(fmt.Sprintf("failed send notification after process complete: %s", err.Error()))
	}
	ctx.JSON(200, true)
	return nil
}

func (cmh *CertificateManagementHandler) RejectCertificateRequest(ctx *wasmplugin.EventContext) error {
	id := ctx.HTTP.Query["id"]
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Incorrect id format. Int or long value is required."})
		return nil
	}

	raw := ctx.HTTP.Body
	var body map[string]string
	err = json.Unmarshal([]byte(raw), &body)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Incorrect JSON format."})
		return nil
	}

	reason, ok := body["reason"]
	if !ok {
		ctx.JSON(400, map[string]string{"error": "Reason field is required."})
		return nil
	}
	if len(strings.TrimSpace(reason)) == 0 {
		ctx.JSON(400, map[string]string{"error": "Rejection reason couldn't be blank."})
		return nil
	}

	orderID, studentID, err := cmh.cmService.RejectCertificateRequest(id64, reason)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(404, map[string]string{"error": "Order not found"})
			return nil
		}
		if errors.Is(err, models.ErrOrderNotFound) {
			ctx.JSON(404, map[string]string{"error": err.Error()})
			return nil
		}
		if errors.Is(err, models.ErrAlreadyRejected) {
			ctx.JSON(400, map[string]string{"error": err.Error()})
			return nil
		}
		if errors.Is(err, models.ErrOrderNotPending) {
			ctx.JSON(400, map[string]string{"error": err.Error()})
			return nil
		}

		ctx.JSON(500, map[string]string{"error": "Internal server error"})
		ctx.LogError(fmt.Sprintf("unexpected rejection order error: %s", err.Error()))
		return nil
	}

	var event = models.OrderEvent{
		UserID:      studentID,
		OrderID:     orderID,
		OrderStatus: string(enums.Rejected),
	}

	err = wasmplugin.PublishEvent("certificate_order.updated", event)
	if err != nil {
		ctx.LogError(fmt.Sprintf("failed send notification after reject: %s", err.Error()))
	}
	ctx.JSON(200, true)
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

	pagination := map[string]interface{}{
		"limit":  filters.Limit,
		"offset": filters.Offset,
		"total":  total,
	}

	body := map[string]interface{}{
		"data":       orders,
		"pagination": pagination,
	}

	ctx.JSON(200, body)
	return nil
}

func (cmh *CertificateManagementHandler) FindRequestDetails(ctx *wasmplugin.EventContext) error {
	userID := ctx.HTTP.Query["id"]
	id64, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		ctx.JSON(400, map[string]string{"error": "Incorrect id formate. Required int or long."})
		return nil
	}

	order, err := cmh.cService.FindByID(id64)
	if err != nil {
		if errors.Is(err, models.ErrOrderNotFound) {
			ctx.JSON(404, map[string]string{"error": "Order not found"})
			return nil
		}

		ctx.JSON(500, map[string]string{"error": "Internal server error"})
		ctx.LogError(fmt.Sprintf("unexpected rejection order error: %s", err.Error()))
		return nil
	}

	ctx.JSON(200, order)
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
		if offset == 0 {
			return nil, fmt.Errorf("The offset value must be greater than zero.")
		}
		filters.Offset = offset
	}

	return &filters, nil
}
