package cmixx

import "encoding/binary"

type MessageType int

const (
	doxxier_init MessageType = 0
	doxxier_end              = 1
)

// Marshal returns the byte representation of the [MessageType].
func (mt MessageType) Marshal() [2]byte {
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], uint16(mt))
	return b
}

// UnmarshalMessageType returns the MessageType from its byte representation.
func UnmarshalMessageType(b [2]byte) MessageType {
	return MessageType(binary.LittleEndian.Uint16(b[:]))
}
