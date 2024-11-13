package main

import (
	"os"
	"syscall/js"

	"doxxier.tech/wasm/internal"
)

func main() {
	jsWriter := internal.JsConsoleWriter{}
	_, _, cleanup := internal.RedirectOutput(jsWriter)
	defer cleanup()

	// Example usage
	println("This will be logged to the JavaScript console as stdout")
	os.Stderr.WriteString("This will be logged to the JavaScript console as stderr\n")

	js.Global().Set("createDoxxier", js.FuncOf(CreateDoxxier))
	js.Global().Set("getDoxxier", js.FuncOf(GetDoxxier))
	js.Global().Set("updateDoxxier", js.FuncOf(UpdateDoxxier))
	js.Global().Set("addPart", js.FuncOf(AddPart))
	js.Global().Set("getPart", js.FuncOf(GetPart))
	<-make(chan bool)
}
