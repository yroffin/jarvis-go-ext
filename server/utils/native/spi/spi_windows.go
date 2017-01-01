package spi

import "os"

// SPIDevice : instance spi device struct
type SPIDevice struct {
	Bus  int      // 0
	Chip int      // 0
	file *os.File // nil

	mode  uint8
	bpw   uint8
	speed uint32
}

// An SPI Device at /dev/spi<bus>.<chip_select>.
func NewSPIDevice(bus int, chipSelect int) *SPIDevice {
	spi := new(SPIDevice)
	spi.Bus = bus
	spi.Chip = chipSelect

	return spi
}

// Open open spi interface
func (spi *SPIDevice) Open() error {
	return nil
}

// Send : Sends bytes over SPI channel and returns []byte response
func (spi *SPIDevice) Send(bytes_to_send []byte) ([]byte, error) {
	ret := []byte{}
	return ret, nil
}

func (spi *SPIDevice) SetMode(mode uint8) error {
	return nil
}

func (spi *SPIDevice) SetBitsPerWord(bpw uint8) error {
	return nil
}

func (spi *SPIDevice) SetSpeed(speed uint32) error {
	return nil
}
