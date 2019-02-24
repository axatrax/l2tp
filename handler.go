package l2tp

import (
	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

var sockHandle = &handler{}

type handler struct {
	genNetlinkSocket *genetlink.Conn
	l2tpFamilyID     uint16
	_initialized     bool
}

func (h *handler) initHander() (err error) {
	h.genNetlinkSocket, err = genetlink.Dial(nil)
	if err != nil {
		return
	}

	family, err := h.genNetlinkSocket.GetFamily("l2tp")
	if err != nil {
		return
	}

	h.l2tpFamilyID = family.ID
	h._initialized = true
	return
}

func (h *handler) communicateWithKernel(nlmsg *genetlink.Message, flags netlink.HeaderFlags) ([]genetlink.Message, error) {
	if !h._initialized {
		if err := h.initHander(); err != nil {
			return nil, err
		}
	}

	return h.genNetlinkSocket.Execute(*nlmsg, h.l2tpFamilyID, flags)
}
