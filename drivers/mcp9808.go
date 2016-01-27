package drivers

import (
	"errors"

	"github.com/tylerchr/icy"
)

const (
	MCP9808_ADDR = 0x18

	MCP9808_REG_UPPER_TEMP   = 0x02
	MCP9808_REG_LOWER_TEMP   = 0x03
	MCP9808_REG_CRIT_TEMP    = 0x04
	MCP9808_REG_AMBIENT_TEMP = 0x05
	MCP9808_REG_MANUF_ID     = 0x06
	MCP9808_REG_DEVICE_ID    = 0x07
	MCP9808_REG_RESOLUTION   = 0x08

	MCP9808_MANUF_ID  = 0x0054
	MCP9808_DEVICE_ID = 0x0400
)

const (
	MCP9808_RESOLUTION_LOW uint8 = iota
	MCP9808_RESOLUTION_MEDIUM
	MCP9808_RESOLUTION_HIGH
	MCP9808_RESOLUTION_VERY_HIGH
)

var (
	ErrUnsupportedResolution = errors.New("unsupported resolution")
)

type MCP9808 struct {
	Bus icy.I2C
}

func (mcp MCP9808) Manufacturer() (uint16, error) {
	return mcp.Bus.ReadUint16(MCP9808_ADDR, MCP9808_REG_MANUF_ID)
}

func (mcp MCP9808) Device() (uint16, error) {
	return mcp.Bus.ReadUint16(MCP9808_ADDR, MCP9808_REG_DEVICE_ID)
}

func (mcp MCP9808) Temperature() (temp float64, err error) {
	var word uint16
	if word, err = mcp.Bus.ReadUint16(MCP9808_ADDR, MCP9808_REG_AMBIENT_TEMP); err != nil {
		return
	}

	temp = float64(word&0x0FFF) / 16
	if (word & 0x1000) == 0x1000 {
		temp -= 256
	}
	return
}

func (mcp MCP9808) Resolution() (res uint8, err error) {
	return mcp.Bus.ReadUint8(MCP9808_ADDR, MCP9808_REG_RESOLUTION)
}

func (mcp MCP9808) SetResolution(res uint8) (err error) {
	if res > 0x03 {
		return ErrUnsupportedResolution
	}
	return mcp.Bus.WriteUint8(MCP9808_ADDR, MCP9808_REG_RESOLUTION, res)
}
