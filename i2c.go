package icy

type I2C interface {
	ReadUint8(addr, reg byte) (uint8, error)
	WriteUint8(addr, reg byte, val uint8) error
	ReadUint16(addr, reg byte) (uint16, error)
	WriteUint16(addr, reg byte, val uint16) error
}
