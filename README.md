# l2tp

Go implementation of l2tp user space control for Linux.  Similar to iproute2: `ip l2tp add|del tunnel|session`

TODO: 
- ~~Convert to structure similar to [`github.com/jsimonetti/rtnetlink`](https://github.com/jsimonetti/rtnetlink)~~
- Add Docs
- Add Tests
- [PROFIT!!!](https://knowyourmeme.com/memes/profit)


## Examples
### Add Tunnel
```Go
    tunnel := &l2tp.Tunnel{
        EncapType:  l2tp.Encap(l2tp.L2TP_ENCAPTYPE_IP),
        ConnId:     l2tp.ID(6),
        PeerConnId: l2tp.ID(7),
        IpSaddr:    l2tp.IP(net.ParseIP("127.0.0.1")),
        IpDaddr:    l2tp.IP(net.ParseIP("127.0.0.2")),
    }

    l2nl, err := l2tp.Dial(nil)
    if err := l2nl.Tunnel.Add(tunnel); err != nil {
        // handle err
    }
```
### Add Session
```Go
    session := &l2tp.Session{
        PwType:        l2tp.PwType(l2tp.L2TP_PWTYPE_ETH),
        Ifname:        l2tp.Ifname("Iterface03"),
        ConnId:        tunnel.ConnId,
        SessionId:     l2tp.ID(9),
        PeerSessionId: l2tp.ID(10),
    }

    if err := l2nl.Session.Add(session); err != nil {
        // handle err
    }
```
I only implementated the pieces I needed for my project; not feature complete.  PRs welcome.
