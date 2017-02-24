package ioctl

import "github.com/yroffin/jarvis-go-ext/logger"

// IOCTL : ioctl wrapper
func IOCTL(fd, op, arg uintptr) error {
	logger.Default.Info("Ioctl [SIMULATED]", log.Fields{
		"fd":  fd,
		"op":  op,
		"arg": arg,
	})

	return nil
}
