package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"doxxier.tech/doxxier/lib"
)

func ProcessDoxxier() {
	manager := lib.NewDoxxierManager(OnEvent)
	defer manager.Close()
}

func OnEvent(event string) {
	fmt.Println(event)
}

func main() {
	ProcessDoxxier()

	// Set up channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	fmt.Println("Shutting down...")
}
