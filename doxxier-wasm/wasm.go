package main

import (
	"syscall/js"
	"time"

	"doxxier.tech/doxxier/lib"
	"doxxier.tech/doxxier/lib/network"
)

var doxxierManager lib.DoxxierManager
var callbackChan chan string
var orchestrator *network.TransportOrchestrator

// JSGlobal is an interface to abstract the js.Global() function
type JSGlobal interface {
	Get(key string) js.Value
	Set(key string, value js.Value)
}

// DefaultJSGlobal is the default implementation of JSGlobal using js.Global()
type DefaultJSGlobal struct{}

func OnEvent(message string) {
	console := global.Get("console")
	console.Call("log", message)
}

func onMessage() {
	console := global.Get("console")
	for message := range callbackChan {
		console.Call("log", message)
	}
}

func (d DefaultJSGlobal) Get(key string) js.Value {
	return js.Global().Get(key)
}

func (d DefaultJSGlobal) Set(key string, value js.Value) {
	js.Global().Set(key, value)
}

var global JSGlobal = DefaultJSGlobal{}

func returnMessage(_ js.Value, args []js.Value) any {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	//output, _ := fmt.Println("The current time is:", formattedTime)
	return "The current time is: " + formattedTime
}

func CreateDoxxier(_ js.Value, args []js.Value) any {
	doxxierManager = *lib.NewDoxxierManager(OnEvent)
	json, _ := doxxierManager.GetDoxxier().ToJson()
	return js.ValueOf(json)
}

func GetDoxxier(_ js.Value, args []js.Value) any {
	json, _ := doxxierManager.GetDoxxier().ToJson()
	return js.ValueOf(json)
}

func UpdateDoxxier(_ js.Value, args []js.Value) any {
	jsDoxxier := args[0]
	doxxier := doxxierManager.GetDoxxier()
	doxxier.Description = jsDoxxier.Get("description").String()
	doxxier.Recipient = jsDoxxier.Get("recipient").String()

	json, _ := doxxier.ToJson()
	return js.ValueOf(json)
}

func AddPart(_ js.Value, args []js.Value) any {
	part := doxxierManager.AddPart()
	byteArray := make([]byte, args[0].Length())
	js.CopyBytesToGo(byteArray, args[0])
	part.Content = byteArray
	json, _ := part.ToJson()
	return js.ValueOf(json)
}

func GetPart(_ js.Value, args []js.Value) any {
	id := args[0].String()
	doxxier := doxxierManager.GetDoxxier()
	part := doxxier.GetPart(id)
	if part == nil {
		return js.Null()
	}
	json, _ := part.ToJson()
	return js.ValueOf(json)
}

func Connect(_ js.Value, args []js.Value) any {
	if callbackChan == nil {
		callbackChan = make(chan string, 10)
	}
	if orchestrator == nil {
		orchestrator = &network.TransportOrchestrator{}
		orchestrator.Initialise()
	}
	transport := orchestrator.GetTransportation(network.CONNECTION_PRIVACY_HIGH)

	contact, err := transport.Connect(callbackChan)
	if err != nil {
		return js.ValueOf(err.Error())
	}
	//go onMessage()
	return js.ValueOf(contact)
}

func SendMessage(_ js.Value, args []js.Value) any {
	recipient := args[0].String()

	transport := orchestrator.GetTransportation(network.CONNECTION_PRIVACY_HIGH)
	err := transport.SendMessage(recipient, "Hello, world!")
	if err != nil {
		return js.ValueOf(err.Error())
	}

	return js.Null()
}
