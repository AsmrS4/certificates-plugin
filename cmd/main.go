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

const INTEGER_REGEX = `^\d+$`

func main() {
	wasmplugin.Run(wasmplugin.Plugin{
		ID:      "certificates",
		Name:    "Certificates Plugin",
		Version: "1.0.0",
		Requirements: []wasmplugin.Requirement{
			wasmplugin.Database("Store applications for a certificate").Name("certificate_applications").Build(),
			wasmplugin.Database("Store uploaded certificates by dean").Name("certificates").Build(),
			wasmplugin.File("Store and serve uploaded documents appendix to the certificate").Build(),
			wasmplugin.NotifyReq("Send notifications").Build(),
		},
		Migrations: wasmplugin.MigrationsFromFS(migrationsFS, "migrations"),
		Triggers: []wasmplugin.Trigger{
			orderCertificateCommand(),
			cancelCertificateOrderCommand(),
			findOrderedCertificateByIDCommand(),
			findAllOrderedCertificatesCommand(),
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
				LocalizedText(cat.L("select_certificate_type"), wasmplugin.StyleHeader).
				LocalizedOptions(cat.L("choose_certificate_type"),
					wasmplugin.Opt("study", "StudyPeriod"),
					wasmplugin.Opt("academ", "Academic"),
					wasmplugin.Opt("recommendation_letter", "Recommendation"),
					wasmplugin.Opt("common", "Common"),
				),

			wasmplugin.NewStep("obtain_method").
				LocalizedText(cat.L("select_certificate_obtain"), wasmplugin.StyleHeader).
				LocalizedOptions(cat.L("choose_obtain_method"),
					wasmplugin.Opt("paper", "Paper"),
					wasmplugin.Opt("electronic", "Electronic"),
				),
		},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return initHandler(ctx).CreateOrder(ctx)
		},
	}
}

// TODO: прикрутить ссылку на скачивание документа
func findOrderedCertificateByIDCommand() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "find_ordered",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Command to find specific ordered certificate by ID",
		Nodes: []wasmplugin.Node{
			wasmplugin.NewStep("enter_id").
				LocalizedText(cat.L("enter_order_id"), wasmplugin.StylePlain).
				Validate(INTEGER_REGEX),
		},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return initHandler(ctx).FindOrderByID(ctx)
		},
	}
}

func cancelCertificateOrderCommand() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "cancel_order",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Command to cancel application for a certificate",
		Nodes: []wasmplugin.Node{
			wasmplugin.NewStep("enter_id").
				LocalizedText(cat.L("enter_order_id"), wasmplugin.StylePlain).
				Validate(INTEGER_REGEX),
		},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return initHandler(ctx).CancelOrderByID(ctx)
		},
	}
}

func findAllOrderedCertificatesCommand() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "find_all",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Command to find ordered certificates",
		Nodes: []wasmplugin.Node{
			wasmplugin.NewStep("status").
				LocalizedText(cat.L("enter_status"), wasmplugin.StyleHeader).
				LocalizedOptions(cat.L("filter_by"),
					wasmplugin.Opt("pending", "Pending"),
					wasmplugin.Opt("prepare", "Prepare"),
					wasmplugin.Opt("done", "Done"),
					wasmplugin.Opt("skip", "Skip"),
				),
		},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return initHandler(ctx).FindAllActive(ctx)
		},
	}
}
