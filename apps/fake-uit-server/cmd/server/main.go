package main

import (
	"log"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/bootstrap"
)

func main() {
	router, err := bootstrap.Init()
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := router.Run(":3000"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
