package l2tp

type l2tpAttr struct {
	attrType  uint16
	attrValue []byte
}

func parseMsgAttrs(d []byte) (attrs []l2tpAttr) {
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
	return (a + 3) & ^(4 - 1)
}
