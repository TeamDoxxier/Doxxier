package testing

import (
	"fmt"
	"reflect"
)

// Value represents a JavaScript value.
type Value struct {
	data interface{}
}

var (
	Undefined = Value{}
	Null      = Value{data: nil}
)

// ValueOf simulates the syscall/js ValueOf function.
func ValueOf(v interface{}) Value {
	return Value{data: v}
}

// Int simulates retrieving an int value.
func (v Value) Int() int {
	if i, ok := v.data.(int); ok {
		return i
	}
	panic("Value is not an int")
}

// String simulates retrieving a string value.
func (v Value) String() string {
	if s, ok := v.data.(string); ok {
		return s
	}
	panic("Value is not a string")
}

// Bool simulates retrieving a boolean value.
func (v Value) Bool() bool {
	if b, ok := v.data.(bool); ok {
		return b
	}
	panic("Value is not a bool")
}

// Call simulates calling a JavaScript method.
func (v Value) Call(method string, args ...Value) Value {
	if obj, ok := v.data.(map[string]interface{}); ok {
		if fn, ok := obj[method].(func(...Value) Value); ok {
			return fn(args...)
		}
		panic(fmt.Sprintf("Method %s not found", method))
	}
	panic("Value is not an object")
}

// Get simulates accessing a property of a JavaScript object.
func (v Value) Get(property string) Value {
	if obj, ok := v.data.(map[string]interface{}); ok {
		if val, exists := obj[property]; exists {
			return Value{data: val}
		}
		panic(fmt.Sprintf("Property %s not found", property))
	}
	panic("Value is not an object")
}

// Set simulates setting a property of a JavaScript object.
func (v Value) Set(property string, value Value) {
	if obj, ok := v.data.(map[string]interface{}); ok {
		obj[property] = value.data
	} else {
		panic("Value is not an object")
	}
}

// InstanceOf simulates checking if a Value is an instance of a constructor.
func (v Value) InstanceOf(constructor Value) bool {
	// For simplicity, we'll assume constructors are represented as strings.
	if constrName, ok := constructor.data.(string); ok {
		if vObj, ok := v.data.(map[string]interface{}); ok {
			if typeName, ok := vObj["_type"].(string); ok {
				return constrName == typeName
			}
		}
	}
	return false
}

// Type simulates the Type method of js.Value.
func (v Value) Type() reflect.Kind {
	return reflect.TypeOf(v.data).Kind()
}
