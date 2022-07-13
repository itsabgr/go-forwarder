package forwarder

import (
	"github.com/itsabgr/go-handy"
	"net"
)

func MustListenPacket(network, addr string) net.PacketConn {
	conn, err := net.ListenPacket(network, addr)
	handy.Throw(err)
	return conn
}
