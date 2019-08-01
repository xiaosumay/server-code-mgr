package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/xiaosumay/server-code-mgr/github"
	"github.com/xiaosumay/server-code-mgr/utils"
)

var (
	SecretToken string
	Version     string

	configPath = flag.String("c", "/etc/code-get/repositories.conf", "配置文件")
	debug      = flag.Bool("debug", false, "调试模式，不验证token")
	update     = flag.Bool("u", false, "手动更新所有代码")
	token      = flag.String("token", SecretToken, "webhook的安全token")
	port       = flag.Int("port", 17293, "监听端口")
)

func main() {
	log.SetFlags(log.LstdFlags)
	log.Println(Version)

	flag.Parse()

	utils.Debug = *debug

	http.HandleFunc("/", HandleFunc)

	utils.ParseConfig(*configPath)

	if *update {
		for name, repo := range utils.Repositories {
			github.DoReposUpdate(name, repo)
		}
		return
	}

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), nil)
	log.Println(err)
}

func HandleFunc(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Println(err)
	}

	signature := request.Header.Get("X-Hub-Signature")

	if len(signature) == 0 {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("无效签名"))
		return
	}

	if !utils.Debug {
		mac := hmac.New(sha1.New, []byte(*token))
		_, _ = mac.Write(data)
		expectedMAC := hex.EncodeToString(mac.Sum(nil))

		if !hmac.Equal([]byte(signature[5:]), []byte(expectedMAC)) {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte("签名不一致"))
			return
		}
	}

	event := request.Header.Get("X-GitHub-Event")

	switch strings.ToLower(event) {
	case "ping":
		if github.PingEvent(data) {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("连接成功"))
			return
		}
	case "push":
		if github.PushEvent(data) {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("更新成功"))
			return
		}
	}

	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte("无效操作"))
}
