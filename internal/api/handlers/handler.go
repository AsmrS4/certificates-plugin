package handlers

import "github.com/gin-gonic/gin"

type Handler struct{}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	//Эндпоинты для ядра системы
	plugin := router.Group("/plugin")
	{
		certificates := plugin.Group("/certificates")
		{
			certificates.POST("/", h.createCertificateApplication)
			certificates.GET("/:certificate_id", h.getCertificateApplicationDetails)
			certificates.GET("/:certificate_id/status", h.getCertificateApplicationStatus)
			certificates.DELETE("/:certificate_id/cancel", h.cancelCertificateApplication)
		}

	}

	return router
}
