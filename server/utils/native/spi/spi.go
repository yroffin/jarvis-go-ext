package spi

// extern int spiOpen();
import "C"

func SpiOpen() int {
	return int(C.spiOpen())
}
