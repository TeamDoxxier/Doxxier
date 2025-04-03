package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"doxxier.tech/doxxier/lib/network"
)

var callbackChan chan string
var orchestrator network.TransportOrchestrator

func OnEvent(event string) {
	fmt.Println(event)
}

func onMessage() {
	for message := range callbackChan {
		fmt.Println(message)
	}
}

func main() {
	recipient := flag.String("recipient", "", "Recipient contact address")

	flag.Parse()

	orchestrator = network.TransportOrchestrator{}
	orchestrator.Initialise()
	transport := orchestrator.GetTransportation(network.CONNECTION_PRIVACY_HIGH)

	callbackChan = make(chan string, 10)
	contact, err := transport.Initialise()
	if err != nil {
		fmt.Println("Error initializing transport")
		os.Exit(1)
	}

	fmt.Println("Generated contact: ", contact)

	if recipient != nil && *recipient != "" {
		fmt.Println("Connecting to recipient:", *recipient)
		connect(*recipient)
	} else {
		showMenu()
	}
	// Set up channel to receive OS signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	fmt.Println("Shutting down...")
}

func showMenu() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("1. Connect to Recipient")
	fmt.Println("2. Exit")
	fmt.Print("Enter your choice: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	choice = choice[:len(choice)-1] // Remove newline character
	switch choice {
	case "1":
		fmt.Println("Enter recipient contact address:")
		recipient, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		recipient = recipient[:len(recipient)-1] // Remove newline character
		fmt.Println("Connecting to recipient:", recipient)
		// Call send function here
		connect(recipient)
		showConnectedMenu()
	case "2":
		fmt.Println("Exiting...")
		os.Exit(0)
	default:
		fmt.Println("Invalid choice. Please try again.")
	}
}

func showConnectedMenu() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("1. Send Doxxier")
	fmt.Println("2. Send Message")
	fmt.Println("3. Disconnect")
	fmt.Println("4. Exit")
	fmt.Print("Enter your choice: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	choice = choice[:len(choice)-1] // Remove newline character
	switch choice {
	case "1":
		fmt.Println("Enter recipient contact address:")
		recipient, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		recipient = recipient[:len(recipient)-1] // Remove newline character
		fmt.Println("Connecting to recipient:", recipient)
		connect(recipient)
	case "2":
		fmt.Println("Enter message to send:")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}
		message = message[:len(message)-1] // Remove newline character
		fmt.Println("Sending message:", message)
		// Call send function here
		orchestrator.GetTransportation(network.CONNECTION_PRIVACY_HIGH).SendMessage(message)
		fmt.Println("Message sent to recipient.")
	case "3":
		fmt.Println("Disconnecting...")
	case "4":
		fmt.Println("Exiting...")
	default:
		fmt.Println("Invalid choice. Please try again.")
	}
}

func connect(recipient string) {
	transport := orchestrator.GetTransportation(network.CONNECTION_PRIVACY_HIGH)
	callbackChan = make(chan string, 10)
	transport.Connect(callbackChan, recipient)

	go onMessage()
	defer transport.Disconnect()
	fmt.Println("Connected to recipient:", recipient)
	fmt.Println("Listening for messages...")
	showConnectedMenu()
}
