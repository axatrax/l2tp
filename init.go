package l2tp

import (
	"encoding/binary"
	"unsafe"
)

var platformEndian binary.ByteOrder

func init() {
	var i uint32 = 0xAABB
	if *(*byte)(unsafe.Pointer(&i)) == 0xAA {
		platformEndian = binary.BigEndian
	} else {
		platformEndian = binary.LittleEndian
	}
}
