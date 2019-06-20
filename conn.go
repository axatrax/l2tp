package l2tp

import (
	"encoding"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

type Conn struct {
	c         conn
	genFamily genetlink.Family

	Tunnel  *TunnelService
	Session *SessionService
}

type conn interface {
	Close() error
	Send(genetlink.Message, uint16, netlink.HeaderFlags) (netlink.Message, error)
	Receive() ([]genetlink.Message, []netlink.Message, error)
	Execute(genetlink.Message, uint16, netlink.HeaderFlags) ([]genetlink.Message, error)
}

func Dial(config *netlink.Config) (*Conn, error) {
	c, err := genetlink.Dial(config)
	if err != nil {
		return nil, err
	}

	family, err := c.GetFamily("l2tp")
	if err != nil {
		return nil, err
	}

	return newConn(c, family), nil
}

func newConn(c conn, f genetlink.Family) *Conn {
	rtc := &Conn{
		c:         c,
		genFamily: f,
	}

	rtc.Tunnel = &TunnelService{c: rtc}
	rtc.Session = &SessionService{c: rtc}

	return rtc
}

func (c *Conn) Close() error {
	return c.c.Close()
}

func (c *Conn) Send(m Message, family uint16, flags netlink.HeaderFlags) (netlink.Message, error) {
	gnlm := genetlink.Message{
		Header: genetlink.Header{
			Command: m.Command(),
			Version: c.genFamily.Version,
		},
	}

	mb, err := m.MarshalBinary()
	if err != nil {
		return netlink.Message{}, err
	}
	gnlm.Data = mb
	reqnm, err := c.c.Send(gnlm, c.genFamily.ID, flags)
	if err != nil {
		return netlink.Message{}, err
	}

	return reqnm, nil
}

func (c *Conn) Receive() ([]Message, []netlink.Message, error) {
	msgs, nlmsgs, err := c.c.Receive()
	if err != nil {
		return nil, nil, err
	}

	genlmsgs, err := unpackMessages(msgs)
	if err != nil {
		return nil, nil, err
	}

	return genlmsgs, nlmsgs, nil
}

func (c *Conn) Execute(m Message, family uint16, flags netlink.HeaderFlags) ([]Message, error) {
	gnlm, err := packMessage(m, c.genFamily.Version, family, flags)
	if err != nil {
		return nil, err
	}

	msgs, err := c.c.Execute(gnlm, c.genFamily.ID, flags)
	if err != nil {
		return nil, err
	}

	return unpackMessages(msgs)
}

type Message interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Command() uint8
	//	rtMessage()
}

func packMessage(m Message, version uint8, family uint16, flags netlink.HeaderFlags) (genetlink.Message, error) {
	gelnm := genetlink.Message{
		Header: genetlink.Header{
			Command: m.Command(),
			Version: version,
		},
	}

	mb, err := m.MarshalBinary()
	if err != nil {
		return genetlink.Message{}, err
	}
	gelnm.Data = mb

	return gelnm, nil
}

func unpackMessages(msgs []genetlink.Message) ([]Message, error) {
	genlmsgs := make([]Message, 0, len(msgs))

	for _, nm := range msgs {
		var m Message
		switch nm.Header.Command {
		case L2TP_CMD_TUNNEL_GET:
			m = &TunnelMessage{command: nm.Header.Command}
		default:
			continue
		}

		if err := (m).UnmarshalBinary(nm.Data); err != nil {
			return nil, err
		}
		genlmsgs = append(genlmsgs, m)
	}

	return genlmsgs, nil
}
