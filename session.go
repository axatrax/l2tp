package l2tp

import (
	"fmt"
	"os"
)

type Session struct {
	PwType *uint16
	//  EncapType *uint16
	//	Offset
	//	DataSeq
	//	L2SpecType
	//	L2SpecLen
	//  ProtoVersion *uint8
	Ifname        *string
	ConnId        *uint32
	PeerConnId    *uint32
	SessionId     *uint32
	PeerSessionId *uint32
	//	UdpCsum
	//	VlanId
	//	Cookie
	//	PeerCookie
	//  Debug
	RecvSeq *uint8
	SendSeq *uint8
	LnsMode *uint8
	//	UsingIpsec
	//	RecvTimeout
	//	Fd
	//  IpSaddr  *net.IP
	//  IpDaddr  *net.IP
	//  UdpSport *uint16
	//  UdpDport *uint16
	//	Mtu
	//	Mru
	//	Stats
	//  Ip6Saddr *net.IP
	//  Ip6Daddr *net.IP
	//	UdpZeroCsum6Tx
	//	UdpZeroCsum6Rx
	//	PAD
}

func parsel2tpSession(d []byte) (session Session) {
	attrs := parseAttrs(d)

	if os.Getenv("DEBUG") != "" {
		fmt.Println("Parsed LTVs: ")
		fmt.Println(attrs)
	}

	for _, attr := range attrs {
		switch attr.attrType {

		case L2TP_ATTR_PW_TYPE:
			session.PwType = platformUint16(attr.attrValue)

		case L2TP_ATTR_IFNAME:
			v := string(attr.attrValue)
			session.Ifname = &v

		case L2TP_ATTR_CONN_ID:
			session.ConnId = platformUint32(attr.attrValue)

		case L2TP_ATTR_PEER_CONN_ID:
			session.PeerConnId = platformUint32(attr.attrValue)

		case L2TP_ATTR_SESSION_ID:
			session.SessionId = platformUint32(attr.attrValue)

		case L2TP_ATTR_PEER_SESSION_ID:
			session.PeerSessionId = platformUint32(attr.attrValue)

		case L2TP_ATTR_RECV_SEQ:
			session.RecvSeq = platformUint8(attr.attrValue)

		case L2TP_ATTR_SEND_SEQ:
			session.SendSeq = platformUint8(attr.attrValue)

		case L2TP_ATTR_LNS_MODE:
			session.LnsMode = platformUint8(attr.attrValue)

		default:
			if os.Getenv("DEBUG") != "" {
				fmt.Printf("Warning: Unknown attr from kernel - %d\n", attr.attrType)
			}
		}
	}
	return
}

/*
func AddSession(session *Session, tunnel *Tunnel) error {

}

func DeleteSession(session *Session, tunnel *Tunnel) error {

}

func GetSessions() ([]Session, error) {

}
*/
