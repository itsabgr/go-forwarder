package forwarder

import (
	"errors"
	"net/netip"
)

type Handler struct {
	AddrCodec AddrCodec
	Default   netip.AddrPort
	Debug     bool
}

var ErrInvalidPort = errors.New("invalid port")
var ErrOwnerIsOrigin = errors.New("owner is origin")
var ErrInvalidAddr = errors.New("invalid addr")
var ErrNoPayload = errors.New("no Payload")
var ErrInvalidAddrPort = errors.New("invalid addr or port")

func (h *Handler) Handle(msg *Msg) error {
	owner, payload := msg.Unpack(h.AddrCodec)
	if len(payload) == 0 {
		return ErrNoPayload
	}
	if owner == nil {
		owner = &h.Default
	} else {
		if !owner.IsValid() {
			return ErrInvalidAddrPort
		}
		if owner.Port() == 0 {
			return ErrInvalidPort
		}
		if owner.Addr().IsUnspecified() {
			return ErrInvalidAddr
		}
	}
	if *owner == msg.Origin {
		return ErrOwnerIsOrigin
	}
	*msg = Pack(h.AddrCodec, &msg.Origin, payload)
	msg.Origin = *owner
	return nil
}
