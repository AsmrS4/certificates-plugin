package main

import (
	"log"

	certificatesplugin "github.com/AsmrS4/certificates-plugin"
)

func main() {
	server := new(certificatesplugin.Server)
	if err := server.Run("9090"); err != nil {
		log.Fatalf("Server start failed... Message: %s", err.Error())
	}
}
