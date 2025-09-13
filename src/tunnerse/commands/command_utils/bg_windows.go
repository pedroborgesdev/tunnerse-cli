//go:build windows
// +build windows

package commands_utils

import "syscall"

func GetSysProcAttrForBackground() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
