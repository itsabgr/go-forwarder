package forwarder

import (
	"errors"
	"net"
	"sync/atomic"
	"time"
)

type Client struct {
	bridge     net.PacketConn
	id         atomic.Value
	bridgeAddr net.Addr
}

func (c *Client) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	if len(p) == 0 {
		return 0, nil, nil
	}
	for {
		data, from, err := recv(c.bridge)
		if err != nil {
			return 0, nil, err
		}
		if len(from.id) == 0 && len(data) > 0 {
			if from.bridge.String() == c.bridgeAddr.String() {
				c.id.Store(data)
			}
			continue
		}
		return copy(p, data), from, nil
	}
}

func (c *Client) Recv() (b []byte, addr net.Addr, err error) {
	for {
		data, from, err := recv(c.bridge)
		if err != nil {
			return nil, nil, err
		}
		if len(from.id) == 0 && len(data) > 0 {
			if from.bridge.String() == c.bridgeAddr.String() {
				c.id.Store(data)
			}
			continue
		}
		return data, from, nil
	}
}

func (c *Client) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	if addr == nil {
		return 0, errors.New("nil addr")
	}
	if len(p) == 0 {
		return 0, nil
	}
	target := addr.(*Addr)
	return send(c.bridge, target, p)
}

func (c *Client) Close() error {
	return c.bridge.Close()
}

func (c *Client) LocalAddr() net.Addr {
	return c.Addr()
}
func (c *Client) Addr() *Addr {
	id := c.id.Load()
	if id != nil {
		return &Addr{bridge: c.bridgeAddr, id: id.([]byte)}
	}
	return nil
}

func (c *Client) SetDeadline(t time.Time) error {
	return c.bridge.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.bridge.SetReadDeadline(t)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	return c.bridge.SetWriteDeadline(t)
}

func (c *Client) Ping() error {
	_, err := send(c.bridge, &Addr{bridge: c.bridgeAddr}, []byte{0, 0, 0, 0})
	return err
}
func New(bridge net.PacketConn, bridgeAddr net.Addr) *Client {
	return &Client{bridge: bridge, bridgeAddr: bridgeAddr}
}
