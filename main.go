package main

import (
	"embed"
	"fmt"
	"sync"

	"github.com/AsmrS4/certificates-plugin/internal/handler"
	"github.com/AsmrS4/certificates-plugin/internal/persistence"
	"github.com/AsmrS4/certificates-plugin/internal/persistence/impl"
	"github.com/AsmrS4/certificates-plugin/internal/service"
	wasmplugin "github.com/StaZisS/SuperBotGo/sdk/go-plugin"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

//go:embed i18n/*.toml
var i18nFS embed.FS

var (
	cat = wasmplugin.NewCatalog("en").
		LoadFS(i18nFS, "i18n")

	once     sync.Once
	cHandler *handler.CertificateHandler
)

const INTEGER_REGEX = `^\d+$`

func main() {
	wasmplugin.Run(wasmplugin.Plugin{
		ID:      "certificates",
		Name:    "Certificates Plugin",
		Version: "1.1.4",
		Requirements: []wasmplugin.Requirement{
			wasmplugin.Database("Store applications for a certificate").Build(),
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
		db, err := persistence.OpenDBConnection()
		if err != nil {
			ctx.LogError("add: db open: " + err.Error())
			ctx.Reply(wasmplugin.NewMessage(tr("error")))
		}
		var exists bool
		row := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'certificate_applications')")
		if err := row.Scan(&exists); err != nil {
			ctx.LogError("check table: " + err.Error())
		} else {
			ctx.Log(fmt.Sprintf("table certificate_applications exists: %v", exists))
		}

		appRepo := impl.NewApplicationRepo(db)
		certRepo := impl.NewCertRepo(db)
		certService := service.New(appRepo, certRepo)
		cHandler = handler.NewHandler(certService, cat)
	})

	if cHandler == nil {
		panic("Handler not initialized.")
	}

	return cHandler
}

func orderCertificateCommand() wasmplugin.Trigger {

	return wasmplugin.Trigger{
		Name:        "order_certificate",
		Type:        wasmplugin.TriggerMessenger,
		Description: "Order certificate",
		Nodes: []wasmplugin.Node{

			wasmplugin.NewStep("type").
				LocalizedText(cat.L("select_certificate_type"), wasmplugin.StylePlain).
				DynamicOptions("",
					func(cbCtx *wasmplugin.CallbackContext) []wasmplugin.Option {
						return []wasmplugin.Option{
							wasmplugin.Opt(cat.L("study")[cbCtx.Locale], "StudyPeriod"),
							wasmplugin.Opt(cat.L("academy")[cbCtx.Locale], "Academic"),
							wasmplugin.Opt(cat.L("recommendation_letter")[cbCtx.Locale], "Recommendation"),
							wasmplugin.Opt(cat.L("common")[cbCtx.Locale], "Common"),
						}
					},
				),

			wasmplugin.NewStep("obtain_method").
				LocalizedText(cat.L("select_certificate_obtain"), wasmplugin.StylePlain).
				DynamicOptions("",
					func(cbCtx *wasmplugin.CallbackContext) []wasmplugin.Option {
						return []wasmplugin.Option{
							wasmplugin.Opt(cat.L("paper")[cbCtx.Locale], "Paper"),
							wasmplugin.Opt(cat.L("electronic")[cbCtx.Locale], "Electronic"),
						}
					},
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
		Description: "Find ordered certificate details",
		Nodes: []wasmplugin.Node{
			wasmplugin.NewStep("enter_id").
				LocalizedText(cat.L("enter_order_id"), wasmplugin.StyleHeader).
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
		Description: "Cancel certificate order",
		Nodes: []wasmplugin.Node{
			wasmplugin.NewStep("enter_id").
				LocalizedText(cat.L("enter_order_id"), wasmplugin.StyleHeader).
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
		Description: "Find ordered certificates",
		Nodes: []wasmplugin.Node{
			wasmplugin.NewStep("status").
				LocalizedText(cat.L("filter_by"), wasmplugin.StyleHeader).
				DynamicOptions("",
					func(cbCtx *wasmplugin.CallbackContext) []wasmplugin.Option {
						return []wasmplugin.Option{
							wasmplugin.Opt(cat.L("pending")[cbCtx.Locale], "Pending"),
							wasmplugin.Opt(cat.L("prepare")[cbCtx.Locale], "Prepare"),
							wasmplugin.Opt(cat.L("done")[cbCtx.Locale], "Done"),
							wasmplugin.Opt(cat.L("skip")[cbCtx.Locale], "Skip"),
						}
					},
				),
		},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return initHandler(ctx).FindAllActive(ctx)
		},
	}
}
