package forwarder

import (
	"encoding/hex"
	"fmt"
	_ "github.com/mr-tron/base58"
	"net"
	"strings"
)

type Addr struct {
	bridge net.Addr
	id     []byte
}

func (a *Addr) Bridge() net.Addr {
	if a == nil {
		return nil
	}
	return a.bridge
}

func (a *Addr) ID() []byte {
	if a == nil {
		return nil
	}
	return append([]byte{}, a.id...)
}

func ParseAddr(addr string) (*Addr, error) {
	chunks := strings.SplitN(addr, "-", 2)
	if len(chunks) != 2 {
		return nil, fmt.Errorf("invalid addr %s", addr)
	}
	bridge, err := net.ResolveUDPAddr("udp", chunks[0])
	if err != nil {
		return nil, err
	}
	id, err := hex.DecodeString(chunks[1])
	if err != nil {
		return nil, err
	}
	return &Addr{bridge, id}, nil
}

func (a *Addr) String() string {
	if a == nil || len(a.id) == 0 || a.bridge == nil {
		return ""
	}
	return a.bridge.String() + "-" + hex.EncodeToString(a.id)
}

func (a *Addr) Network() string {
	return "forwarder"
}
