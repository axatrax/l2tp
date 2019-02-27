package l2tp

import "net"

type l2tpAttr struct {
	attrType  uint16
	attrValue []byte
}

/*
func (msg l2tpMessage) toWireFmt() []byte {

	b := []byte{}
	if msg.EncapType != 65535 { // Default value
		b = append(b)

	}

}

*/
func parseAttrs(d []byte) (attrs []l2tpAttr) {
	for len(d) > 0 {
		attrlen := platformEndian.Uint16(d[:2])

		a := l2tpAttr{
			attrType:  platformEndian.Uint16(d[2:4]),
			attrValue: d[4:attrlen],
		}

		attrs = append(attrs, a)
		d = d[alignAttr(int(attrlen)):]
	}

	return
}

func alignAttr(a int) int {
	return (a + 3) & ^3
}

func platformUint8(b []byte) *uint8 {
	v := uint8(b[0])
	return &v
}

func platformUint16(b []byte) *uint16 {
	v := platformEndian.Uint16(b)
	return &v
}

func platformUint32(b []byte) *uint32 {
	v := platformEndian.Uint32(b)
	return &v
}

func platformPutUint8(i uint8) []byte {
	return []byte{i}
}

func platformPutUint16(i uint16) []byte {
	b := make([]byte, 2)
	platformEndian.PutUint16(b, i)
	return b
}

func platformPutUint32(i uint32) []byte {
	b := make([]byte, 4)
	platformEndian.PutUint32(b, i)
	return b
}

func paddedAttr8(attrType uint16, attrValue uint8) (b []byte) {
	b = append(b, platformPutUint16(5)...)        // Length
	b = append(b, platformPutUint16(attrType)...) // Type
	b = append(b, platformPutUint8(attrValue)...) // Value
	b = append(b, []byte{0, 0, 0}...)
	return
}

func paddedAttr16(attrType uint16, attrValue uint16) (b []byte) {
	b = append(b, platformPutUint16(6)...)         // Length
	b = append(b, platformPutUint16(attrType)...)  // Type
	b = append(b, platformPutUint16(attrValue)...) // Value
	b = append(b, []byte{0, 0}...)
	return
}

func paddedAttr32(attrType uint16, attrValue uint32) (b []byte) {
	b = append(b, platformPutUint16(8)...)         // Length
	b = append(b, platformPutUint16(attrType)...)  // Type
	b = append(b, platformPutUint32(attrValue)...) // Value
	return
}

func paddedIP(attrType uint16, ipAddr net.IP) (b []byte) {
	if ip4 := ipAddr.To4(); ip4 != nil {
		b = append(b, platformPutUint16(8)...)        // Length
		b = append(b, platformPutUint16(attrType)...) // Type
		b = append(b, ipAddr[12:]...)                 // Value

	} else if len(ipAddr) == 16 {
		b = append(b, platformPutUint16(20)...)       // Length
		b = append(b, platformPutUint16(attrType)...) // Type
		b = append(b, ipAddr...)                      // Value

	}
	return
}

func paddedAttrString(attrType uint16, attrValue string) []byte {
	b := []byte{}

	value := []byte(attrValue)
	value = append(value, 0) // Null terminating

	msgLen := len(value) + 4 // Plus Header
	withPadding := alignAttr(msgLen)
	pad := withPadding - msgLen

	b = append(b, platformPutUint16(uint16(msgLen))...)
	b = append(b, platformPutUint16(attrType)...)
	b = append(b, value...)

	if withPadding != 0 {
		b = append(b, make([]byte, pad)...)
	}

	return b
}

/*
func paddedAttrBytes(attrType uint16, attrValue []byte) []byte {

}
*/
