package icy

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
)

type (
	CmdCode uint8

	IcydRequest struct {
		Command  CmdCode
		Address  byte
		Register byte
		Value    uint16
	}

	IcydResponse struct {
		Command  CmdCode
		Response uint16
		Error    error
	}

	RemoteI2C struct {
		c   net.Conn
		enc *gob.Encoder
		dec *gob.Decoder
	}

	RemoteClient struct {
		RemoteI2C
	}

	RemoteServer struct {
		Handler I2C
	}
)

var (
	ErrCommunicationFailure = errors.New("communication failure")
)

const (
	CmdReadUint8 CmdCode = iota
	CmdWriteUint8
	CmdReadUint16
	CmdWriteUint16
)

func (rs RemoteServer) Serve(l net.Listener) (err error) {

	for {
		var conn net.Conn
		conn, err = l.Accept()
		if err != nil {
			break
		}
		go func() {
			rs.handleConnection(RemoteI2C{
				c:   conn,
				enc: gob.NewEncoder(conn),
				dec: gob.NewDecoder(conn),
			})
			conn.Close()
		}()
	}

	return

}

func (rs RemoteServer) handleConnection(comm RemoteI2C) (err error) {
	for {
		req := IcydRequest{}
		if err = comm.dec.Decode(&req); err != nil {
			fmt.Printf("[ experienced error ] %s\n", err)
			break
		}
		res := IcydResponse{}
		// fmt.Printf("[ received ] %#v\n", req)
		switch req.Command {
		case CmdReadUint8:
			r, e := rs.Handler.ReadUint8(req.Address, req.Register)
			res.Response, res.Error = uint16(r), e
		case CmdWriteUint8:
			res.Error = rs.Handler.WriteUint8(req.Address, req.Register, uint8(req.Value))
		case CmdReadUint16:
			res.Response, res.Error = rs.Handler.ReadUint16(req.Address, req.Register)
			// fmt.Printf("[ response ] %#v\n", res)
		case CmdWriteUint16:
			res.Error = rs.Handler.WriteUint16(req.Address, req.Register, req.Value)
		}
		if err = comm.enc.Encode(res); err != nil {
			break
		}
	}
	fmt.Printf("[ connection ended ] %s\n", comm.c.RemoteAddr())
	return
}

func NewRemoteClient(target string) (rem *RemoteClient, err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", target)
	if err != nil {
		return
	}

	rem = &RemoteClient{
		RemoteI2C: RemoteI2C{
			c:   conn,
			enc: gob.NewEncoder(conn),
			dec: gob.NewDecoder(conn),
		},
	}
	return
}

func (r RemoteClient) executeRemoteCommand(cmd CmdCode, addr, reg byte, val uint16) (res uint16, err error) {
	if err = r.enc.Encode(IcydRequest{
		Command:  cmd,
		Address:  addr,
		Register: reg,
		Value:    val,
	}); err != nil {
		return
	}
	response := IcydResponse{}
	if err = r.dec.Decode(&response); err != nil {
		return
	}
	res, err = response.Response, response.Error
	return
}

func (r RemoteClient) ReadUint8(addr, reg byte) (res uint8, err error) {
	var r8 uint16
	r8, err = r.executeRemoteCommand(CmdReadUint8, addr, reg, 0)
	res = uint8(r8)
	return
}

func (r RemoteClient) WriteUint8(addr, reg byte, val uint8) (err error) {
	_, err = r.executeRemoteCommand(CmdWriteUint8, addr, reg, uint16(val))
	return
}

func (r RemoteClient) ReadUint16(addr, reg byte) (res uint16, err error) {
	return r.executeRemoteCommand(CmdReadUint16, addr, reg, 0)
}

func (r RemoteClient) WriteUint16(addr, reg byte, val uint16) (err error) {
	_, err = r.executeRemoteCommand(CmdWriteUint16, addr, reg, val)
	return
}

func (r RemoteClient) Close() (err error) {
	err = r.c.Close()
	return
}
