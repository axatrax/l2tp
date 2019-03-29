package l2tp

import (
	"fmt"
	"net"
	"os"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

type Tunnel struct {
	//  PwType    *uint16
	EncapType *uint16
	//	Offset
	//	DataSeq
	//	L2SpecType
	//	L2SpecLen
	ProtoVersion *uint8
	//  Ifname       *string
	ConnId     *uint32
	PeerConnId *uint32
	//  SessionId     *uint32
	// 	PeerSessionId *uint32
	//	UdpCsum
	//	VlanId
	//	Cookie
	//	PeerCookie
	//  Debug
	//  RecvSeq *uint8
	//  SendSeq *uint8
	//  LnsMode *uint8
	//	UsingIpsec
	//	RecvTimeout
	//	Fd
	IpSaddr  *net.IP
	IpDaddr  *net.IP
	UdpSport *uint16
	UdpDport *uint16
	//	Mtu
	//	Mru
	//	Stats
	Ip6Saddr *net.IP
	Ip6Daddr *net.IP
	//	UdpZeroCsum6Tx
	//	UdpZeroCsum6Rx
	//	PAD
}

func (t Tunnel) toLTV() (b []byte) {
	if t.EncapType != nil {
		b = append(b, paddedAttr16(L2TP_ATTR_ENCAP_TYPE, *t.EncapType)...)
	}

	if t.ProtoVersion != nil {
		b = append(b, paddedAttr8(L2TP_ATTR_PROTO_VERSION, *t.ProtoVersion)...)
	}

	if t.ConnId != nil {
		b = append(b, paddedAttr32(L2TP_ATTR_CONN_ID, *t.ConnId)...)
	}

	if t.PeerConnId != nil {
		b = append(b, paddedAttr32(L2TP_ATTR_PEER_CONN_ID, *t.PeerConnId)...)
	}

	if t.IpSaddr != nil {
		b = append(b, paddedIP(L2TP_ATTR_IP_SADDR, *t.IpSaddr)...)
	}

	if t.IpDaddr != nil {
		b = append(b, paddedIP(L2TP_ATTR_IP_DADDR, *t.IpDaddr)...)
	}

	if t.UdpSport != nil {
		b = append(b, paddedAttr16(L2TP_ATTR_UDP_SPORT, *t.UdpSport)...)
	}

	if t.UdpDport != nil {
		b = append(b, paddedAttr16(L2TP_ATTR_UDP_DPORT, *t.UdpDport)...)
	}

	if t.Ip6Saddr != nil {
		b = append(b, paddedIP(L2TP_ATTR_IP6_SADDR, *t.Ip6Saddr)...)
	}

	if t.Ip6Daddr != nil {
		b = append(b, paddedIP(L2TP_ATTR_IP6_DADDR, *t.Ip6Daddr)...)
	}

	return
}

func parsel2tpTunnel(d []byte) (tunnel Tunnel) {
	attrs := parseAttrs(d)

	if os.Getenv("DEBUG") != "" {
		fmt.Println("Parsed LTVs: ")
		fmt.Println(attrs)
	}

	for _, attr := range attrs {
		switch attr.attrType {

		case L2TP_ATTR_ENCAP_TYPE:
			tunnel.EncapType = platformUint16(attr.attrValue)

		case L2TP_ATTR_PROTO_VERSION:
			tunnel.ProtoVersion = platformUint8(attr.attrValue)

		case L2TP_ATTR_CONN_ID:
			tunnel.ConnId = platformUint32(attr.attrValue)

		case L2TP_ATTR_PEER_CONN_ID:
			tunnel.PeerConnId = platformUint32(attr.attrValue)

		case L2TP_ATTR_IP6_SADDR:
			v := net.IP(attr.attrValue)
			tunnel.Ip6Saddr = &v

		case L2TP_ATTR_IP6_DADDR:
			v := net.IP(attr.attrValue)
			tunnel.Ip6Daddr = &v

		case L2TP_ATTR_IP_SADDR:
			v := net.IP(attr.attrValue)
			tunnel.IpSaddr = &v

		case L2TP_ATTR_IP_DADDR:
			v := net.IP(attr.attrValue)
			tunnel.IpDaddr = &v

		default:
			if os.Getenv("DEBUG") != "" {
				fmt.Printf("Warning: Unknown attr from kernel - %d\n", attr.attrType)
			}
		}
	}
	return
}

func AddTunnel(tunnel *Tunnel) error {
	if tunnel.ProtoVersion == nil {
		v := uint8(3)
		tunnel.ProtoVersion = &v
	}

	msg := &genetlink.Message{
		Header: genetlink.Header{
			Command: L2TP_CMD_TUNNEL_CREATE,
		},
		Data: tunnel.toLTV(),
	}

	_, err := sockHandle.communicateWithKernel(
		msg,
		netlink.Request|netlink.Acknowledge,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTunnel(tunnel *Tunnel) error {
	msg := &genetlink.Message{
		Header: genetlink.Header{
			Command: L2TP_CMD_TUNNEL_DELETE,
		},
		Data: tunnel.toLTV(),
	}

	_, err := sockHandle.communicateWithKernel(
		msg,
		netlink.Request|netlink.Acknowledge,
	)
	if err != nil {
		return err
	}

	return nil
}

func GetTunnels() ([]Tunnel, error) {
	var tunnels []Tunnel

	msg := &genetlink.Message{
		Header: genetlink.Header{
			Command: L2TP_CMD_TUNNEL_GET,
		},
	}

	resp, err := sockHandle.communicateWithKernel(
		msg,
		netlink.Request|netlink.Dump,
	)
	if err != nil {
		return tunnels, err
	}

	for _, rmsg := range resp {
		tunnels = append(tunnels, parsel2tpTunnel(rmsg.Data))
	}

	return tunnels, nil
}
