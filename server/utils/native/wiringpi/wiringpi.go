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

package wiringpi

// #cgo arm CFLAGS: -marm
// #cgo arm LDFLAGS: -lwiringPi
// extern int wiringPiSetupInit();
// extern void  delay             	(unsigned int howLong);
// extern void  delayMicroseconds 	(unsigned int howLong);
// extern void digitalWrite        (int pin, int value);
// extern int  wiringPiSetup       (void);
// extern int  wiringPiMode;
// extern void pinMode             (int pin, int mode);
// extern int  setuid      		(int uid);
import "C"
import (
	"sync"

	"github.com/Sirupsen/logrus"
)

// WiringPiDriver : wiring pi instance
type WiringPiDriver struct {
}

var instance *WiringPiDriver
var once sync.Once

// GetInstance : singleton instance
func GetInstance() *WiringPiDriver {
	once.Do(func() {
		instance = new(WiringPiDriver)
		instance.init()
	})
	return instance
}

// WiringPiInit : initialize the library
func (wiringPi *WiringPiDriver) init() int {
	var res = int(C.wiringPiSetupInit())
	logrus.WithFields(logrus.Fields{
		"Init": "on",
	}).Info("WiringPiDriver")
	return res
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
