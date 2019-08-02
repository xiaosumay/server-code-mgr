# 基于github中转的代码管理工具


- 编译脚本参考
	```bat
	@echo off

	set scrip_path=%~dp0

	set work_path=%scrip_path%

	set GOROOT=C:\Go
	set GOPATH=%work_path%
	set GO111MODULE=on

	rem set ALL_PROXY=socks5://127.0.0.1:1080

	set Name=code-get

	set MAIN_VER=%DATE:~2,2%%DATE:~5,2%.%DATE:~8,2%
	set MINI_VER=%TIME:~0,2%
	set Version=%MAIN_VER: =0%%MINI_VER: =0%

	set FLAGS="-s -w -X main.Version=v%Version% -X main.SecretToken={your-github-webhook-secret}"

	pushd %work_path%src\github.com\xiaosumay\server-code-mgr

	set GOOS=windows
	"%GOROOT%\bin\go.exe" build -i -ldflags %FLAGS% -o %scrip_path%%Name%.exe .

	set GOOS=linux
	"%GOROOT%\bin\go.exe" build -i -ldflags %FLAGS% -o %scrip_path%%Name% .

	popd
	```
