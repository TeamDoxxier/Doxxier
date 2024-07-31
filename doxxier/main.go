package main

import (
	"io"
	"os"
	"syscall/js"

	"doxxier.tech/doxxier/wasm"
)

// jsConsoleWriter writes to the JavaScript console
type jsConsoleWriter struct{}

func (w jsConsoleWriter) Write(p []byte) (n int, err error) {
	js.Global().Get("console").Call("log", string(p))
	return len(p), nil
}

// redirectOutput redirects os.Stdout and os.Stderr to the JavaScript console
func redirectOutput(w io.Writer) (*os.File, *os.File, func()) {
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	r, wpipe, _ := os.Pipe()
	rErr, wErrPipe, _ := os.Pipe()

	os.Stdout = wpipe
	os.Stderr = wErrPipe

	done := make(chan struct{})
	go func() {
		io.Copy(w, r)
	}()
	go func() {
		io.Copy(w, rErr)
	}()

	return oldStdout, oldStderr, func() {
		wpipe.Close()
		wErrPipe.Close()
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		close(done)
	}
}

func main() {
	jsWriter := jsConsoleWriter{}
	_, _, cleanup := redirectOutput(jsWriter)
	defer cleanup()

	// Example usage
	println("This will be logged to the JavaScript console as stdout")
	os.Stderr.WriteString("This will be logged to the JavaScript console as stderr\n")
	js.Global().Set("compress", js.FuncOf(wasm.Compress))
	js.Global().Set("transformImage", js.FuncOf(wasm.TransformImage))
	<-make(chan bool)
}
