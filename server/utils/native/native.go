package native

// #cgo arm CFLAGS: -marm
// #cgo arm LDFLAGS: -lwiringPi
// extern int dioOn(int pin, int sender, int interruptor);
// extern int dioOff(int pin, int sender, int interruptor);
// extern int dioInit();
import "C"

func InitWiringPi() int {
	return int(C.dioInit())
}

/**
 * push ON
 */
func DioOn(pin int, sender int, interruptor int) int {
	return int(C.dioOn(C.int(pin), C.int(sender), C.int(interruptor)))
}

/**
 * push OFF
 */
func DioOff(pin int, sender int, interruptor int) int {
	return int(C.dioOff(C.int(pin), C.int(sender), C.int(interruptor)))
}
