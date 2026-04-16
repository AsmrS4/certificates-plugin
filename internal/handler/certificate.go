package handler

import (
	"fmt"

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
