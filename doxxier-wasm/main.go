package main

import (
	"fmt"
	"os"
	"syscall/js"

	"github.com/spf13/cobra"
)

func main() {
	// Set to os.Args because the default is os.Args[1:] and in WASM, args start
	// at 0, not 1.
	wasmCmd.SetArgs(os.Args)

	err := wasmCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var wasmCmd = &cobra.Command{
	Use:     "xxdk-wasm",
	Short:   "WebAssembly bindings for xxDK.",
	Example: "const go = new Go();\ngo.argv = [\"--logLevel=1\"]",
	Run: func(cmd *cobra.Command, args []string) {

		// Enable all top level bindings functions
		println("Registering global functions")
		registerGlobals()

		js.Global().Get("onWasmInitialized").Invoke()

		<-make(chan bool)
		os.Exit(0)
	},
}

func registerGlobals() {
	// Example usage
	println("Initializing Javascript bindings")

	js.Global().Set("returnMessage", js.FuncOf(returnMessage))
	js.Global().Set("createDoxxier", js.FuncOf(CreateDoxxier))
	js.Global().Set("getDoxxier", js.FuncOf(GetDoxxier))
	js.Global().Set("updateDoxxier", js.FuncOf(UpdateDoxxier))
	js.Global().Set("addPart", js.FuncOf(AddPart))
	js.Global().Set("getPart", js.FuncOf(GetPart))
	js.Global().Set("sendMessage", js.FuncOf(SendMessage))
	js.Global().Set("connect", js.FuncOf(Connect))
	println("Javascript bindings initialized")
	<-make(chan bool)
	os.Exit(0)
}
