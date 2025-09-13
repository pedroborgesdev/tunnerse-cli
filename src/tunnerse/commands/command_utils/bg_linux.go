//go:build linux || darwin
// +build linux darwin

package commands_utils

import "syscall"

func GetSysProcAttrForBackground() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}
