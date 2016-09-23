package wiringpi

// #include <stdio.h>
// void  delay             	(unsigned int howLong) {fprintf(stderr,"delay: %d\n", howLong);}
// void  delayMicroseconds 	(unsigned int howLong) {fprintf(stderr,"delayMicroseconds: %d\n", howLong);}
// unsigned int millis      (void) {}
// unsigned int micros      (void) {}
// void digitalWrite        (int pin, int value) {fprintf(stderr,"digitalWrite: %d,%d\n", pin, value);}
// int  wiringPiSetup       (void) {return 0;}
// void pinMode             (int pin, int mode) {}
// int  setuid      		(int uid) {return 0;}
import "C"
