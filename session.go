package l2tp

import (
	"fmt"
	"os"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
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

func (s Session) toLTV() (b []byte) {
	if s.PwType != nil {
		b = append(b, paddedAttr16(L2TP_ATTR_PW_TYPE, *s.PwType)...)
	}

	if s.ConnId != nil {
		b = append(b, paddedAttr32(L2TP_ATTR_CONN_ID, *s.ConnId)...)
	}

	if s.PeerConnId != nil {
		b = append(b, paddedAttr32(L2TP_ATTR_PEER_CONN_ID, *s.PeerConnId)...)
	}

	if s.SessionId != nil {
		b = append(b, paddedAttr32(L2TP_ATTR_SESSION_ID, *s.SessionId)...)
	}

	if s.PeerSessionId != nil {
		b = append(b, paddedAttr32(L2TP_ATTR_PEER_SESSION_ID, *s.PeerSessionId)...)
	}

	if s.RecvSeq != nil {
		b = append(b, paddedAttr8(L2TP_ATTR_RECV_SEQ, *s.RecvSeq)...)
	}

	if s.SendSeq != nil {
		b = append(b, paddedAttr8(L2TP_ATTR_SEND_SEQ, *s.SendSeq)...)
	}

	if s.LnsMode != nil {
		b = append(b, paddedAttr8(L2TP_ATTR_LNS_MODE, *s.LnsMode)...)
	}

	if s.Ifname != nil {
		b = append(b, paddedAttrString(L2TP_ATTR_IFNAME, *s.Ifname)...)
	}
	return
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

func AddSession(session *Session) error {
	if session.PwType == nil {
		v := uint16(L2TP_PWTYPE_ETH)
		session.PwType = &v
	}

	msg := &genetlink.Message{
		Header: genetlink.Header{
			Command: L2TP_CMD_SESSION_CREATE,
		},
		Data: session.toLTV(),
	}

	_, err := sockHandle.communicateWithKernel(
		msg,
		netlink.HeaderFlagsRequest|netlink.HeaderFlagsAcknowledge,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSession(session *Session) error {
	msg := &genetlink.Message{
		Header: genetlink.Header{
			Command: L2TP_CMD_SESSION_DELETE,
		},
		Data: session.toLTV(),
	}

	_, err := sockHandle.communicateWithKernel(
		msg,
		netlink.HeaderFlagsRequest|netlink.HeaderFlagsAcknowledge,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetSessions() ([]Session, error) {
	var sessions []Session

	msg := &genetlink.Message{
		Header: genetlink.Header{
			Command: L2TP_CMD_SESSION_GET,
		},
	}

	resp, err := sockHandle.communicateWithKernel(
		msg,
		netlink.HeaderFlagsRequest|netlink.HeaderFlagsDump,
	)
	if err != nil {
		return sessions, err
	}

	for _, rmsg := range resp {
		sessions = append(sessions, parsel2tpSession(rmsg.Data))
	}

	return sessions, nil
}
