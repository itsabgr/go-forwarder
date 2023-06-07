package forwarder

import (
	"bytes"
	"net/netip"
	"testing"
)

func TestHandler_Default(t *testing.T) {
	payload := []byte{190, 210}
	addr1 := netip.MustParseAddrPort("172.1.1.2:60606")
	addr2 := netip.MustParseAddrPort("10.0.1.1:1140")
	handler := Handler{AddrCodec: RawAddrCodec{}, Default: addr1}
	msg := Pack(handler.AddrCodec, nil, payload)
	msg.Origin = addr2
	err := handler.Handle(&msg)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Origin != addr1 {
		t.FailNow()
	}
	owner, payload2 := msg.Unpack(handler.AddrCodec)
	if !bytes.Equal(payload2, payload) {
		t.FailNow()
	}
	if *owner != addr2 {
		t.FailNow()
	}
}
func TestHandler_Handle(t *testing.T) {
	payload := []byte{190, 210}
	addr1 := netip.MustParseAddrPort("172.1.1.2:60606")
	addr2 := netip.MustParseAddrPort("10.0.1.1:1140")
	handler := Handler{AddrCodec: RawAddrCodec{}}
	msg := Pack(handler.AddrCodec, &addr1, payload)
	msg.Origin = addr2
	err := handler.Handle(&msg)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Origin != addr1 {
		t.FailNow()
	}
	owner, payload2 := msg.Unpack(handler.AddrCodec)
	if !bytes.Equal(payload2, payload) {
		t.FailNow()
	}
	if *owner != addr2 {
		t.FailNow()
	}
}
