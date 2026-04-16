package cmd

import (
	"embed"
	"sync"

	"github.com/AsmrS4/certificates-plugin/internal/handler"
	"github.com/AsmrS4/certificates-plugin/internal/persistence"
	"github.com/AsmrS4/certificates-plugin/internal/persistence/impl"
	"github.com/AsmrS4/certificates-plugin/internal/service"
	wasmplugin "github.com/StaZisS/SuperBotGo/sdk/go-plugin"
)

var (
	migrationsFS embed.FS
	i18nFS       embed.FS
	cat          = wasmplugin.NewCatalog("en").
			LoadFS(i18nFS, "i18n")

	once     sync.Once
	cHandler *handler.CertificateHandler
)

func main() {

	wasmplugin.Run(wasmplugin.Plugin{
		ID:      "certificates",
		Name:    "Certificates Plugin",
		Version: "1.0.0",
		Requirements: []wasmplugin.Requirement{
			wasmplugin.Database("Store applications for a certificate").Name("certificate_applications").Build(),
			wasmplugin.Database("Store uploaded certificates by dean").Name("certificates").Build(),
			wasmplugin.File("Store and serve uploaded documents appendix to the certificate").Build(),
			wasmplugin.NotifyReq("Send notificatons").Build(),
		},
		Triggers: []wasmplugin.Trigger{
			orderCertificateCommand(),
			cancelCertificateOrderCommand(),
			findOrderedCertificateByIDCommand(),
			findAllOrderedCertificatesCommand(),
			showCommands(),
		},
	})
}

func initHandler(ctx *wasmplugin.EventContext) *handler.CertificateHandler {
	tr := cat.Tr(ctx.Locale())
	once.Do(func() {

		db, err := persistence.OpenDBConnection("certificate_applications")
		if err != nil {
			ctx.LogError("add: db open: " + err.Error())
			ctx.Reply(wasmplugin.NewMessage(tr("error")))
		}
		defer db.Close()

		appRepo := impl.NewApplicationRepo(db)
		certRepo := impl.NewCertRepo(db)
		certService := service.New(appRepo, certRepo)
		cHandler = handler.NewHandler(certService, cat)
	})
	return cHandler
}

func orderCertificateCommand() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "order_certificate",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Command to start creating application for a certificate",
		Nodes: []wasmplugin.Node{
			wasmplugin.NewStep("type").
				LocalizedText(cat.L("select_certificate_type"), wasmplugin.StylePlain),

			wasmplugin.NewStep("obtain_method").
				LocalizedText(cat.L("select_certificate_obtain"), wasmplugin.StylePlain),
		},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return initHandler(ctx).CreateOrder(ctx)
		},
	}
}

func cancelCertificateOrderCommand() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "cancel_order",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Command to cancel application for a certificate",
		Nodes:       []wasmplugin.Node{},
		Handler: func(ctx *wasmplugin.EventContext) error {
			ctx.Reply(wasmplugin.NewMessage("Привет, мир!"))
			return nil
		},
	}
}

func findOrderedCertificateByIDCommand() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "find_ordered",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Command to find specific ordered certificate by ID",
		Nodes:       []wasmplugin.Node{},
		Handler: func(ctx *wasmplugin.EventContext) error {
			ctx.Reply(wasmplugin.NewMessage("Привет, мир!"))
			return nil
		},
	}
}

func findAllOrderedCertificatesCommand() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "find_all",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Command to find ordered certificates",
		Nodes:       []wasmplugin.Node{},
		Handler: func(ctx *wasmplugin.EventContext) error {
			ctx.Reply(wasmplugin.NewMessage("Привет, мир!"))
			return nil
		},
	}
}

func showCommands() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "list",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Command to show plugin commands",
		Nodes:       []wasmplugin.Node{},
		Handler: func(ctx *wasmplugin.EventContext) error {
			ctx.Reply(wasmplugin.NewMessage("Привет, мир!"))
			return nil
		},
	}
}
