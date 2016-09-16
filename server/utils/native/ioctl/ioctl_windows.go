package ioctl

import (
	log "github.com/Sirupsen/logrus"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
)

// IOCTL : ioctl wrapper
func IOCTL(fd, op, arg uintptr) error {
	logger.NewLogger().WithFields(log.Fields{
		"fd":  fd,
		"op":  op,
		"arg": arg,
	}).Info("Ioctl [SIMULATED]")

	return nil
}
