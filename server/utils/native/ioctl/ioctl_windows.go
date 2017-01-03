package ioctl

import log "github.com/Sirupsen/logrus"

// IOCTL : ioctl wrapper
func IOCTL(fd, op, arg uintptr) error {
	logrus.WithFields(log.Fields{
		"fd":  fd,
		"op":  op,
		"arg": arg,
	}).Info("Ioctl [SIMULATED]")

	return nil
}
