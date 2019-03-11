package l2tp

import "net"

func typeUint8(i uint) *uint8 {
	v := uint8(i)
	return &v
}

func typeUint16(i uint) *uint16 {
	v := uint16(i)
	return &v
}

func typeUint32(i uint) *uint32 {
	v := uint32(i)
	return &v
}

func typeIP(ip net.IP) *net.IP {
	return &ip
}

func typeString(s string) *string {
	return &s
}

var Port = typeUint16
var ID = typeUint32
var PwType = typeUint16
var RecvSeq = typeUint8
var SendSeq = typeUint8
var LnsMode = typeUint8
var Encap = typeUint16
var Version = typeUint8
var IP = typeIP
var Ifname = typeString
