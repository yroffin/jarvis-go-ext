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
