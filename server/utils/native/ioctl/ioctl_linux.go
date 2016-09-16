package ioctl

import (
	"syscall"

	"github.com/yroffin/jarvis-go-ext/server/utils/logger"
)

// IOCTL : ioctl wrapper
func IOCTL(fd, op, arg uintptr) error {
	logger.NewLogger().WithFields(log.Fields{
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
