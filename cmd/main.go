package cmd

import (
	"log"

	"github.com/AsmrS4/certificates-plugin/internal/api/handlers"
)

func main() {
	server := &Server{}
	handlers := &handlers.Handler{}

	if err := server.Run("9090", handlers.InitRoutes()); err != nil {
		log.Fatalf("Server start failed... Message: %s", err.Error())
	}
}
