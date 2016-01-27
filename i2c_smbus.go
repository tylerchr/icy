package icy

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type (
	// one of the constants from v3.12:include/uapi/linux/i2c.h
	SMBusTransactionType uint32

	// from v3.12:include/uapi/linux/i2c-dev.h
	i2c_smbus_ioctl_data struct {
		readWrite byte
		command   byte
		size      SMBusTransactionType
		data      uintptr
	}

	SMBus struct {
		handle *os.File
		addr   byte
	}
)

const (
	// from v3.12:include/uapi/linux/i2c.h
	I2C_SMBUS_READ  = 1
	I2C_SMBUS_WRITE = 0

	// from v3.12:include/uapi/linux/i2c-dev.h
	I2C_SLAVE = 0x0703
	I2C_FUNCS = 0x0705
	I2C_SMBUS = 0x0720
)

const (
	// from v3.12:include/uapi/linux/i2c.h
	I2C_SMBUS_QUICK SMBusTransactionType = iota
	I2C_SMBUS_BYTE
	I2C_SMBUS_BYTE_DATA
	I2C_SMBUS_WORD_DATA
	I2C_SMBUS_PROC_CALL
	I2C_SMBUS_BLOCK_DATA
	I2C_SMBUS_I2C_BLOCK_BROKEN
	I2C_SMBUS_BLOCK_PROC_CALL
	I2C_SMBUS_I2C_BLOCK_DATA
)

func NewSMBus(bus byte) (*SMBus, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil, err
	}
	return &SMBus{handle: f}, nil
}

func (smbus *SMBus) setSlaveAddress(addr byte) error {
	if smbus.addr != addr {
		if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, smbus.handle.Fd(), I2C_SLAVE, uintptr(addr)); err != 0 {
			return err
		}
		smbus.addr = addr
	}
	return nil
}

func (smbus *SMBus) i2c_smbus_access(read_write, command byte, size SMBusTransactionType, data []byte) (err error) {

	args := i2c_smbus_ioctl_data{
		readWrite: read_write,
		command:   command,
		size:      size,
		data:      uintptr(unsafe.Pointer(&data[0])),
	}

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, smbus.handle.Fd(), I2C_SMBUS, uintptr(unsafe.Pointer(&args))); errno != 0 {
		err = syscall.Errno(errno)
	}
	return
}

func (smbus *SMBus) ReadUint8(addr, reg byte) (data uint8, err error) {
	if err = smbus.setSlaveAddress(addr); err == nil {
		blockData := []byte{0, 0}
		err = smbus.i2c_smbus_access(I2C_SMBUS_READ, reg, I2C_SMBUS_BYTE_DATA, blockData)
		data = uint8(blockData[0])
	}
	return
}

func (smbus *SMBus) WriteUint8(addr, reg byte, val uint8) (err error) {
	if err = smbus.setSlaveAddress(addr); err == nil {
		blockData := []byte{val, 0}
		err = smbus.i2c_smbus_access(I2C_SMBUS_WRITE, reg, I2C_SMBUS_BYTE_DATA, blockData)
	}
	return
}

func (smbus *SMBus) ReadUint16(addr, reg byte) (data uint16, err error) {
	if err = smbus.setSlaveAddress(addr); err == nil {
		blockData := []byte{0, 0}
		err = smbus.i2c_smbus_access(I2C_SMBUS_READ, reg, I2C_SMBUS_WORD_DATA, blockData)
		data = (uint16(blockData[0]) << 8) | uint16(blockData[1])
	}
	return
}

func (smbus *SMBus) WriteUint16(addr, reg byte, val uint16) (err error) {
	if err = smbus.setSlaveAddress(addr); err == nil {
		blockData := []byte{byte((val >> 8) & 0x00FF), byte(val & 0x00FF)}
		err = smbus.i2c_smbus_access(I2C_SMBUS_WRITE, reg, I2C_SMBUS_WORD_DATA, blockData)
	}
	return
}
