package models

type Cloneable[T any] interface {
	Clone() T
}
