package l2tp

import (
	"fmt"

	"github.com/mdlayher/genetlink"
	"github.com/mdlayher/netlink"
)

/*
type Tunnel struct{}

func AddTunnel(tunnel *Tunnel) error {

}

func DeleteTunnel(tunnel *Tunnel) error {

}

func GetTunnels() ([]Tunnel, error) {*/
func GetTunnels() error {
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
		panic(err)
	}

	for _, rmsg := range resp {
		fmt.Printf("%+v\n", rmsg.Header)
		fmt.Printf("%+v\n", parseMsgAttrs(rmsg.Data))
	}

	return nil
}
