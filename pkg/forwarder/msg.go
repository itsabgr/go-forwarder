package forwarder

import (
	"net/netip"
)

type Msg struct {
	Data   []byte
	Origin netip.AddrPort
}

func Pack(codec AddrCodec, owner *netip.AddrPort, payload []byte) Msg {
	if owner == nil {
		return Msg{Data: append(make([]byte, 1, 2048), payload...)}
	}
	data := make([]byte, 2048)
	data[0] = byte(codec.Encode(*owner, data[1:]))
	data = data[:1+int(data[0])+copy(data[1+int(data[0]):], payload)]
	return Msg{Data: data}
}
func (msg *Msg) Unpack(codec AddrCodec) (*netip.AddrPort, []byte) {
	var owner []byte
	var payload []byte
	func() {
		defer func() {
			recover()
		}()
		if msg.Data[0] == 0 {
			owner = msg.Data[1 : 1+int(msg.Data[0])]
		} else {
			owner = msg.Data[1 : 1+int(msg.Data[0])]
		}
		payload = msg.Data[1+int(msg.Data[0]):]
	}()
	if len(payload) == 0 {
		return nil, nil
	}
	if len(owner) == 0 {
		return nil, payload
	}
	addrPort := codec.Decode(owner)
	return &addrPort, payload
}
