package wasm

import "syscall/js"

func Initialise(_ js.Value, args []js.Value) interface{} {
	js.Global().Set("doxxier", map[Doxxier]interface{})
}

func AddPart(_ js.Value, args []js.Value) interface{} {

}

func RemovePart(_ js.Value, args []js.Value) interface{} {

}
