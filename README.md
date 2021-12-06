![Doge-Assembly](https://socialify.git.ci/timwhitez/Doge-Assembly/image?description=1&font=Raleway&forks=1&issues=1&language=1&logo=https%3A%2F%2Favatars1.githubusercontent.com%2Fu%2F36320909&owner=1&pattern=Circuit%20Board&stargazers=1&theme=Light)

- 🐸Frog For Automatic Scan

- 🐶Doge For Defense Evasion&Offensive Security

# Doge-Assembly
Golang evasion tool, execute-assembly .Net file

## 20211206 加入非.net exe支持，fix bug
支持data文件夹内aeskey.exe与任意文件名exe加密文件识别

支持字符串混淆，详见main.go


## 使用go-donut重构 兼容性更好
原版见[old_version](./old_version)

## Intro
Are you still worrying about antivirus?


## feature
更新etw bypass相关代码，full dll unhooking相关代码，请重新获取依赖

go get -u github.com/timwhitez/Doge-Assembly

使用Golang execute assembly加载C#程序

C#程序编译为静态资源文件，使用AES加密，动态生成密钥

shellcode注入的过程采用direct syscall进行api调用

若想增强免杀效果可自行添加:
```
反沙箱反调试等相关代码

Blockdlls

parent-process-id-spoofing
```

### 写的较为仓促，希望能有大佬帮忙优化


## Usage
注意，若源程序需要多个参数执行，请使用如下方式:
```
in powershell:
Doge-Assembly.exe '-t schtask -c \"C:\Windows\System32\cmd.exe\" -a \"/c calc\" -n Test -m add -o hourly'

in cmd:
Doge-Assembly.exe -t schtask -c \"C:\Windows\System32\cmd.exe\" -a \"/c calc\" -n Test -m add -o hourly
```

```
cd encrypt
go build

you can change sharp.exe to other C# exe file

./encrypt.exe ./sharp.exe

copy version.txt to data/

copy aeskey.txt to data/

copy sharp.exe.cipher to data/

cd ..

go-bindata data/
go build

```

```
demo sharp.exe is SharpChromium.exe
```


## ref
资源文件加载:

- https://github.com/go-bindata/go-bindata

go-donut:

- https://github.com/Binject/go-donut

golang 的 execute assembly 实现:

- https://github.com/lesnuages/go-execute-assembly

bananaphone, golang hells gate:

- https://github.com/C-Sto/BananaPhone

etw bypass:

-https://blog.xpnsec.com/hiding-your-dotnet-etw/

-https://idiotc4t.com/defense-evasion/memory-pacth-bypass-etw


## todo



## screenshot
```
PS D:\Doge-Assembly> .\Doge-Assembly.exe
2021/03/29 17:08:14 Reloading c:\windows\system32\kernel32.dll...
2021/03/29 17:08:14 Made memory map RWX
2021/03/29 17:08:14 DLL overwritten
2021/03/29 17:08:14 Restored memory map permissions
2021/03/29 17:08:14 Reloading c:\windows\system32\kernelbase.dll...
2021/03/29 17:08:14 Made memory map RWX
2021/03/29 17:08:14 DLL overwritten
2021/03/29 17:08:14 Restored memory map permissions

All Dll Unhooked!

2021/03/29 17:08:14 Reloading c:\windows\system32\ntdll.dll...
2021/03/29 17:08:14 Made memory map RWX
2021/03/29 17:08:14 DLL overwritten
2021/03/29 17:08:14 Restored memory map permissions
Mess with the banana, die like the... banana?
patching .NET ETW ......
ETW patched!!

[X] Invalid argument passed:

Usage:
    .\SharpChromium.exe arg0 [arg1 arg2 ...]

Arguments:
    all       - Retrieve all Chromium Cookies, History and Logins.
    full      - The same as 'all'
    logins    - Retrieve all saved credentials that have non-empty passwords.
    history   - Retrieve user's history with a count of each time the URL was
                visited, along with cookies matching those items.
    cookies [domain1.com domain2.com] - Retrieve the user's cookies in JSON format.
                                        If domains are passed, then return only
                                        cookies matching those domains. Otherwise,
                                        all cookies are saved into a temp file of
                                        the format "%TEMP%\$browser-cookies.json"
```

# 🚀Star Trend
[![Stargazers over time](https://starchart.cc/timwhitez/Doge-Assembly.svg)](https://starchart.cc/timwhitez/Doge-Assembly)


## etc
1. 开源的样本大部分可能已经无法免杀,需要自行修改

2. 我认为基础核心代码的开源能够帮助想学习的人
 
3. 本人从github大佬项目中学到了很多
 
4. 若用本人项目去进行：HW演练/红蓝对抗/APT/黑产/恶意行为/违法行为/割韭菜，等行为，本人概不负责，也与本人无关

5. 本人已不参与大小HW活动的攻击方了，若溯源到timwhite id与本人无关
