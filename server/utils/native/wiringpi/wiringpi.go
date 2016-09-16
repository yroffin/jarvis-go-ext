package wiringpi

// #cgo arm CFLAGS: -marm
// #cgo arm LDFLAGS: -lwiringPi
// extern int wiringPiSetupInit();
// extern int pinMode(int pin, int value);
import "C"

// WiringPiInit : initialize the library
func Init() int {
	return int(C.wiringPiSetupInit())
}

// WinringPiPinMode : set pin mode
func PinMode(pin int, value int) {
	C.pinMode(C.int(pin), C.int(value))
}
