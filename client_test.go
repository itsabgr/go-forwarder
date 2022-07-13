package forwarder

import (
	"context"
	"github.com/itsabgr/go-handy"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	t.Parallel()
	for range handy.N(3) {
		testClient(t)
	}
}

func testClient(t *testing.T) {
	bridgeConn := MustListenPacket("udp", "127.0.0.1:0")
	go Serve(bridgeConn)
	defer bridgeConn.Close()
	wait, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	cli, err := ListenCtx(wait, "udp", "127.0.0.1:0", bridgeConn.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()
	_, err = ParseAddr(cli.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	cli2, err := ListenCtx(wait, "udp", "127.0.0.1:0", cli.bridgeAddr.String())
	if err != nil {
		t.Fatal(err)
	}
	defer cli2.Close()
	data := handy.Rand(100)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			_, err = cli2.WriteTo(data, cli.LocalAddr())
			if err != nil {
				return
			}
		}
	}()
	err = cli.SetDeadline(time.Now().Add(time.Second * 2))
	handy.Throw(err)
	data2, fromCli2, err := cli.Recv()
	if err != nil {
		t.Fatal(err)
	}
	handy.Assert(cli2.Addr().String() == fromCli2.String())
	handy.Assert(string(data2) == string(data))
	data = handy.Rand(99)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			_, err = cli.WriteTo(data, fromCli2)
			if err != nil {
				return
			}
		}
	}()
	err = cli2.SetDeadline(time.Now().Add(time.Second * 2))
	handy.Throw(err)
	data2, fromCli, err := cli2.Recv()
	if err != nil {
		t.Fatal(err)
	}
	handy.Assert(cli.Addr().String() == fromCli.String())
	handy.Assert(string(data2) == string(data))
}
