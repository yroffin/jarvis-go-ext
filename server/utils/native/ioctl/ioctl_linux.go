package ioctl

// extern int ioctl_debug(void * pointer);
import "C"

import "syscall"

// IOCTL send ioctl
func IOCTL(fd, name, data uintptr) error {
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, fd, name, data)
	if err != 0 {
		return syscall.Errno(err)
	}
	//C.ioctl_debug(unsafe.Pointer(data))
	return nil
}
