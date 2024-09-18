package wasm

import (
	"syscall/js"

	"doxxier.tech/doxxier/pkg/models"
)

var doxxier *models.Doxxier

// JSGlobal is an interface to abstract the js.Global() function
type JSGlobal interface {
	Get(key string) js.Value
	Set(key string, value js.Value)
}

// DefaultJSGlobal is the default implementation of JSGlobal using js.Global()
type DefaultJSGlobal struct{}

func (d DefaultJSGlobal) Get(key string) js.Value {
	return js.Global().Get(key)
}

func (d DefaultJSGlobal) Set(key string, value js.Value) {
	js.Global().Set(key, value)
}

var global JSGlobal = DefaultJSGlobal{}

func CreateDoxxier(_ js.Value, args []js.Value) any {
	doxxier = models.NewDoxxier()
	json, _ := doxxier.ToJson()
	return js.ValueOf(json)
}

func GetDoxxier(_ js.Value, args []js.Value) any {
	if doxxier == nil {
		return CreateDoxxier(js.Null(), nil)
	}

	json, _ := doxxier.ToJson()
	return js.ValueOf(json)
}

func UpdateDoxxier(_ js.Value, args []js.Value) any {
	if doxxier == nil {
		return js.Null()
	}

	jsDoxxier := args[0]

	doxxier.Description = jsDoxxier.Get("description").String()
	doxxier.Recipient = jsDoxxier.Get("recipient").String()

	json, _ := doxxier.ToJson()
	return js.ValueOf(json)
}

func AddPart(_ js.Value, args []js.Value) any {
	part := models.NewDoxxierPart()
	byteArray := make([]byte, args[0].Length())
	js.CopyBytesToGo(byteArray, args[0])
	part.Content = byteArray
	doxxier.AddPart(*part)
	json, _ := part.ToJson()
	return js.ValueOf(json)
}

func GetPart(_ js.Value, args []js.Value) any {
	id := args[0].String()

	part := doxxier.GetPart(id)
	if part == nil {
		return js.Null()
	}
	json, _ := part.ToJson()
	return js.ValueOf(json)
}
