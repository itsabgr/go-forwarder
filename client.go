package forwarder

import (
	"context"
	"errors"
	"net"
	"net/netip"
	"runtime"
	"sync/atomic"
	"time"
)

type Client struct {
	conn       net.PacketConn
	id         atomic.Value
	bridgeAddr net.Addr
}

func (c *Client) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	if len(p) == 0 {
		return 0, nil, nil
	}
	for {
		data, from, err := recv(c.conn)
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
		data, from, err := recv(c.conn)
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
	return send(c.conn, target, p)
}

func (c *Client) Close() error {
	return c.conn.Close()
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
	return c.conn.SetDeadline(t)
}

func (c *Client) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *Client) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *Client) Ping() error {
	_, err := send(c.conn, &Addr{bridge: c.bridgeAddr}, []byte{0, 0, 0, 0})
	return err
}
func (c *Client) waitForAddr(ctx context.Context) error {
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer cancel()
		for {
			err := ctx2.Err()
			if err != nil {
				break
			}
			err = c.Ping()
			if err != nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	go func() {
		defer cancel()
		for {
			err := ctx2.Err()
			if err != nil {
				break
			}
			_, _, err = c.Recv()
			if err != nil {
				break
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	for c.Addr() == nil {
		err := ctx2.Err()
		if err != nil {
			return err
		}
		runtime.Gosched()
	}
	return nil
}
func ListenCtx(wait context.Context, network, addr, bridge string) (*Client, error) {
	addrPort, err := netip.ParseAddrPort(bridge)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenPacket(network, addr)
	if err != nil {
		return nil, err
	}
	cli := New(conn, net.UDPAddrFromAddrPort(addrPort))
	err = cli.waitForAddr(wait)
	if err != nil {
		_ = cli.Close()
	}
	return cli, err
}
func New(conn net.PacketConn, bridgeAddr net.Addr) *Client {
	return &Client{conn: conn, bridgeAddr: bridgeAddr}
}
