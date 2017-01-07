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

package mfrc522

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/yroffin/jarvis-go-ext/server/types"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/spi"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/wiringpi"
)

/**
 * Cf. https://github.com/pkourany/RC522_RFID/blob/master/RFID.cpp
 * Cf. http://www.nxp.com/documents/data_sheet/MFRC522.pdf
 */

const (
	MAX_LEN = 16

	PCD_IDLE       = 0x00
	PCD_AUTHENT    = 0x0E
	PCD_RECEIVE    = 0x08
	PCD_TRANSMIT   = 0x04
	PCD_TRANSCEIVE = 0x0C
	PCD_RESETPHASE = 0x0F
	PCD_CALCCRC    = 0x03

	PICC_REQIDL    = 0x26
	PICC_REQALL    = 0x52
	PICC_ANTICOLL  = 0x93
	PICC_SELECTTAG = 0x93
	PICC_AUTHENT1A = 0x60
	PICC_AUTHENT1B = 0x61
	PICC_READ      = 0x30
	PICC_WRITE     = 0xA0
	PICC_DECREMENT = 0xC0
	PICC_INCREMENT = 0xC1
	PICC_RESTORE   = 0xC2
	PICC_TRANSFER  = 0xB0
	PICC_HALT      = 0x50
)

const (
	MI_OK          = 0
	MI_NOTAGERR    = 1
	MI_ERR         = 2
	MI_ERR_CRC     = 3
	MI_ERR_CRC_LEN = 4
	MI_ERR_SEND    = 5
	MI_ERR_REQUEST = 6
	MI_ERR_TIMEOUT = 7
	// Status2Reg register MFCrypto1On bit not set
	// indicates that the MIFARE Crypto1 unit is switched on and
	// therefore all data communication with the card is encrypted
	MI_ERR_CRYPTO = 8
)

const (
	Reserved00    = 0x00
	CommandReg    = 0x01
	ComIEnReg     = 0x02
	DivlEnReg     = 0x03
	ComIrqReg     = 0x04
	DivIrqReg     = 0x05
	ErrorReg      = 0x06
	Status1Reg    = 0x07
	Status2Reg    = 0x08
	FIFODataReg   = 0x09
	FIFOLevelReg  = 0x0A
	WaterLevelReg = 0x0B
	ControlReg    = 0x0C
	BitFramingReg = 0x0D
	CollReg       = 0x0E
	Reserved01    = 0x0F

	Reserved10     = 0x10
	ModeReg        = 0x11
	TxModeReg      = 0x12
	RxModeReg      = 0x13
	TxControlReg   = 0x14
	TxAutoReg      = 0x15
	TxSelReg       = 0x16
	RxSelReg       = 0x17
	RxThresholdReg = 0x18
	DemodReg       = 0x19
	Reserved11     = 0x1A
	Reserved12     = 0x1B
	MifareReg      = 0x1C
	Reserved13     = 0x1D
	Reserved14     = 0x1E
	SerialSpeedReg = 0x1F

	Reserved20        = 0x20
	CRCResultRegM     = 0x21
	CRCResultRegL     = 0x22
	Reserved21        = 0x23
	ModWidthReg       = 0x24
	Reserved22        = 0x25
	RFCfgReg          = 0x26
	GsNReg            = 0x27
	CWGsPReg          = 0x28
	ModGsPReg         = 0x29
	TModeReg          = 0x2A
	TPrescalerReg     = 0x2B
	TReloadRegH       = 0x2C
	TReloadRegL       = 0x2D
	TCounterValueRegH = 0x2E
	TCounterValueRegL = 0x2F

	Reserved30      = 0x30
	TestSel1Reg     = 0x31
	TestSel2Reg     = 0x32
	TestPinEnReg    = 0x33
	TestPinValueReg = 0x34
	TestBusReg      = 0x35
	AutoTestReg     = 0x36
	VersionReg      = 0x37
	AnalogTestReg   = 0x38
	TestDAC1Reg     = 0x39
	TestDAC2Reg     = 0x3A
	TestADCReg      = 0x3B
	Reserved31      = 0x3C
	Reserved32      = 0x3D
	Reserved33      = 0x3E
	Reserved34      = 0x3F
)

// Mfrc522 : instance Mfrc522 device struct
type Mfrc522 struct {
	spiDevice (*spi.SPIDevice)
}

// Write : write on mfrc522 component
func (mfrc522 *Mfrc522) Write(addr byte, val byte) {
	var err error
	_, err = mfrc522.spiDevice.Send([]byte{(addr << 1) & 0x7E, val})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"addr":  fmt.Sprintf("%02x", addr),
			"Error": err,
		}).Error("Write")
	} else {
		logrus.WithFields(logrus.Fields{
			"addr": fmt.Sprintf("%02x", addr),
			"val":  fmt.Sprintf("%02x", val),
		}).Debug("Write")
	}
}

// Read : read on mfrc522 component
func (mfrc522 *Mfrc522) Read(addr byte) (byte, error) {
	var value []byte
	var err error
	value, err = mfrc522.spiDevice.Send([]byte{((addr << 1) & 0x7E) | 0x80, 0})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"addr":  fmt.Sprintf("%02x", addr),
			"Error": err,
		}).Error("Read")
		return 0, err
	}
	logrus.WithFields(logrus.Fields{
		"addr": fmt.Sprintf("%02x", addr),
		"val":  fmt.Sprintf("%02x", value[1]),
	}).Debug("Read")
	return value[1], nil
}

// Reset : reset mfrc522 component
func (mfrc522 *Mfrc522) Reset() {
	logrus.WithFields(logrus.Fields{}).Info("Reset")
	mfrc522.Write(CommandReg, PCD_RESETPHASE)
}

// SetBitMask : reset bit mask
func (mfrc522 *Mfrc522) SetBitMask(reg byte, mask byte) {
	var value byte
	value, _ = mfrc522.Read(reg)
	mfrc522.Write(reg, value|mask)
	logrus.WithFields(logrus.Fields{
		"reg":  fmt.Sprintf("%02x", reg),
		"mask": fmt.Sprintf("%02x", value|mask),
	}).Debug("SetBitMask")
}

// ClearBitMask : reset bit mask
func (mfrc522 *Mfrc522) ClearBitMask(reg byte, mask byte) {
	var value byte
	value, _ = mfrc522.Read(reg)
	mfrc522.Write(reg, value&(^mask))
	logrus.WithFields(logrus.Fields{
		"reg":           fmt.Sprintf("%02x", reg),
		"value":         fmt.Sprintf("%02x", value),
		"value/bitwise": fmt.Sprintf("%02x", ^mask),
		"mask":          fmt.Sprintf("%02x", mask),
		"mask/write":    fmt.Sprintf("%02x", value&(^mask)),
	}).Debug("ClearBitMask")
}

// AntennaOn : set antenna on
func (mfrc522 *Mfrc522) AntennaOn() {
	var value byte
	value, _ = mfrc522.Read(TxControlReg)
	if ^(value & 0x03) > 0 {
		mfrc522.SetBitMask(TxControlReg, 0x03)
	}
	logrus.WithFields(logrus.Fields{}).Info("AntennaOn")
}

// AntennaOff : set antenna off
func (mfrc522 *Mfrc522) AntennaOff() {
	mfrc522.ClearBitMask(TxControlReg, 0x03)
	logrus.WithFields(logrus.Fields{}).Info("AntennaOff")
}

// StopCrypto1 : set antenna off
func (mfrc522 *Mfrc522) StopCrypto1() {
	mfrc522.ClearBitMask(Status2Reg, 0x08)
	logrus.WithFields(logrus.Fields{}).Debug("StopCrypto1")
}

// appendByte : internal append function to array
func appendByte(slice []byte, data ...byte) []byte {
	m := len(slice)
	n := m + len(data)
	if n > cap(slice) { // if necessary, reallocate
		// allocate double what's needed, for future growth.
		newSlice := make([]byte, (n+1)*2)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0:n]
	copy(slice[m:n], data)
	return slice
}

// dataDump, convert to string
func dataDump(sendData []byte) string {
	var hexadump, asciidump string
	for i := 0; i < len(sendData); i++ {
		if i == 0 {
			hexadump = hexadump + fmt.Sprintf("0x%02x", sendData[i])
			asciidump = asciidump + fmt.Sprintf("%d", sendData[i])
		} else {
			hexadump = hexadump + fmt.Sprintf(",0x%02x", sendData[i])
			asciidump = asciidump + fmt.Sprintf(",%d", sendData[i])
		}
	}
	return hexadump + "::" + asciidump
}

// ToCard : write to card
func (mfrc522 *Mfrc522) ToCard(command byte, sendData []byte, expected int) (int, []byte) {
	var status = MI_ERR
	var commIrqValue byte = 0x00
	var backData = make([]byte, 0)
	var irqEn byte = 0x00
	var waitIRq byte = 0x00

	logrus.WithFields(logrus.Fields{
		"command":  command,
		"sendData": dataDump(sendData),
		"expected": expected,
	}).Debug("ToCard::input")

	if command == PCD_AUTHENT {
		// Bit 4 IdleIEn : allows the idle interrupt request (IdleIRq bit) to be propagated to pin IRQ
		// Bit 1 ErrIEn : allows the error interrupt request (ErrIRq bit) to be propagated to pin IRQ
		irqEn = 0x12
		waitIRq = 0x10
	}

	if command == PCD_TRANSCEIVE {
		// Bit 6 TxIEn : allows the transmitter interrupt request (TxIRq bit) to be propagated to pin IRQ
		// Bit 5 RxIEn : allows the receiver interrupt request (RxIRq bit) to be propagated to pin IRQ
		// Bit 4 IdleIEn : allows the idle interrupt request (IdleIRq bit) to be propagated to pin IRQ
		// Bit 2 LoAlertIEn : allows the low alert interrupt request (LoAlertIRq bit) to be propagated to pin IRQ
		// Bit 1 ErrIEn : allows the error interrupt request (ErrIRq bit) to be propagated to pin IRQ
		// Bit 0 TimerIEn : allows the timer interrupt request (TimerIRq bit) to be propagated to pin IRQ
		irqEn = 0x77
		waitIRq = 0x30
	}

	// 0x80 : signal on pin IRQ is inverted with respect to the Status1Reg register’s IRq bit
	mfrc522.Write(ComIEnReg, irqEn|0x80)
	// 0x80 : indicates that the marked bits in the ComIrqReg register are set
	mfrc522.ClearBitMask(ComIrqReg, 0x80)
	// 0x80 : immediately clears the internal FIFO buffer’s read and write pointer and ErrorReg register’s BufferOvfl bit
	mfrc522.SetBitMask(FIFOLevelReg, 0x80)
	// Clear any command
	mfrc522.Write(CommandReg, PCD_IDLE)
	// Write data in fifo tunnel
	for i := 0; i < len(sendData); i++ {
		mfrc522.Write(FIFODataReg, sendData[i])
	}

	mfrc522.Write(CommandReg, command)

	if command == PCD_TRANSCEIVE {
		// StartSend : starts the transmission of data, only valid in combination with the Transceive command
		mfrc522.SetBitMask(BitFramingReg, 0x80)
	}

	var index = 2000
	for {
		// ComIrqReg register bit descriptions
		// CommIrqReg[7..0]
		// Set1 TxIRq RxIRq IdleIRq HiAlerIRq LoAlertIRq ErrIRq TimerIRq
		commIrqValue, _ = mfrc522.Read(ComIrqReg)
		index--
		if !((index != 0) && !((commIrqValue & 0x01) != 0x00) && !((commIrqValue & waitIRq) != 0x00)) {
			break
		}
	}

	logrus.WithFields(logrus.Fields{
		"TimerIRq":   (commIrqValue & 0x01) == 0x01,
		"ErrIRq":     (commIrqValue & 0x02) == 0x02,
		"LoAlertIRq": (commIrqValue & 0x04) == 0x04,
		"HiAlertIRq": (commIrqValue & 0x08) == 0x08,
		"IdleIRq":    (commIrqValue & 0x10) == 0x10,
		"RxIRq":      (commIrqValue & 0x20) == 0x20,
		"TxIRq":      (commIrqValue & 0x40) == 0x40,
		"Set1":       (commIrqValue & 0x80) == 0x80,
	}).Debug("ToCard::commIrqValue")

	logrus.WithFields(logrus.Fields{
		"index": index,
	}).Debug("ToCard")

	// StartSend : stop the transmission of data
	mfrc522.ClearBitMask(BitFramingReg, 0x80)

	if index != 0 {
		var errorValue byte
		errorValue, _ = mfrc522.Read(ErrorReg)

		if (errorValue & 0x1B) == 0x00 {
			status = MI_OK

			if (commIrqValue & irqEn & 0x01) != 0x00 {
				// The ComIrqReg register’s ErrIRq bit indicates an error detected by the contactless UART
				// during send or receive. This is indicated when any bit is set to logic 1 in register ErrorReg
				status = MI_NOTAGERR
				logrus.WithFields(logrus.Fields{
					"status": status,
				}).Debug("ToCard::MI_NOTAGERR")
				return MI_NOTAGERR, nil
			} else {
				logrus.WithFields(logrus.Fields{
					"status": status,
				}).Debug("ToCard::MI_OK")
			}

			if command == PCD_TRANSCEIVE {
				var overflow int
				for overflow = 0; len(backData) != expected && overflow <= 16; overflow++ {
					var fifoLevelValue, _ = mfrc522.Read(FIFOLevelReg)
					// RxLastBits[2:0]
					// if this value is 000b, the whole byte is valid
					var lastBits, _ = mfrc522.Read(ControlReg)
					for (lastBits & 0x07) != 0x00 {
						lastBits, _ = mfrc522.Read(ControlReg)
					}
					// Append to current buffer
					for i := 0; i < int(fifoLevelValue); i++ {
						var data, _ = mfrc522.Read(FIFODataReg)
						backData = append(backData, data)
					}
					// Compute size in bits
					logrus.WithFields(logrus.Fields{
						"overflow":       overflow,
						"fifoLevelValue": fmt.Sprintf("%08x", fifoLevelValue),
						"lastBits":       fmt.Sprintf("%08x", lastBits),
						"len(backData)":  fmt.Sprintf("%08x", len(backData)),
					}).Debug("ToCard::PCD_TRANSCEIVE")
				}
				// Verify timeout
				if overflow > 16 {
					logrus.WithFields(logrus.Fields{
						"overflow": overflow,
					}).Error("ToCard::timeout")
					status = MI_ERR_TIMEOUT
				}
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"ProtocolErr": (errorValue & 0x01) == 0x01,
				"ParityErr":   (errorValue & 0x02) == 0x02,
				"CrcErr":      (errorValue & 0x04) == 0x04,
				"ColErr":      (errorValue & 0x08) == 0x08,
				"BufferOvfl":  (errorValue & 0x10) == 0x10,
				"Reserved":    (errorValue & 0x20) == 0x20,
				"TempErr":     (errorValue & 0x40) == 0x40,
				"WrErr":       (errorValue & 0x80) == 0x80,
			}).Error("ToCard::errorValue")
			status = MI_ERR_SEND
		}
	}

	logrus.WithFields(logrus.Fields{
		"status":   status,
		"backData": dataDump(backData),
	}).Debug("ToCard::result")

	return status, backData
}

// ReadCard : read from card
func (mfrc522 *Mfrc522) ReadCard(blockAddr byte) (byte, []byte, error) {
	var transceive []byte = make([]byte, 0)
	transceive = append(transceive, PICC_READ)
	transceive = append(transceive, blockAddr)
	var crc = mfrc522.calulateCRC(transceive)
	transceive = append(transceive, crc[0])
	transceive = append(transceive, crc[1])
	var status, backData = mfrc522.ToCard(PCD_TRANSCEIVE, transceive, 18)
	if status != MI_OK {
		return blockAddr, backData, fmt.Errorf("Error while reading !")
	} else {
		logrus.WithFields(logrus.Fields{
			"blockAddr": blockAddr,
			"backData":  dataDump(backData),
		}).Debug("ReadCard")
		return blockAddr, backData, nil
	}
}

// WriteCard : write from card
func (mfrc522 *Mfrc522) WriteCard(blockAddr byte, writeData []byte) (byte, []byte, error) {
	var transceive []byte = make([]byte, 0)
	transceive = append(transceive, PICC_WRITE)
	transceive = append(transceive, blockAddr)
	var crc = mfrc522.calulateCRC(transceive)
	transceive = append(transceive, crc[0])
	transceive = append(transceive, crc[1])
	var status, backData = mfrc522.ToCard(PCD_TRANSCEIVE, transceive, 4)
	if status != MI_OK {
		return blockAddr, transceive, fmt.Errorf("Error while writing !")
	} else {
		logrus.WithFields(logrus.Fields{
			"blockAddr": blockAddr,
			"backData":  dataDump(backData),
		}).Debug("WriteCard")
		return blockAddr, backData, nil
	}

	// transceive ok
	for index := 0; index < len(writeData); index++ {
		transceive = append(transceive, writeData[index])
	}
	crc = mfrc522.calulateCRC(transceive)
	transceive = append(transceive, crc[0])
	transceive = append(transceive, crc[1])
	status, backData = mfrc522.ToCard(PCD_TRANSCEIVE, transceive, 4)
	if status != MI_OK {
		return blockAddr, transceive, fmt.Errorf("Error while writing !")
	} else {
		logrus.WithFields(logrus.Fields{
			"blockAddr": blockAddr,
			"backData":  dataDump(backData),
		}).Debug("WriteCard data written")
		return blockAddr, backData, nil
	}
}

// Request : request
func (mfrc522 *Mfrc522) Request(reqMode byte) (int, []byte) {
	var tagType []byte = make([]byte, 1)
	tagType[0] = reqMode

	logrus.WithFields(logrus.Fields{
		"reqMode": fmt.Sprintf("%02x", reqMode),
	}).Debug("Request")

	// TxLastBits[2:0]
	// used for transmission of bit oriented frames: defines the number of bits of the last byte that will be transmitted
	mfrc522.Write(BitFramingReg, 0x07)
	var status, backData = mfrc522.ToCard(PCD_TRANSCEIVE, tagType, 2)

	if status != MI_OK {
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Debug("Request")

		return status, nil
	}

	logrus.WithFields(logrus.Fields{
		"status":   status,
		"backData": dataDump(backData),
	}).Debug("Request")

	return MI_OK, backData
}

// Anticoll : dump anticoll
func (mfrc522 *Mfrc522) Anticoll() ([]byte, error) {
	var serNum = make([]byte, 2)

	logrus.WithFields(logrus.Fields{}).Debug("Anticoll")

	mfrc522.Write(BitFramingReg, 0x00)

	serNum[0] = PICC_ANTICOLL
	serNum[1] = 0x20

	var status, backData = mfrc522.ToCard(PCD_TRANSCEIVE, serNum, 5)

	if status == MI_OK {
		if len(backData) == 5 {
			var index = 0
			var serNumCheck byte = 0
			for i := 0; i < 4; i++ {
				serNumCheck = serNumCheck ^ backData[i]
				index++
			}
			if serNumCheck != backData[index] {
				logrus.WithFields(logrus.Fields{
					"error": "CRC Error",
				}).Error("Anticoll")
				return nil, fmt.Errorf("CRC Error")
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"error":  "Invalid length",
				"length": len(backData),
			}).Error("Anticoll")
			return nil, fmt.Errorf("Invalid length %d", len(backData))
		}
	}

	logrus.WithFields(logrus.Fields{
		"backData": dataDump(backData),
	}).Debug("Anticoll")

	return backData, nil
}

// CalulateCRC : compute CRC
func (mfrc522 *Mfrc522) calulateCRC(pIndata []byte) []byte {
	mfrc522.ClearBitMask(DivIrqReg, 0x04)  // CRCIrq = 0
	mfrc522.SetBitMask(FIFOLevelReg, 0x80) // FIFO

	// write data
	for i := 0; i < len(pIndata); i++ {
		mfrc522.Write(FIFODataReg, pIndata[i])
	}
	mfrc522.Write(CommandReg, PCD_CALCCRC)

	// Iterate on CRC compute
	var n, _ = mfrc522.Read(DivIrqReg)
	for i := 255; (i > 0) && ((int(n) & 0x04) == 0x00); i-- {
		n, _ = mfrc522.Read(DivIrqReg)
	}

	// Store CRC result
	var pOutData = make([]byte, 2)
	pOutData[0], _ = mfrc522.Read(CRCResultRegL)
	pOutData[1], _ = mfrc522.Read(CRCResultRegM)

	logrus.WithFields(logrus.Fields{
		"pIndata":  pIndata,
		"pOutData": pOutData,
	}).Debug("calulateCRC")

	return pOutData
}

// SelectTag : Tag
func (mfrc522 *Mfrc522) SelectTag(serNum []byte) (int, int) {
	var buf = make([]byte, 0)

	buf = append(buf, PICC_SELECTTAG)
	buf = append(buf, 0x70)
	for i := 0; i < len(serNum); i++ {
		buf = append(buf, serNum[i])
	}
	var crc = mfrc522.calulateCRC(buf)
	buf = append(buf, crc[0])
	buf = append(buf, crc[1])

	logrus.WithFields(logrus.Fields{
		"serNum": serNum,
		"buf":    dataDump(buf),
	}).Debug("SelectTag")

	var status, backData = mfrc522.ToCard(PCD_TRANSCEIVE, buf, 3)
	if status == MI_OK {
		logrus.WithFields(logrus.Fields{
			"backData": dataDump(backData),
		}).Debug("SelectTag")
		return MI_OK, int(backData[0])
	}
	logrus.WithFields(logrus.Fields{
		"backData": dataDump(backData),
	}).Error("SelectTag")
	return MI_ERR, 0
}

// Auth : Auth
// This command manages MIFARE authentication to enable a secure communication to
// any MIFARE Mini, MIFARE 1K and MIFARE 4K card. The following data is written to the
// FIFO buffer before the command can be activated
// • Authentication command code (60h, 61h)
// • Block address
// • Sector key byte 0
// • Sector key byte 1
// • Sector key byte 2
// • Sector key byte 3
// • Sector key byte 4
// • Sector key byte 5
// • Card serial number byte 0
// • Card serial number byte 1
// • Card serial number byte 2
// • Card serial number byte 3
func (mfrc522 *Mfrc522) Auth(authMode byte, blockAddr byte, sectorkey []byte, serNum []byte) int {
	logrus.WithFields(logrus.Fields{
		"authMode":  authMode,
		"blockAddr": blockAddr,
		"sectorkey": sectorkey,
		"serNum":    serNum,
	}).Debug("Auth")

	var buff = make([]byte, 0)
	// First byte should be the authMode (A or B)
	buff = append(buff, authMode)
	// Second byte is the trailerBlock (usually 7)
	buff = append(buff, blockAddr)
	// Now we need to append the authKey which usually is 6 bytes of 0xFF
	for i := 0; i < len(sectorkey); i++ {
		buff = append(buff, sectorkey[i])
	}
	// Next we append the first 4 bytes of the UID
	for i := 0; i < 4; i++ {
		buff = append(buff, serNum[i])
	}

	// Now we start the authentication itself
	var status, _ = mfrc522.ToCard(PCD_AUTHENT, buff, 0)

	// Check result
	if status != MI_OK {
		// Error
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Error("Auth an error occured")
		return MI_ERR
	}
	// Check Status2Reg
	var statusValue, _ = mfrc522.Read(Status2Reg)
	if (statusValue & 0x08) != 0x08 {
		// Error
		logrus.WithFields(logrus.Fields{
			"statusValue": statusValue,
		}).Error("Auth MFCrypto1On not set")
		return MI_ERR_CRYPTO
	}

	logrus.WithFields(logrus.Fields{
		"status": MI_OK,
	}).Debug("Auth successful")

	return MI_OK
}

// DumpClassic1K : DumpClassic1K
func (mfrc522 *Mfrc522) dumpClassic1K(key []byte, uid []byte) ([]types.Mfrc522Sector16, error) {
	var Sectors []types.Mfrc522Sector16 = make([]types.Mfrc522Sector16, 64)

	logrus.WithFields(logrus.Fields{
		"key": key,
		"uid": uid,
	}).Debug("DumpClassic1K")

	for i := 0; i < 64; i++ {
		var status int = mfrc522.Auth(PICC_AUTHENT1A, byte(i), key, uid)
		// Check if authenticated
		if status == MI_OK {
			var _, dumpArray, err = mfrc522.ReadCard(byte(i))
			if err == nil {
				for element := 0; element < 16; element++ {
					Sectors[i].Values[element] = dumpArray[element]
				}
			} else {
				// Error
				logrus.WithFields(logrus.Fields{
					"status": status,
				}).Error("DumpClassic1K")
				return Sectors, fmt.Errorf("Error while reading sector %d", i)
			}
		} else {
			// Error
			logrus.WithFields(logrus.Fields{
				"status": status,
			}).Error("DumpClassic1K")
			return Sectors, fmt.Errorf("Error while reading authentification")
		}
	}

	return Sectors, nil
}

// WriteClassic1K : WriteClassic1K
func (mfrc522 *Mfrc522) WriteClassic1K(key []byte, uid []byte, sector byte, data []byte) error {
	logrus.WithFields(logrus.Fields{
		"key": key,
		"uid": uid,
	}).Debug("WriteClassic1K")

	var status int = mfrc522.Auth(PICC_AUTHENT1A, sector, key, uid)
	// Check if authenticated
	if status == MI_OK {
		var _, _, err = mfrc522.WriteCard(sector, data)
		if err != nil {
			// Error
			logrus.WithFields(logrus.Fields{
				"status": status,
			}).Error("WriteClassic1K")
			return fmt.Errorf("Error while writing")
		}
	} else {
		// Error
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Error("WriteClassic1K")
		return fmt.Errorf("Error while reading authentification")
	}

	return nil
}

// RequestIdle : RequestIdle
// Output tagType
// 0x4400 = Mifare_UltraLight
// 0x0400 = Mifare_One(S50)
// 0x0200 = Mifare_One(S70)
// 0x0800 = Mifare_Pro(X)
// 0x4403 = Mifare_DESFire
func (mfrc522 *Mfrc522) RequestIdle() (string, error) {
	// request for status
	for index := 0; index < 2; index++ {
		var status, backData = mfrc522.Request(PICC_REQIDL)
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Debug("RequestIdle::result")
		if status == MI_OK {
			if backData[0] == 0x04 && backData[1] == 0x00 {
				return "Mifare_One(S50)", nil
			}
			return "Other", nil
		}
		if status == MI_NOTAGERR {
			continue
		}
		return "", fmt.Errorf("Error while requesting for no idle with error %d", status)
	}
	return "", fmt.Errorf("Error while requesting for no idle")
}

var instance *Mfrc522
var once sync.Once

// GetInstance : singleton instance
func GetInstance() *Mfrc522 {
	once.Do(func() {
		instance = new(Mfrc522)
		instance.init()
	})
	return instance
}

// HandleRequestIdle : verify card/nfc tag presence
func (mfrc522 *Mfrc522) handleRequestIdle() (error, string) {
	// request for status
	var tagType, status = mfrc522.RequestIdle()
	if status != nil {
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Debug("Unable to detect tag")
		return fmt.Errorf("Unable to detect tag"), ""
	}

	return nil, tagType
}

// HandleAnticoll : read uuid
func (mfrc522 *Mfrc522) handleAnticoll() (error, []byte) {
	// request for uuid
	var status error
	var data []byte
	data, status = mfrc522.Anticoll()
	if status != nil {
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Error("Unable to detect uuid")
		return fmt.Errorf("Unable to detect uuid"), nil
	}

	return nil, data
}

// handleSelectTag : read uuid
func (mfrc522 *Mfrc522) handleSelectTag(Uid []byte) error {
	var statusSelectTag, _ = mfrc522.SelectTag(Uid)
	if statusSelectTag != 0 {
		logrus.WithFields(logrus.Fields{}).Error("Unable to select tag")
		return fmt.Errorf("Unable to select tag")
	}

	return nil
}

// handleDumpClassic1K : rdump nfc tag
func (mfrc522 *Mfrc522) handleDumpClassic1K(Key []byte, Uid []byte) (error, []types.Mfrc522Sector16) {
	var status error
	var Sectors []types.Mfrc522Sector16
	Sectors, _ = mfrc522.dumpClassic1K(Key, Uid)
	if status != nil {
		// stop any recent auth
		mfrc522.StopCrypto1()
		logrus.WithFields(logrus.Fields{
			"error": status,
		}).Error("Unable to read tag")
		return fmt.Errorf("Unable to read tag"), nil
	}

	return nil, Sectors
}

// WaitForTag : handler for WaitForTag
func (mfrc522 *Mfrc522) WaitForTag() (error, string, []byte) {
	var result error

	logrus.WithFields(logrus.Fields{}).Debug("WaitForTag")

	// request for status
	var tagType string
	result, tagType = mfrc522.handleRequestIdle()
	if result != nil {
		return result, "", nil
	}

	// request for uuid
	var Uuid []byte
	result, Uuid = mfrc522.handleAnticoll()
	if result != nil {
		return result, "", nil
	}

	return nil, tagType, Uuid
}

// DumpClassic1K : handler for DumpClassic1K
func (mfrc522 *Mfrc522) DumpClassic1K(Key []byte) (error, string, []byte, []types.Mfrc522Sector16) {
	var result error

	logrus.WithFields(logrus.Fields{
		"key": Key,
	}).Info("DumpClassic1K")

	// request for status
	var tagType string
	result, tagType = mfrc522.handleRequestIdle()
	if result != nil {
		return result, "", nil, nil
	}

	// request for uuid
	var Uuid []byte
	result, Uuid = mfrc522.handleAnticoll()
	if result != nil {
		return result, "", nil, nil
	}

	// select tag
	result = mfrc522.handleSelectTag(Uuid)
	if result != nil {
		// stop any recent auth
		mfrc522.StopCrypto1()
		return result, "", nil, nil
	}

	// DumpClassic1K
	var sectors []types.Mfrc522Sector16
	result, sectors = mfrc522.handleDumpClassic1K(Key, Uuid)
	if result != nil {
		// stop any recent auth
		mfrc522.StopCrypto1()
		return result, "", nil, nil
	}

	// stop any recent auth
	mfrc522.StopCrypto1()

	return nil, tagType, Uuid, sectors
}

// Init : Init
func (mfrc522 *Mfrc522) init() {
	mfrc522.spiDevice = spi.NewSPIDevice(0, 0)
	mfrc522.spiDevice.Open()
	mfrc522.spiDevice.SetSpeed(1000000)
	mfrc522.spiDevice.SetBitsPerWord(8)
	mfrc522.spiDevice.SetMode(0)

	wiringpi.PinMode(22, 1)
	wiringpi.DigitalWrite(22, 1)

	mfrc522.Reset()
	mfrc522.Write(TModeReg, 0x8D)
	mfrc522.Write(TPrescalerReg, 0x3E)
	mfrc522.Write(TReloadRegL, 0x1E)
	mfrc522.Write(TReloadRegH, 0x00)
	mfrc522.Write(TxAutoReg, 0x40)
	mfrc522.Write(ModeReg, 0x3D)
	mfrc522.AntennaOn()

	logrus.WithFields(logrus.Fields{
		"Config": "on",
	}).Info("Mfrc522")
}
