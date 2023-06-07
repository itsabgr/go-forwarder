package main

import (
	"flag"
	"github.com/itsabgr/go-forwarder/pkg/forwarder"
	"log"
	"net"
	"net/netip"
)

var flagAddr = flag.String("addr", "", "")
var flagDefault = flag.String("default", "", "")
var flagDebug = flag.Bool("debug", false, "")

func main() {
	flag.Parse()
	handler := forwarder.Handler{
		AddrCodec: forwarder.RawAddrCodec{},
		Default:   netip.MustParseAddrPort(*flagDefault),
		Debug:     *flagDebug,
	}
	conn := must(net.ListenUDP("udp", must(net.ResolveUDPAddr("udp", *flagAddr))))
	defer func() { _ = conn.Close() }()
	err := handler.Serve(conn)
	log.Println(err)
}
