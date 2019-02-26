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

/*

func AddTunnel(tunnel *Tunnel) error {

}

func DeleteTunnel(tunnel *Tunnel) error {

}

*/
func GetTunnels() ([]Tunnel, error) {
	var tunnels []Tunnel

	msg := &genetlink.Message{
		Header: genetlink.Header{
			Command: L2TP_CMD_TUNNEL_GET,
		},
	}

	resp, err := sockHandle.communicateWithKernel(
		msg,
		netlink.HeaderFlagsRequest|netlink.HeaderFlagsDump,
	)
	if err != nil {
		return tunnels, err
	}

	for _, rmsg := range resp {
		tunnels = append(tunnels, parsel2tpTunnel(rmsg.Data))
	}

	return tunnels, nil
}
