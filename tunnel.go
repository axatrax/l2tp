package l2tp

import (
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
)

type TunnelMessage struct {
	command uint8

	EncapType    *uint16
	ProtoVersion *uint8
	ConnId       *uint32
	PeerConnId   *uint32
	IpSaddr      *net.IP
	IpDaddr      *net.IP
	UdpSport     *uint16
	UdpDport     *uint16
	Ip6Saddr     *net.IP
	Ip6Daddr     *net.IP
}

func (t *TunnelMessage) Command() uint8 {
	return t.command
}

func (t *TunnelMessage) MarshalBinary() ([]byte, error) {
	ae := netlink.NewAttributeEncoder()

	if t.EncapType != nil {
		ae.Uint16(L2TP_ATTR_ENCAP_TYPE, *t.EncapType)
	}

	if t.ProtoVersion != nil {
		ae.Uint8(L2TP_ATTR_PROTO_VERSION, *t.ProtoVersion)
	} else {
		// Unlikely any other version is desired
		ae.Uint8(L2TP_ATTR_PROTO_VERSION, 3)
	}

	if t.ConnId != nil {
		ae.Uint32(L2TP_ATTR_CONN_ID, *t.ConnId)
	}

	if t.PeerConnId != nil {
		ae.Uint32(L2TP_ATTR_PEER_CONN_ID, *t.PeerConnId)
	}

	if t.IpSaddr != nil {
		ip4 := t.IpSaddr.To4()
		if ip4 == nil {
			return nil, fmt.Errorf("dst addr (%s) is not an ipv4 address", t.IpSaddr)
		}
		ae.Bytes(L2TP_ATTR_IP_SADDR, ip4)
	}

	if t.IpDaddr != nil {
		ip4 := t.IpSaddr.To4()
		if ip4 == nil {
			return nil, fmt.Errorf("src addr (%s) is not an ipv4 address", t.IpDaddr)
		}
		ae.Bytes(L2TP_ATTR_IP_DADDR, *t.IpDaddr)
	}

	if t.UdpSport != nil {
		ae.Uint16(L2TP_ATTR_UDP_SPORT, *t.UdpSport)
	}

	if t.UdpDport != nil {
		ae.Uint16(L2TP_ATTR_UDP_DPORT, *t.UdpDport)
	}

	if t.Ip6Saddr != nil {
		ae.Bytes(L2TP_ATTR_IP6_SADDR, *t.Ip6Saddr)
	}

	if t.Ip6Daddr != nil {
		ae.Bytes(L2TP_ATTR_IP6_DADDR, *t.Ip6Daddr)
	}

	return ae.Encode()
}

func (t *TunnelMessage) UnmarshalBinary(b []byte) error {
	ad, err := netlink.NewAttributeDecoder(b)
	if err != nil {
		return err
	}

	for ad.Next() {
		switch ad.Type() {
		case L2TP_ATTR_ENCAP_TYPE:
			v := ad.Uint16()
			t.EncapType = &v

		case L2TP_ATTR_PROTO_VERSION:
			v := ad.Uint8()
			t.ProtoVersion = &v

		case L2TP_ATTR_CONN_ID:
			v := ad.Uint32()
			t.ConnId = &v

		case L2TP_ATTR_PEER_CONN_ID:
			v := ad.Uint32()
			t.PeerConnId = &v

		case L2TP_ATTR_IP6_SADDR:
			v := net.IP(ad.Bytes())
			t.Ip6Saddr = &v

		case L2TP_ATTR_IP6_DADDR:
			v := net.IP(ad.Bytes())
			t.Ip6Daddr = &v

		case L2TP_ATTR_IP_SADDR:
			v := net.IP(ad.Bytes())
			t.IpSaddr = &v

		case L2TP_ATTR_IP_DADDR:
			v := net.IP(ad.Bytes())
			t.IpDaddr = &v

		case L2TP_ATTR_UDP_SPORT:
			v := ad.Uint16()
			t.UdpSport = &v

		case L2TP_ATTR_UDP_DPORT:
			v := ad.Uint16()
			t.UdpDport = &v
		}
	}

	return nil
}

type TunnelService struct {
	c *Conn
}

func (t *TunnelService) Add(tun *TunnelMessage) error {
	tun.command = L2TP_CMD_TUNNEL_CREATE

	_, err := t.c.Execute(tun, t.c.genFamily.ID, netlink.Request|netlink.Acknowledge)
	if err != nil {
		return err
	}

	return nil
}

func (t *TunnelService) Delete(tun *TunnelMessage) error {
	tun.command = L2TP_CMD_TUNNEL_DELETE

	_, err := t.c.Execute(tun, t.c.genFamily.ID, netlink.Request|netlink.Acknowledge)
	if err != nil {
		return err
	}

	return nil
}

func (t *TunnelService) List() ([]TunnelMessage, error) {
	req := &TunnelMessage{
		command: L2TP_CMD_TUNNEL_GET,
	}

	resp, err := t.c.Execute(req, t.c.genFamily.ID, netlink.Request|netlink.Dump)
	if err != nil {
		return []TunnelMessage{}, err
	}

	tunnels := make([]TunnelMessage, len(resp))

	for i, t := range resp {
		tunnels[i] = *(t).(*TunnelMessage)
	}

	return tunnels, nil
}

func (t *TunnelService) Get(tun *TunnelMessage) ([]TunnelMessage, error) {
	tunnels, err := t.List()
	if err != nil {
		return nil, err
	}

	result := make([]TunnelMessage, 0, len(tunnels))
	for _, t := range tunnels {
		if tunnelFilterMatch(tun, &t) {
			result = append(result, t)
		}
	}

	return result, nil
}

func tunnelFilterMatch(f, t *TunnelMessage) bool {
	if f.EncapType != nil && *f.EncapType != *t.EncapType {
		return false
	}

	if f.ProtoVersion != nil && *f.ProtoVersion != *t.ProtoVersion {
		return false
	}

	if f.ConnId != nil && *f.ConnId != *t.ConnId {
		return false
	}

	if f.PeerConnId != nil && *f.PeerConnId != *t.PeerConnId {
		return false
	}

	if f.IpSaddr != nil && !f.IpSaddr.Equal(*t.IpSaddr) {
		return false
	}

	if f.IpDaddr != nil && !f.IpDaddr.Equal(*t.IpDaddr) {
		return false
	}

	if f.UdpSport != nil && *f.UdpSport != *t.UdpSport {
		return false
	}

	if f.UdpDport != nil && *f.UdpDport != *t.UdpDport {
		return false
	}

	if f.Ip6Saddr != nil && !f.Ip6Saddr.Equal(*t.Ip6Saddr) {
		return false
	}

	if f.Ip6Daddr != nil && !f.Ip6Daddr.Equal(*t.Ip6Daddr) {
		return false
	}

	return true
}
