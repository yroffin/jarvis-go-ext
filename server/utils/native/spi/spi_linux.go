/**
 * Copyright 2017 Yannick Roffin
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package spi

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/Sirupsen/logrus"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/ioctl"
)

/*
	Original source here
	Cf. https://github.com/luismesas/goPi/blob/master/spi/SPIDevice.go
*/

const (
	// SPIDEV : default spi path
	SPIDEV            = "/dev/spidev"
	SPI_HARDWARE_ADDR = 0
	SPI_BUS           = 0
	SPI_CHIP          = 0
	SPI_DELAY         = 0
)

// SpiIocTransfert : spi struct
type SpiIocTransfert struct {
	txBuf       uint64
	rxBuf       uint64
	length      uint32
	speedHz     uint32
	delayUsecs  uint16
	bitsPerWord uint8
	csChange    uint8
	pad         uint32
}

// SPIDevice : instance spi device struct
type SPIDevice struct {
	Bus  int      // 0
	Chip int      // 0
	file *os.File // nil

	mode  uint8
	bpw   uint8
	speed uint32
}

const SPI_IOC_MAGIC = 107

// Read of SPI mode (SPI_MODE_0..SPI_MODE_3)
func SPI_IOC_RD_MODE() uintptr {
	return ioctl.IOR(SPI_IOC_MAGIC, 1, 1)
}

// Write of SPI mode (SPI_MODE_0..SPI_MODE_3)
func SPI_IOC_WR_MODE() uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 1, 1)
}

// Read SPI bit justification
func SPI_IOC_RD_LSB_FIRST() uintptr {
	return ioctl.IOR(SPI_IOC_MAGIC, 2, 1)
}

// Write SPI bit justification
func SPI_IOC_WR_LSB_FIRST() uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 2, 1)
}

// Read SPI device word length (1..N)
func SPI_IOC_RD_BITS_PER_WORD() uintptr {
	return ioctl.IOR(SPI_IOC_MAGIC, 3, 1)
}

// Write SPI device word length (1..N)
func SPI_IOC_WR_BITS_PER_WORD() uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 3, 1)
}

// Read SPI device default max speed hz
func SPI_IOC_RD_MAX_SPEED_HZ() uintptr {
	return ioctl.IOR(SPI_IOC_MAGIC, 4, 4)
}

// Write SPI device default max speed hz
func SPI_IOC_WR_MAX_SPEED_HZ() uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 4, 4)
}

// SPI_IOC_MESSAGE: Write custom SPI message
func SPI_IOC_MESSAGE(n uintptr) uintptr {
	return ioctl.IOW(SPI_IOC_MAGIC, 0, uintptr(SPI_MESSAGE_SIZE(n)))
}

func SPI_MESSAGE_SIZE(n uintptr) uintptr {
	if (n * unsafe.Sizeof(SpiIocTransfert{})) < (1 << ioctl.IOC_SIZEBITS) {
		return (n * unsafe.Sizeof(SpiIocTransfert{}))
	}
	return 0
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
	spiDevice := fmt.Sprintf("%s%d.%d", SPIDEV, spi.Bus, spi.Chip)

	logger.NewLogger().WithFields(logrus.Fields{
		"device": spiDevice,
	}).Info("spi.Open")

	var err error
	spi.file, err = os.OpenFile(spiDevice, os.O_RDWR, 0)
	// spi.file, err = os.Create(spiDevice)
	if err != nil {
		return fmt.Errorf("I can't see %s. Have you enabled the SPI module?", spiDevice)
	}

	return nil
}

// Closes SPI device
func (spi *SPIDevice) Close() error {
	err := spi.file.Close()
	if err != nil {
		return fmt.Errorf("Error closing spi", err)
	}
	return nil
}

// Send : Sends bytes over SPI channel and returns []byte response
func (spi *SPIDevice) Send(bytes_to_send []byte) ([]byte, error) {
	wBuffer := bytes_to_send
	rBuffer := [3]byte{}

	// generates message
	transfer := SpiIocTransfert{}
	transfer.txBuf = uint64(uintptr(unsafe.Pointer(&wBuffer[0])))
	transfer.rxBuf = uint64(uintptr(unsafe.Pointer(&rBuffer[0])))
	transfer.length = uint32(len(wBuffer))
	transfer.delayUsecs = SPI_DELAY
	transfer.csChange = 1
	transfer.bitsPerWord = spi.bpw
	transfer.speedHz = spi.speed

	// sends message over SPI
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_MESSAGE(1), uintptr(unsafe.Pointer(&transfer)))
	if err != nil {
		return nil, fmt.Errorf("Error on sending: %s\n", err)
	}

	// generates a valid response
	ret := make([]byte, unsafe.Sizeof(rBuffer))
	for i := range ret {
		ret[i] = rBuffer[i]
	}

	return ret, nil
}

func (spi *SPIDevice) SetMode(mode uint8) error {
	spi.mode = mode
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_WR_MODE(), uintptr(unsafe.Pointer(&mode)))
	if err != nil {
		return fmt.Errorf("Error setting mode: %s\n", err)
	}
	return nil
}

func (spi *SPIDevice) SetBitsPerWord(bpw uint8) error {
	spi.bpw = bpw
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_WR_BITS_PER_WORD(), uintptr(unsafe.Pointer(&bpw)))
	if err != nil {
		return fmt.Errorf("Error setting bits per word: %s\n", err)
	}
	return nil
}

func (spi *SPIDevice) SetSpeed(speed uint32) error {
	spi.speed = speed
	err := ioctl.IOCTL(spi.file.Fd(), SPI_IOC_WR_MAX_SPEED_HZ(), uintptr(unsafe.Pointer(&speed)))
	if err != nil {
		return fmt.Errorf("Error setting speed: %s\n", err)
	}
	return nil
}
