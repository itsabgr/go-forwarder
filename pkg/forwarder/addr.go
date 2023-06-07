package forwarder

import (
	"crypto/cipher"
	"encoding/binary"
	"net/netip"
)

type AddrCodec interface {
	Encode(addr netip.AddrPort, b []byte) int8
	Decode(b []byte) netip.AddrPort
}
type CipherAddrCodec struct {
	Cipher cipher.Stream
	raw    RawAddrCodec
}

func (c *CipherAddrCodec) Decode(b []byte) netip.AddrPort {
	c.Cipher.XORKeyStream(b, b)
	return c.raw.Decode(b)
}
func (c *CipherAddrCodec) Encode(addr netip.AddrPort, b []byte) int8 {
	n := c.raw.Encode(addr, b)
	if n > 0 {
		c.Cipher.XORKeyStream(b[:n], b[:n])
	}
	return n
}

type RawAddrCodec struct{}

func (RawAddrCodec) Decode(b []byte) netip.AddrPort {
	switch len(b) {
	case 4 + 2, 16 + 2:
	default:
		return netip.AddrPort{}
	}
	port := binary.BigEndian.Uint16(b[:2])
	addr, _ := netip.AddrFromSlice(b[2:])
	return netip.AddrPortFrom(addr, port)
}
func (RawAddrCodec) Encode(addr netip.AddrPort, b []byte) int8 {
	switch {
	case addr.Addr().Is4():
		ip := addr.Addr().As4()
		{
			copy(b[2:], ip[:])
			binary.BigEndian.PutUint16(b[:2], addr.Port())
		}
		return 4 + 2
	case addr.Addr().Is6():
		ip := addr.Addr().As16()
		{
			copy(b[2:], ip[:])
			binary.BigEndian.PutUint16(b[:2], addr.Port())
		}
		return 16 + 2
	default:
		return -1
	}
}
