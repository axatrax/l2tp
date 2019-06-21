package l2tp

import (
	"github.com/mdlayher/netlink"
)

type SessionMessage struct {
	command uint8

	PwType        *uint16
	Ifname        *string
	ConnId        *uint32
	PeerConnId    *uint32
	SessionId     *uint32
	PeerSessionId *uint32
	RecvSeq       *uint8
	SendSeq       *uint8
	LnsMode       *uint8
}

func (s *SessionMessage) Command() uint8 {
	return s.command
}

func (s *SessionMessage) MarshalBinary() ([]byte, error) {
	ae := netlink.NewAttributeEncoder()

	if s.PwType != nil {
		ae.Uint16(L2TP_ATTR_PW_TYPE, *s.PwType)
	}

	if s.ConnId != nil {
		ae.Uint32(L2TP_ATTR_CONN_ID, *s.ConnId)
	}

	if s.PeerConnId != nil {
		ae.Uint32(L2TP_ATTR_PEER_CONN_ID, *s.PeerConnId)
	}

	if s.SessionId != nil {
		ae.Uint32(L2TP_ATTR_SESSION_ID, *s.SessionId)
	}

	if s.PeerSessionId != nil {
		ae.Uint32(L2TP_ATTR_PEER_SESSION_ID, *s.PeerSessionId)
	}

	if s.RecvSeq != nil {
		ae.Uint8(L2TP_ATTR_RECV_SEQ, *s.RecvSeq)
	}

	if s.SendSeq != nil {
		ae.Uint8(L2TP_ATTR_SEND_SEQ, *s.SendSeq)
	}

	if s.LnsMode != nil {
		ae.Uint8(L2TP_ATTR_LNS_MODE, *s.LnsMode)
	}

	if s.Ifname != nil {
		ae.String(L2TP_ATTR_IFNAME, *s.Ifname)
	}

	return ae.Encode()
}

func (s *SessionMessage) UnmarshalBinary(b []byte) error {
	ad, err := netlink.NewAttributeDecoder(b)
	if err != nil {
		return err
	}

	for ad.Next() {
		switch ad.Type() {
		case L2TP_ATTR_PW_TYPE:
			v := ad.Uint16()
			s.PwType = &v

		case L2TP_ATTR_IFNAME:
			v := ad.String()
			s.Ifname = &v

		case L2TP_ATTR_CONN_ID:
			v := ad.Uint32()
			s.ConnId = &v

		case L2TP_ATTR_PEER_CONN_ID:
			v := ad.Uint32()
			s.PeerConnId = &v

		case L2TP_ATTR_SESSION_ID:
			v := ad.Uint32()
			s.SessionId = &v

		case L2TP_ATTR_PEER_SESSION_ID:
			v := ad.Uint32()
			s.PeerSessionId = &v

		case L2TP_ATTR_RECV_SEQ:
			v := ad.Uint8()
			s.RecvSeq = &v

		case L2TP_ATTR_SEND_SEQ:
			v := ad.Uint8()
			s.SendSeq = &v

		case L2TP_ATTR_LNS_MODE:
			v := ad.Uint8()
			s.LnsMode = &v
		}
	}

	return nil
}

type SessionService struct {
	c *Conn
}

func (s *SessionService) Add(sess *SessionMessage) error {
	sess.command = L2TP_CMD_SESSION_CREATE

	_, err := s.c.Execute(sess, s.c.genFamily.ID, netlink.Request|netlink.Acknowledge)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionService) Delete(sess *SessionMessage) error {
	sess.command = L2TP_CMD_SESSION_DELETE

	_, err := s.c.Execute(sess, s.c.genFamily.ID, netlink.Request|netlink.Acknowledge)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionService) List() ([]SessionMessage, error) {
	req := &SessionMessage{
		command: L2TP_CMD_SESSION_GET,
	}

	resp, err := s.c.Execute(req, s.c.genFamily.ID, netlink.Request|netlink.Dump)
	if err != nil {
		return []SessionMessage{}, err
	}

	sessions := make([]SessionMessage, len(resp))

	for i, s := range resp {
		sessions[i] = *(s).(*SessionMessage)
	}

	return sessions, nil
}

func (s *SessionService) Get(sess *SessionMessage) ([]SessionMessage, error) {
	sessions, err := s.List()
	if err != nil {
		return nil, err
	}

	result := make([]SessionMessage, 0, len(sessions))
	for _, s := range sessions {
		if sessionFilterMatch(sess, &s) {
			result = append(result, s)
		}
	}

	return result, nil
}

func sessionFilterMatch(f, s *SessionMessage) bool {
	if f.PwType != nil && *f.PwType != *s.PwType {
		return false
	}

	if f.Ifname != nil && *f.Ifname != *s.Ifname {
		return false
	}

	if f.ConnId != nil && *f.ConnId != *s.ConnId {
		return false
	}

	if f.PeerConnId != nil && *f.PeerConnId != *s.PeerConnId {
		return false
	}

	if f.SessionId != nil && *f.SessionId != *s.SessionId {
		return false
	}

	if f.PeerSessionId != nil && *f.PeerSessionId != *s.PeerSessionId {
		return false
	}

	if f.RecvSeq != nil && *f.RecvSeq != *s.RecvSeq {
		return false
	}

	if f.SendSeq != nil && *f.SendSeq != *s.SendSeq {
		return false
	}

	if f.LnsMode != nil && *f.LnsMode != *s.LnsMode {
		return false
	}

	return true
}
