package l2tp

import (
	"fmt"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

var sockHandle = &handler{}

type handler struct {
	genNetlinkSocket *genetlink.Conn
	l2tpFamilyID     uint16
	version          uint8
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

	if family.Version != 1 {
		return fmt.Errorf("Unsupported family version: %d", family.Version)
	}

	h.l2tpFamilyID = family.ID
	h.version = family.Version
	h._initialized = true
	return
}

func (h *handler) communicateWithKernel(nlmsg *genetlink.Message, flags netlink.HeaderFlags) ([]genetlink.Message, error) {
	if !h._initialized {
		if err := h.initHander(); err != nil {
			return nil, err
		}
	}

	nlmsg.Header.Version = h.version
	return h.genNetlinkSocket.Execute(*nlmsg, h.l2tpFamilyID, flags)
}
