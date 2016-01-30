package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/tylerchr/icy"
)

type (
	CmdCode uint8

	IcydRequest struct {
		Command           CmdCode
		Address, Register byte
	}
)

const (
	CmdReadUint8 CmdCode = iota
	CmdWriteUint8
	CmdReadUint16
	CmdWriteUint16
)

var port = flag.Int("port", 6000, "on which port to listen for clients")

func main() {

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	// dummy := DummyServer{}
	handler, err := icy.NewSMBus(0x01)
	if err != nil {
		panic(err)
	}
	server := icy.RemoteServer{
		Handler: handler,
	}
	err = server.Serve(listener)
	fmt.Printf("[ server ended ] %s\n", err)

}

// type DummyServer struct{}

// func (ds DummyServer) ReadUint8(addr, reg byte) (uint8, error) {
// 	return 0xAB, nil
// }

// func (ds DummyServer) WriteUint8(addr, reg byte, val uint8) error {
// 	return nil
// }

// func (ds DummyServer) ReadUint16(addr, reg byte) (uint16, error) {
// 	return 0xABCD, nil
// }

// func (ds DummyServer) WriteUint16(addr, reg byte, val uint16) error {
// 	return nil
// }
