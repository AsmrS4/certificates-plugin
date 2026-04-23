package main

import (
	"embed"
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
	aHandler *handler.CertificateManagementHandler
)

const INTEGER_REGEX = `^\d+$`

func main() {
	wasmplugin.Run(wasmplugin.Plugin{
		ID:      "certificates",
		Name:    "Certificates Plugin",
		Version: "1.2.6",
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
			findRequests(),
			findRequestByID(),
			processCertificateRequest(),
			rejectCertificateRequest(),
		},
	})
}

func initHandlers(ctx *wasmplugin.EventContext) {
	tr := cat.Tr(ctx.Locale())

	once.Do(func() {
		db, err := persistence.OpenDBConnection()
		if err != nil {
			ctx.LogError("add: db open: " + err.Error())
			ctx.Reply(wasmplugin.NewMessage(tr("error")))
		}

		appRepo := impl.NewApplicationRepo(db)
		certRepo := impl.NewCertRepo(db)
		certService := service.New(appRepo, certRepo)
		cmService := service.NewManagementService(appRepo, certRepo)
		cHandler = handler.NewHandler(certService, cat)
		aHandler = handler.NewManagementHandler(certService, cmService, cat)
	})

}

func consumerHandler(ctx *wasmplugin.EventContext) *handler.CertificateHandler {
	initHandlers(ctx)
	if cHandler == nil {
		ctx.LogError("certificate-plugin/main.go: Consumer Handler not initialized.")
		panic("Handler not initialized.")
	}

	return cHandler
}

func deanHandler(ctx *wasmplugin.EventContext) *handler.CertificateManagementHandler {
	initHandlers(ctx)
	if aHandler == nil {
		ctx.LogError("certificate-plugin/main.go: Dean Handler not initialized.")
		panic("Handler not initialized.")
	}

	return aHandler
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
			return consumerHandler(ctx).CreateOrder(ctx)
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
			return consumerHandler(ctx).FindOrderByID(ctx)
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
			return consumerHandler(ctx).CancelOrderByID(ctx)
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
			return consumerHandler(ctx).FindAllActive(ctx)
		},
	}
}

func findRequests() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "Find all certificate orders",
		Type:        wasmplugin.TriggerHTTP,
		Description: "Find all ordered certificate requests from users.",
		Path:        "/api/certificates/all",
		Methods:     []string{"GET"},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return deanHandler(ctx).FindRequests(ctx)
		},
	}
}

func findRequestByID() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "Certificate order details",
		Type:        wasmplugin.TriggerHTTP,
		Description: "Find concrete certificate order details",
		Path:        "/api/certificates",
		Methods:     []string{"GET"},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return deanHandler(ctx).FindRequestDetails(ctx)
		},
	}
}

func processCertificateRequest() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "Start process certificate order",
		Type:        wasmplugin.TriggerHTTP,
		Description: "Start process certificate order",
		Path:        "/api/certificates/process",
		Methods:     []string{"POST"},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return deanHandler(ctx).ProcessRequest(ctx)
		},
	}
}

func rejectCertificateRequest() wasmplugin.Trigger {
	return wasmplugin.Trigger{
		Name:        "Reject certificate order",
		Type:        wasmplugin.TriggerHTTP,
		Description: "Rejection process certificate order",
		Path:        "/api/certificates/reject",
		Methods:     []string{"DELETE"},
		Handler: func(ctx *wasmplugin.EventContext) error {
			return deanHandler(ctx).RejectCertificateRequest(ctx)
		},
	}
}
