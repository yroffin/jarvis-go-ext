package wiringpi

// #cgo arm CFLAGS: -marm
// #cgo arm LDFLAGS: -lwiringPi
// extern int wiringPiSetupInit();
// extern void  delay             	(unsigned int howLong);
// extern void  delayMicroseconds 	(unsigned int howLong);
// extern void digitalWrite        (int pin, int value);
// extern int  wiringPiSetup       (void);
// extern void pinMode             (int pin, int mode);
// extern int  setuid      		(int uid);
import "C"

// WiringPiInit : initialize the library
func Init() int {
	return int(C.wiringPiSetupInit())
}

// PinMode : call wiringpi pinMode
func PinMode(pin int, value int) {
	C.pinMode(C.int(pin), C.int(value))
}

// DigitalWrite : call wiringpi digitalWrite
func DigitalWrite(pin int, value int) {
	C.digitalWrite(C.int(pin), C.int(value))
}

// DelayMicroseconds : call wiringpi delayMicroseconds
func DelayMicroseconds(delay uint) {
	C.delayMicroseconds(C.uint(delay))
}

// Delay : call wiringpi delay
func Delay(delay uint) {
	C.delay(C.uint(delay))
}
