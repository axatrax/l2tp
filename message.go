package l2tp

import (
	"fmt"
	"net"
)

type l2tpAttr struct {
	attrType  uint16
	attrValue []byte
}

type l2tpMessage struct {
	//	PwType
	EncapType *uint16
	//	Offset
	//	DataSeq
	//	L2SpecType
	//	L2SpecLen
	ProtoVersion *uint8
	//  Ifname        string
	ConnId     *uint32
	PeerConnId *uint32
	//  SessionId     uint32 // ?
	//  PeerSessionId uint32 // ?
	//	UdpCsum
	//	VlanId
	//	Cookie
	//	PeerCookie
	//  Debug
	//	RecvSeq
	//	SendSeq
	//	LnsMode
	//	UsingIpsec
	//	RecvTimeout
	//	Fd
	//  IpSaddr  net.IP
	//  IpDaddr  net.IP
	//  UdpSport uint16
	//  UdpDport uint16
	//	Mtu
	//	Mru
	//	Stats
	Ip6Saddr *net.IP
	Ip6Daddr *net.IP
	//	UdpZeroCsum6Tx
	//	UdpZeroCsum6Rx
	//	PAD
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

func parsel2tpMsgAttrs(d []byte) (l2tpMsg l2tpMessage) {
	attrs := parseAttrs(d)

	for _, attr := range attrs {
		switch attr.attrType {
		case L2TP_ATTR_ENCAP_TYPE:
			v := platformEndian.Uint16(attr.attrValue)
			l2tpMsg.EncapType = &v
		case L2TP_ATTR_PROTO_VERSION:
			v := uint8(attr.attrValue[0])
			l2tpMsg.ProtoVersion = &v
		case L2TP_ATTR_CONN_ID:
			v := platformEndian.Uint32(attr.attrValue)
			l2tpMsg.ConnId = &v
		case L2TP_ATTR_PEER_CONN_ID:
			v := platformEndian.Uint32(attr.attrValue)
			l2tpMsg.PeerConnId = &v
		case L2TP_ATTR_IP6_SADDR:
			v := net.IP(attr.attrValue)
			l2tpMsg.Ip6Saddr = &v
		case L2TP_ATTR_IP6_DADDR:
			v := net.IP(attr.attrValue)
			l2tpMsg.Ip6Daddr = &v
		default:
			fmt.Printf("Warning: Unknown attr from kernel - %d\n", attr.attrType)
		}
	}
	return
}

func alignAttr(a int) int {
	return (a + 3) & ^3
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
