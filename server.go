package forwarder

import (
	"github.com/itsabgr/go-handy"
	"net"
	"net/netip"
)

func Serve(conn net.PacketConn) error {
	b := make([]byte, MaxSize)
	for {
		n, from, err := conn.ReadFrom(b)
		if err != nil {
			return err
		}
		pack := Packet(b[:n])
		if len(pack.Data()) == 0 {
			continue
		}
		if len(pack.Addr()) == 0 {
			addrPort, _ := netip.MustParseAddrPort(from.String()).MarshalBinary()
			_, err := conn.WriteTo(NewPacket(nil, addrPort), from)
			handy.Throw(err)
		} else {
			addrPort := netip.AddrPort{}
			err := addrPort.UnmarshalBinary(pack.Addr())
			handy.Throw(err)
			sender, _ := netip.MustParseAddrPort(from.String()).MarshalBinary()
			_, err = conn.WriteTo(NewPacket(sender, pack.Data()), net.UDPAddrFromAddrPort(addrPort))
			handy.Throw(err)
		}
	}
}
