# l2tp

Go implementation of l2tp user space control for Linux.  Similar to iproute2: `ip l2tp add|del tunnel|session`


### Add Tunnel

    tunnel := &l2tp.Tunnel{
        EncapType:  l2tp.Encap(l2tp.L2TP_ENCAPTYPE_IP)
        ConnId:     l2tp.ID(6),
        PeerConnId: l2tp.ID(7),
        IpSaddr:    l2tp.IP(net.ParseIP("127.0.0.1")),
        IpDaddr:    l2tp.IP(net.ParseIP("127.0.0.2")),
    }

    l2tp.AddTunnel(tunnel)

### Add Session

    session := &l2tp.Session{
        Ifname:        l2tp.Ifname("Iterface03"),
        ConnId:        l2tp.ID(6),
        SessionId:     l2tp.ID(9),
        PeerSessionId: l2tp.ID(10),
    }

    l2tp.AddSession(session)

I only implementated the pieces I needed for my project; not feature complete!