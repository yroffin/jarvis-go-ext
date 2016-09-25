package mfrc522

import (
	"fmt"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/spi"
	"github.com/yroffin/jarvis-go-ext/server/utils/native/wiringpi"
)

/**
 * Cf. https://github.com/pkourany/RC522_RFID/blob/master/RFID.cpp
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
	PICC_SElECTTAG = 0x93
	PICC_AUTHENT1A = 0x60
	PICC_AUTHENT1B = 0x61
	PICC_READ      = 0x30
	PICC_WRITE     = 0xA0
	PICC_DECREMENT = 0xC0
	PICC_INCREMENT = 0xC1
	PICC_RESTORE   = 0xC2
	PICC_TRANSFER  = 0xB0
	PICC_HALT      = 0x50

	MI_OK       = 0
	MI_NOTAGERR = 1
	MI_ERR      = 2
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
	wiringpi.DigitalWrite(22, 0)
	_, err = mfrc522.spiDevice.Send([]byte{(addr << 1) & 0x7E, val})
	wiringpi.DigitalWrite(22, 1)
	if err != nil {
		logger.NewLogger().WithFields(logrus.Fields{
			"addr":  fmt.Sprintf("%02x", addr),
			"Error": err,
		}).Error("Write")
	} else {
		logger.NewLogger().WithFields(logrus.Fields{
			"addr": fmt.Sprintf("%02x", addr),
			"val":  fmt.Sprintf("%02x", val),
		}).Debug("Write")
	}
}

// Read : read on mfrc522 component
func (mfrc522 *Mfrc522) Read(addr byte) (byte, error) {
	var value []byte
	var err error
	wiringpi.DigitalWrite(22, 0)
	value, err = mfrc522.spiDevice.Send([]byte{((addr << 1) & 0x7E) | 0x80, 0})
	wiringpi.DigitalWrite(22, 1)
	if err != nil {
		logger.NewLogger().WithFields(logrus.Fields{
			"addr":  fmt.Sprintf("%02x", addr),
			"Error": err,
		}).Error("Read")
		return 0, err
	}
	logger.NewLogger().WithFields(logrus.Fields{
		"addr": fmt.Sprintf("%02x", addr),
		"val":  fmt.Sprintf("%02x", value[1]),
	}).Debug("Read")
	return value[1], nil
}

// Reset : reset mfrc522 component
func (mfrc522 *Mfrc522) Reset() {
	logger.NewLogger().WithFields(logrus.Fields{}).Info("Reset")
	mfrc522.Write(CommandReg, PCD_RESETPHASE)
}

// SetBitMask : reset bit mask
func (mfrc522 *Mfrc522) SetBitMask(reg byte, mask byte) {
	var value byte
	value, _ = mfrc522.Read(reg)
	mfrc522.Write(reg, value|mask)
	logger.NewLogger().WithFields(logrus.Fields{
		"reg":  fmt.Sprintf("%02x", reg),
		"mask": fmt.Sprintf("%02x", value|mask),
	}).Info("SetBitMask")
}

// ClearBitMask : reset bit mask
func (mfrc522 *Mfrc522) ClearBitMask(reg byte, mask byte) {
	var value byte
	value, _ = mfrc522.Read(reg)
	mfrc522.Write(reg, value&(^mask))
	logger.NewLogger().WithFields(logrus.Fields{
		"reg":           fmt.Sprintf("%02x", reg),
		"value":         fmt.Sprintf("%02x", value),
		"value/bitwise": fmt.Sprintf("%02x", ^mask),
		"mask":          fmt.Sprintf("%02x", mask),
		"mask/write":    fmt.Sprintf("%02x", value&(^mask)),
	}).Info("ClearBitMask")
}

// AntennaOn : set antenna on
func (mfrc522 *Mfrc522) AntennaOn() {
	var value byte
	value, _ = mfrc522.Read(TxControlReg)
	if ^(value & 0x03) > 0 {
		mfrc522.SetBitMask(TxControlReg, 0x03)
	}
	logger.NewLogger().WithFields(logrus.Fields{}).Info("AntennaOn")
}

// AntennaOff : set antenna off
func (mfrc522 *Mfrc522) AntennaOff() {
	mfrc522.ClearBitMask(TxControlReg, 0x03)
	logger.NewLogger().WithFields(logrus.Fields{}).Info("AntennaOff")
}

// StopCrypto1 : set antenna off
func (mfrc522 *Mfrc522) StopCrypto1() {
	mfrc522.ClearBitMask(Status2Reg, 0x08)
	logger.NewLogger().WithFields(logrus.Fields{}).Info("StopCrypto1")
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

// ToCard : write to card
func (mfrc522 *Mfrc522) ToCard(command byte, sendData []byte) (int, []byte, int) {
	var status = MI_ERR
	var commIrqValue byte = 0x00
	var backData = make([]byte, 0)
	var backLen int
	var irqEn byte = 0x00
	var waitIRq byte = 0x00

	var buffer string
	for i := 0; i < len(sendData); i++ {
		if i == 0 {
			buffer = buffer + fmt.Sprintf("0x%02x", sendData[i])
		} else {
			buffer = buffer + fmt.Sprintf(",0x%02x", sendData[i])
		}
	}

	logger.NewLogger().WithFields(logrus.Fields{
		"command":  fmt.Sprintf("%02x", command),
		"sendData": buffer,
	}).Info("ToCard")

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
	var stop = false
	for ; index > 0 && stop == false; index-- {
		// ComIrqReg register bit descriptions
		// CommIrqReg[7..0]
		// Set1 TxIRq RxIRq IdleIRq HiAlerIRq LoAlertIRq ErrIRq TimerIRq
		commIrqValue, _ = mfrc522.Read(ComIrqReg)
		stop = ((commIrqValue & 0x01) == 0x00) && ((commIrqValue & waitIRq) == 0x00)
	}
	logger.NewLogger().WithFields(logrus.Fields{
		"TimerIRq":   (commIrqValue & 0x01) == 0x01,
		"ErrIRq":     (commIrqValue & 0x02) == 0x02,
		"LoAlertIRq": (commIrqValue & 0x04) == 0x04,
		"HiAlertIRq": (commIrqValue & 0x08) == 0x08,
		"IdleIRq":    (commIrqValue & 0x10) == 0x10,
		"RxIRq":      (commIrqValue & 0x20) == 0x20,
		"TxIRq":      (commIrqValue & 0x40) == 0x40,
		"Set1":       (commIrqValue & 0x80) == 0x80,
	}).Info("ToCard::commIrqValue")

	logger.NewLogger().WithFields(logrus.Fields{
		"index": index,
	}).Info("ToCard")

	// StartSend : stop the transmission of data
	mfrc522.ClearBitMask(BitFramingReg, 0x80)

	if index != 0 {
		var errorValue byte
		errorValue, _ = mfrc522.Read(ErrorReg)
		logger.NewLogger().WithFields(logrus.Fields{
			"ProtocolErr": (errorValue & 0x01) == 0x01,
			"ParityErr":   (errorValue & 0x02) == 0x02,
			"CrcErr":      (errorValue & 0x04) == 0x04,
			"ColErr":      (errorValue & 0x08) == 0x08,
			"BufferOvfl":  (errorValue & 0x10) == 0x10,
			"Reserved":    (errorValue & 0x20) == 0x20,
			"TempErr":     (errorValue & 0x40) == 0x40,
			"WrErr":       (errorValue & 0x80) == 0x80,
		}).Info("ToCard::errorValue")
		if (errorValue & 0x1B) == 0x00 {
			status = MI_OK

			if (commIrqValue & irqEn & 0x01) != 0x00 {
				status = MI_NOTAGERR
				logger.NewLogger().WithFields(logrus.Fields{
					"status": status,
				}).Info("ToCard::MI_NOTAGERR")
			}

			if command == PCD_TRANSCEIVE {
				var lastBits byte
				var fifoLevelValue, _ = mfrc522.Read(FIFOLevelReg)
				lastBits, _ = mfrc522.Read(ControlReg)
				logger.NewLogger().WithFields(logrus.Fields{
					"fifoLevelValue": fmt.Sprintf("%08x", fifoLevelValue),
					"lastBits":       fmt.Sprintf("%08x", lastBits),
				}).Info("ToCard::PCD_TRANSCEIVE::check")
				// RxLastBits[2:0]
				// if this value is 000b, the whole byte is valid
				if (lastBits & 0x07) != 0x00 {
					backLen = int((fifoLevelValue-1)*8 + lastBits)
				} else {
					backLen = int(fifoLevelValue * 8)
				}
				if fifoLevelValue == 0 {
					fifoLevelValue = 1
				}
				if fifoLevelValue > MAX_LEN {
					fifoLevelValue = MAX_LEN
				}
				logger.NewLogger().WithFields(logrus.Fields{
					"fifoLevelValue": fmt.Sprintf("%08x", fifoLevelValue),
					"lastBits":       fmt.Sprintf("%08x", lastBits),
					"backLen":        fmt.Sprintf("%08x", backLen),
				}).Info("ToCard::PCD_TRANSCEIVE")
				for i := 0; i < int(fifoLevelValue); i++ {
					var data, _ = mfrc522.Read(FIFODataReg)
					backData = append(backData, data)
				}
			}
		} else {
			status = MI_ERR
		}
	}

	logger.NewLogger().WithFields(logrus.Fields{
		"status":   status,
		"backData": backData,
		"backLen":  backLen,
	}).Info("ToCard::result")

	return status, backData, backLen
}

// ReadCard : read from card
func (mfrc522 *Mfrc522) ReadCard(blockAddr byte) (byte, []byte, error) {
	var recvData []byte = make([]byte, 0)
	recvData = append(recvData, PICC_READ)
	recvData = append(recvData, blockAddr)
	var crc = mfrc522.calulateCRC(recvData)
	recvData = append(recvData, crc[0])
	recvData = append(recvData, crc[1])
	var status, backData, _ = mfrc522.ToCard(PCD_TRANSCEIVE, recvData)
	if status != MI_OK {
		return blockAddr, recvData, fmt.Errorf("Error while reading !")
	}
	if len(backData) == 16 {
		return blockAddr, backData, nil
	}
	return blockAddr, recvData, fmt.Errorf("Internal error")
}

// WriteCard : write from card
func (mfrc522 *Mfrc522) WriteCard() {
}

// Request : request
// Output tagType
// 0x4400 = Mifare_UltraLight
// 0x0400 = Mifare_One(S50)
// 0x0200 = Mifare_One(S70)
// 0x0800 = Mifare_Pro(X)
// 0x4403 = Mifare_DESFire
func (mfrc522 *Mfrc522) Request(reqMode byte) (int, []byte) {
	var tagType []byte = make([]byte, 1)
	tagType[0] = reqMode

	// TxLastBits[2:0]
	// used for transmission of bit oriented frames: defines the number of bits of the last byte that will be transmitted
	mfrc522.Write(BitFramingReg, 0x07)
	var status, backData, backLen = mfrc522.ToCard(PCD_TRANSCEIVE, tagType)

	if (status != MI_OK) || (backLen != 0x10) {
		status, backData, backLen = mfrc522.ToCard(PCD_TRANSCEIVE, tagType)
	}

	var buffer string
	for i := 0; i < len(backData); i++ {
		if i == 0 {
			buffer = buffer + fmt.Sprintf("0x%02x", backData[i])
		} else {
			buffer = buffer + fmt.Sprintf(",0x%02x", backData[i])
		}
	}

	logger.NewLogger().WithFields(logrus.Fields{
		"reqMode":  fmt.Sprintf("%02x", reqMode),
		"status":   fmt.Sprintf("%02x", status),
		"backData": buffer,
		"backLen":  fmt.Sprintf("%02x", backLen),
	}).Info("Request")

	if (status != MI_OK) || (backLen != 0x10) {
		status = MI_ERR
	} else {
		status = MI_OK
	}

	return status, backData
}

// Anticoll : dump anticoll
func (mfrc522 *Mfrc522) Anticoll() (int, []byte) {
	var serNum = make([]byte, 0)

	mfrc522.Write(BitFramingReg, 0x00)

	serNum = append(serNum, PICC_ANTICOLL)
	serNum = append(serNum, 0x20)

	var status, backData, backBits = mfrc522.ToCard(PCD_TRANSCEIVE, serNum)

	var buffer string
	for i := 0; i < len(backData); i++ {
		if i == 0 {
			buffer = buffer + fmt.Sprintf("0x%02x", backData[i])
		} else {
			buffer = buffer + fmt.Sprintf(",0x%02x", backData[i])
		}
	}

	if status == MI_OK {
		if len(backData) == 5 {
			var index = 0
			var serNumCheck byte = 0
			for i := 0; i < 4; i++ {
				serNumCheck = serNumCheck ^ backData[i]
				index++
			}
			if serNumCheck != backData[index] {
				status = MI_ERR
			}
		} else {
			status = MI_ERR
		}
	}

	logger.NewLogger().WithFields(logrus.Fields{
		"status":   fmt.Sprintf("%02x", status),
		"backData": buffer,
		"backBits": fmt.Sprintf("%02x", backBits),
	}).Info("Anticoll")

	return status, backData
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

	logger.NewLogger().WithFields(logrus.Fields{
		"pIndata":  pIndata,
		"pOutData": pOutData,
	}).Info("calulateCRC")

	return pOutData
}

// SelectTag : Tag
func (mfrc522 *Mfrc522) SelectTag(serNum [5]byte) int {
	var buf = make([]byte, 0)

	buf = append(buf, PICC_SElECTTAG)
	buf = append(buf, 0x70)
	for i := 0; i < len(serNum); i++ {
		buf = append(buf, serNum[i])
	}
	var crc = mfrc522.calulateCRC(buf)
	buf = append(buf, crc[0])
	buf = append(buf, crc[1])

	var status, backData, backLen = mfrc522.ToCard(PCD_TRANSCEIVE, buf)
	if (status == MI_OK) && (backLen == 0x18) {
		logger.NewLogger().WithFields(logrus.Fields{
			"serNum":   serNum,
			"backData": backData,
		}).Info("SelectTag")
		return int(backData[0])
	}
	return 0
}

// Auth : Auth
func (mfrc522 *Mfrc522) Auth(authMode byte, blockAddr byte, sectorkey [6]byte, serNum [5]byte) int {
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
	var status, _, _ = mfrc522.ToCard(PCD_AUTHENT, buff)
	// Check result
	if status != MI_OK {
		// Error
		logger.NewLogger().WithFields(logrus.Fields{
			"status": status,
		}).Error("Auth::ToCard")
	}
	var statusValue, _ = mfrc522.Read(Status2Reg)
	if (statusValue & 0x08) != 0 {
		// Error
		logger.NewLogger().WithFields(logrus.Fields{
			"statusValue": statusValue,
		}).Error("Auth::Read")
	}
	return status
}

// DumpClassic1K : DumpClassic1K
func (mfrc522 *Mfrc522) DumpClassic1K(key [6]byte, uid [5]byte) {
	logger.NewLogger().WithFields(logrus.Fields{
		"key": key,
		"uid": uid,
	}).Info("DumpClassic1K")

	for i := 0; i < 64; i++ {
		var status int = mfrc522.Auth(PICC_AUTHENT1A, byte(i), key, uid)
		// Check if authenticated
		if status == MI_OK {
			var blockAdr, dumpArray, err = mfrc522.ReadCard(byte(i))
			if err != nil {
				logger.NewLogger().WithFields(logrus.Fields{
					"blockAdr":  blockAdr,
					"dumpArray": dumpArray,
				}).Info("DumpClassic1K")
			}
		} else {
			// Error
			logger.NewLogger().WithFields(logrus.Fields{
				"status": status,
			}).Error("DumpClassic1K")
		}
	}
}

var instance *Mfrc522
var once sync.Once

// GetInstance : singleton instance
func GetInstance() *Mfrc522 {
	once.Do(func() {
		instance = new(Mfrc522)
		instance.Init()
	})
	return instance
}

// Init : Init
func (mfrc522 *Mfrc522) Init() {
	mfrc522.spiDevice = spi.NewSPIDevice(0, 0)
	mfrc522.spiDevice.Open()
	mfrc522.spiDevice.SetSpeed(1000000)
	mfrc522.spiDevice.SetBitsPerWord(8)
	mfrc522.spiDevice.SetMode(0)

	logger.NewLogger().WithFields(logrus.Fields{}).Info("SetPinMode")
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

	logger.NewLogger().WithFields(logrus.Fields{}).Info("Init")
}
