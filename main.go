package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/timwhitez/Doge-Assembly/assembly"
)

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
