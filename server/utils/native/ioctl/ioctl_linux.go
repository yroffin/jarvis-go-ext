package ioctl

import (
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
)

// IOCTL : ioctl wrapper
func IOCTL(fd, op, arg uintptr) error {
	logger.NewLogger().WithFields(logrus.Fields{
		"fd":  fd,
		"op":  op,
		"arg": arg,
	}).Info("Ioctl")

	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, fd, op, arg)
	if ep != 0 {
		return syscall.Errno(ep)
	}
	return nil
}
