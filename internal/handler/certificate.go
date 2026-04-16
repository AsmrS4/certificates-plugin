package handler

import (
	"fmt"
	"strconv"

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
	id := ctx.Param("id")

	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.Reply(wasmplugin.NewMessage("Неверный формат ID. Введите число."))
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
		ctx.Reply(wasmplugin.NewMessage(tr("order_not_found")))
		return nil
	}

	localizedMessage := h.foundMessageToString(order, tr)

	ctx.Reply(wasmplugin.NewMessage(localizedMessage))
	ctx.Log(fmt.Sprintf("get: found certificate order #%d  %d", order.ID, ctx.Messenger.UserID))

	return nil
}

func (h *CertificateHandler) foundMessageToString(foundOrder *models.CertificateApplication, trans func(key string, args ...any) string) string {
	tr := trans
	typeName := tr("certificate_type_" + string(foundOrder.CertificateType))
	methodName := tr("obtain_method_" + string(foundOrder.ObtainMethod))
	statusName := tr("status_" + string(foundOrder.ApplicationStatus))

	return fmt.Sprintf(
		"%s: %d\n%s: %s\n%s: %s\n%s: %s",
		tr("order_info_id"), foundOrder.ID,
		tr("order_info_type"), typeName,
		tr("order_info_obtain_method"), methodName,
		tr("order_info_status"), statusName,
	)
}
