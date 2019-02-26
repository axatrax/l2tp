package l2tp

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

/*
func paddedAttr8(attrType uint16, attrValue uint8) []byte {

}

func paddedAttr16(attrType uint16, attrValue uint16) []byte {

}

func paddedAttr32(attrType uint16, attrValue uint32) []byte {

}

func paddedAttrString(attrType uint16, attrValue string) []byte {

}

func paddedAttrBytes(attrType uint16, attrValue []byte) []byte {

}
*/
