# Doge-Assembly
Golang evasion tools, use execute-assembly to excute C# tools

## feature
使用Golang execute assembly加载C#程序

C#程序编译进资源文件内，使用AES加密

进程注入的过程采用direct syscall

## Usage

```
cd encrypt
go build

you can change sharp.exe to your own C# exe file

./encrypt.exe ./sharp.exe

copy sharp.exe.cipher to bin/

cd ..
go-bindata data/
go build

```

```
demo sharp.exe is SharpChromium.exe
```
