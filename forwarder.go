package forwarder

import "net"

type IWriteTo interface {
	WriteTo(p []byte, addr net.Addr) (n int, err error)
}

type IReadFrom interface {
	ReadFrom(p []byte) (n int, addr net.Addr, err error)
}

func recv(conn IReadFrom) ([]byte, *Addr, error) {
	b := make([]byte, MaxSize)
	for {
		n, bridge, err := conn.ReadFrom(b)
		if err != nil {
			return nil, nil, err
		}
		pack := Packet(b[:n])
		sender := pack.Addr()
		data := pack.Data()
		if len(data) == 0 {
			continue
		}
		addr := &Addr{
			bridge: bridge,
			id:     sender,
		}
		return data, addr, nil
	}
}
func send(conn IWriteTo,
	peer *Addr,
	data []byte,
) (int, error) {
	pack := NewPacket(peer.id, data)
	return conn.WriteTo(pack, peer.bridge)
}
