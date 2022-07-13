package forwarder

import (
	"errors"
	"math"
)

const MaxSize = 512

type Packet []byte

func (p Packet) Addr() []byte {
	defer func() {
		recover()
	}()
	return p[1 : 1+p[0]]
}

func (p Packet) Data() []byte {
	defer func() {
		recover()
	}()
	return p[1+p[0]:]
}

func NewPacket(addr, data []byte) Packet {
	if len(addr) > math.MaxUint8 {
		panic(errors.New("too long addr"))
	}
	b := make([]byte, 1, 1+len(addr)+len(data))
	b[0] = byte(len(addr))
	b = append(b, addr...)
	b = append(b, data...)
	return Packet(b)
}
