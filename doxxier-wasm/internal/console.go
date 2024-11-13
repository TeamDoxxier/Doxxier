package internal

import (
	"io"
	"os"
	"syscall/js"
)

type JsConsoleWriter struct{}

func (w JsConsoleWriter) Write(p []byte) (n int, err error) {
	js.Global().Get("console").Call("log", string(p))
	return len(p), nil
}

// redirectOutput redirects os.Stdout and os.Stderr to the JavaScript console
func RedirectOutput(w io.Writer) (*os.File, *os.File, func()) {
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
