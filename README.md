# 基于github中转的代码管理工具


- build go-github-server
	```bat
	set GOROOT=C:\Go
	set GOPATH=C:\server-code-mgr;C:\Go
	C:\Go\bin\go.exe build -i -ldflags "-s -w -X main.Version=v3.0.0 -X main.SecretToken={your-github-webhook-secret}" -o go-github-server .
	```

- build repositories-mgr
	```bat
	set GOROOT=C:\Go
	set GOPATH=C:\server-code-mgr;C:\Go
	C:\Go\bin\go.exe build -i -ldflags "-s -w -X main.Token={your-github-token} -X main.version=v2.2.0 -X main.Owner={your-github-name} -X main.Secret={your-github-webhook-secret}" -o repositories-mgr repositories-mgr.go
	```
