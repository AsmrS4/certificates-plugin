package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AsmrS4/certificates-plugin/internal/enums"
	"github.com/AsmrS4/certificates-plugin/internal/models"
	"github.com/AsmrS4/certificates-plugin/internal/service"

	wasmplugin "github.com/StaZisS/SuperBotGo/sdk/go-plugin"
)

type CertificateHandler struct {
	service *service.CertificateService
	cat     *wasmplugin.Catalog
}

func NewHandler(service *service.CertificateService, cat *wasmplugin.Catalog) *CertificateHandler {
	return &CertificateHandler{service: service, cat: cat}
}

func (h *CertificateHandler) CreateOrder(ctx *wasmplugin.EventContext) error {

	tr := h.cat.Tr(ctx.Locale())
	id, err := h.service.CreateCertificateOrder(ctx)

	if err != nil {
		ctx.LogError(err.Error())
		ctx.Reply(wasmplugin.NewMessage(tr("order_creation_error")))
		return nil
	}

	localizedMessage := fmt.Sprintf(tr("order_created"), id)

	ctx.Reply(wasmplugin.NewMessage(localizedMessage))
	ctx.Log(fmt.Sprintf("add: certificate order #%d created by user %d", id, ctx.Messenger.UserID))

	return nil
}

func (h *CertificateHandler) FindOrderByID(ctx *wasmplugin.EventContext) error {

	tr := h.cat.Tr(ctx.Locale())
	id := ctx.Param("enter_id")

	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.Reply(wasmplugin.NewMessage(tr("incorrect_id")))
		return nil
	}
	order, err := h.service.FindByID(id64)

	if err != nil {
		ctx.LogError(err.Error())
		ctx.Reply(wasmplugin.NewMessage(tr("error_default")))
		return nil
	}

	if order == nil {
		ctx.LogError(fmt.Sprintf("error: not found certificate order #%d", id64))
		ctx.Reply(wasmplugin.NewMessage(fmt.Sprintf((tr("order_not_found")), id64)))
		return nil
	}

	localizedMessage := h.foundMessageToString(order, tr)

	ctx.Reply(wasmplugin.NewMessage(localizedMessage))
	ctx.Log(fmt.Sprintf("get: found certificate order #%d  %d", order.ID, ctx.Messenger.UserID))

	return nil
}

// TODO: подумать над пагинацией
func (h *CertificateHandler) FindAllActive(ctx *wasmplugin.EventContext) error {
	tr := h.cat.Tr(ctx.Locale())
	status := ctx.Param("status")
	studentID := ctx.Messenger.UserID

	var orders []models.CertificateApplication
	var err error

	if len(strings.TrimSpace(status)) == 0 || status == "Skip" {
		orders, err = h.service.FindAllActive(studentID)
	} else {
		orders, err = h.service.FindAllWithStatus(studentID, status)
	}

	if err != nil {
		ctx.LogError(fmt.Sprintf("server error: %s", err.Error()))
		ctx.Reply(wasmplugin.NewMessage(tr("error_default")))
	}

	if len(orders) == 0 {
		ctx.Reply(wasmplugin.NewMessage(tr("no_orders")))
		return nil
	}

	res := fmt.Sprintf(tr("orders_header"), len(orders)) + "\n\n"
	for _, order := range orders {
		strOrder := h.foundMessageToString(&order, tr)
		res += strOrder + "\n_______________\n"
	}

	ctx.Reply(wasmplugin.NewMessage(res))
	return nil
}

func (h *CertificateHandler) CancelOrderByID(ctx *wasmplugin.EventContext) error {

	tr := h.cat.Tr(ctx.Locale())
	id := ctx.Param("enter_id")
	cancel := ctx.Param("confirm_cancellation")

	if cancel == "no" {
		return nil
	}

	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.Reply(wasmplugin.NewMessage(tr("incorrect_id")))
		return nil
	}

	order, err := h.service.FindByID(id64)
	if err != nil {
		ctx.LogError(err.Error())
		ctx.Reply(wasmplugin.NewMessage(tr("error_default")))
		return nil
	}

	if order == nil {
		ctx.LogError(fmt.Sprintf("error: not found certificate order #%d", id64))
		ctx.Reply(wasmplugin.NewMessage(fmt.Sprintf((tr("order_not_found")), id64)))
		return nil
	}

	status := order.ApplicationStatus
	if status != enums.Pending {
		if status == enums.Cancelled {
			ctx.LogError(fmt.Sprintf("bad request: order #%d already cancelled", id64))
			ctx.Reply(wasmplugin.NewMessage(fmt.Sprintf((tr("order_already_cancelled")), id64)))
			return nil
		}
		if status == enums.Rejected {
			ctx.LogError(fmt.Sprintf("bad request: order #%d was rejected", id64))
			ctx.Reply(wasmplugin.NewMessage(fmt.Sprintf((tr("order_already_rejected")), id64)))
			return nil
		}
		ctx.LogError(fmt.Sprintf("bad request: order #%d has active status", id64))
		ctx.Reply(wasmplugin.NewMessage(fmt.Sprintf((tr("order_is_not_pending")), id64)))
		return nil
	}

	_, err = h.service.CancelCertificateOrder(order.ID)
	if err != nil {
		ctx.LogError(fmt.Sprintf("server error: %s", err.Error()))
		ctx.Reply(wasmplugin.NewMessage(tr("error_default")))
	}

	ctx.Log(fmt.Sprintf("patch: certificate order #%d was cancelled by %d", order.ID, ctx.Messenger.UserID))
	ctx.Reply(wasmplugin.NewMessage(fmt.Sprintf((tr("order_cancelled_successfully")), id64)))
	return nil
}

func (h *CertificateHandler) foundMessageToString(foundOrder *models.CertificateApplication, trans func(key string, args ...any) string) string {
	tr := trans
	var parts []string

	parts = append(parts, fmt.Sprintf("%s: %d", tr("order_info_id"), foundOrder.ID))
	parts = append(parts, fmt.Sprintf("%s: %s", tr("order_info_type"), tr("certificate_type_"+string(foundOrder.CertificateType))))
	parts = append(parts, fmt.Sprintf("%s: %s", tr("order_info_obtain_method"), tr("obtain_method_"+string(foundOrder.ObtainMethod))))
	parts = append(parts, fmt.Sprintf("%s: %s", tr("order_info_status"), tr("status_"+string(foundOrder.ApplicationStatus))))

	if foundOrder.ApplicationStatus == enums.Rejected && foundOrder.RejectionReason != "" {
		parts = append(parts, fmt.Sprintf("%s: %s", tr("order_info_rejection_reason"), foundOrder.RejectionReason))
	}

	return strings.Join(parts, "\n")
}
