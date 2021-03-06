module github.com/xiaosumay/server-code-mgr

go 1.12

require (
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gopkg.in/ini.v1 v1.44.0
	gopkg.in/src-d/go-git.v4 v4.12.0
)

replace golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 => github.com/golang/oauth2 v0.0.0-20190604053449-0f29369cfe45

replace golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4

replace golang.org/x/net v0.0.0-20190502183928-7f726cade0ab => github.com/golang/net v0.0.0-20190502183928-7f726cade0ab

replace golang.org/x/sys v0.0.0-20190422165155-953cdadca894 => github.com/golang/sys v0.0.0-20190422165155-953cdadca894
