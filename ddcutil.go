// Package main provides ...
package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	MOUSE_X = 1000
	MOUSE_Y = 1000
	M_POINT = ( MOUSE_X & 0xFFFFFFFF ) | ( MOUSE_Y << 32 )
	VPC_SOURCE_SELECT = '\x60'
	VGA1 = 1
	VGA2 = 2
	DVI1 = 3
	DVI2 = 4
	COMPOSITE_VIDEO1 = 5
	COMPOSITE_VIDEO2 = 6
	SVIDEO1 = 7
	SVIDEO2 = 8
	TUNER1 = 9
	TUNER2 = 10
	TUNER3 = 11
	COMPONENT_VIDEO1 = 12
	COMPONENT_VIDEO2 = 13
	COMPONENT_VIDEO3 = 14
	DISPLAY_PORT1 = 15
	DISPLAY_PORT2 = 16
	HDMI1 = 17
	HDMI2 = 18
	USBC = 27
)

var (
	user32, _ = syscall.LoadLibrary("User32.dll")
	dxva2, _ = syscall.LoadLibrary("dxva2.dll")
	monitorFromPoint, _ = syscall.GetProcAddress(user32, "MonitorFromPoint")
	GetPhysicalMonitorsFromHMONITOR, _ = syscall.GetProcAddress(dxva2, "GetPhysicalMonitorsFromHMONITOR")
	SetVCPFeature, _ = syscall.GetProcAddress(dxva2, "SetVCPFeature")
	DestroyPhysicalMonitor, _  = syscall.GetProcAddress(dxva2, "DestroyPhysicalMonitor")
)

func getMonitorHandle() (result uintptr) {
		var nargs uintptr = 2
		ret, _, callErr := syscall.Syscall(
			uintptr(monitorFromPoint),
			nargs,
			uintptr(M_POINT),
			uintptr(1),
			0)
		
		if callErr != 0 {
			abort("Call getMonitorHandle", callErr)	
		}
		
		res := uintptr(ret)
		result = getPhysicalMonitor(res)
		return
}

func getPhysicalMonitor(handle uintptr) (result uintptr) {
	b := make([]byte, 256)
	var nargs uintptr = 3
	_, _, callErr := syscall.Syscall(
		uintptr(GetPhysicalMonitorsFromHMONITOR),
		nargs,
		handle,
		uintptr(1),
		uintptr(unsafe.Pointer(&b[0])))
	
	if callErr != 0 {
			abort("Call getPhysicalMonitor", callErr)
	}
	result = uintptr(b[0])
	return
}

func setMonitorInputSource(source int) {
	mHandle := getMonitorHandle()
	
	var nargs uintptr = 3
	_, _, callErr := syscall.Syscall(
		uintptr(SetVCPFeature),
		nargs,
		uintptr(mHandle),
		VPC_SOURCE_SELECT,
		uintptr(source))
		
	if callErr != 0 {
			abort("Call SetVCPFeature", callErr)
	}
	destroyPhysicalMonitor(mHandle)
}

func destroyPhysicalMonitor(monitor uintptr) {
	var nargs uintptr = 1
	_, _, callErr := syscall.Syscall(
		uintptr(DestroyPhysicalMonitor),
		nargs,
		monitor,
		0,
		0)
		
	if callErr != 0 {
			abort("Call destroyPhysicalMonitor", callErr)
	}
}

func abort(funcname string, err error) {
    panic(fmt.Sprintf("%s failed: %v", funcname, err))
}



func main() {
	defer syscall.FreeLibrary(user32)
	defer syscall.FreeLibrary(dxva2)

	setMonitorInputSource(HDMI1)
}
