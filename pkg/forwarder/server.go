package forwarder

import (
	"log"
	"net"
)

func (h *Handler) Serve(conn *net.UDPConn) (err error) {
	msg := Pack(h.AddrCodec, nil, make([]byte, 1024))
	var n int
	for {
		n, msg.Origin, err = conn.ReadFromUDPAddrPort(msg.Data)
		if err != nil {
			return err
		}
		msg.Data = msg.Data[:n]
		if err := h.Handle(&msg); err != nil {
			if h.Debug {
				log.Println(err)
			}
			continue
		}
		_, err = conn.WriteToUDPAddrPort(msg.Data, msg.Origin)
		if err != nil {
			return err
		}
	}
}
