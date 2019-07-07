package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/xiaosumay/server-code-mgr/github"
	"github.com/xiaosumay/server-code-mgr/utils"
)

var (
	configPath  = flag.String("config", "", "Name of repo to create in authenticated user's GitHub account.")
	SecretToken string
	Version     string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println(Version)

	flag.Parse()

	http.HandleFunc("/", HandleFunc)

	for {
		utils.ParseConfig(*configPath)

		err := http.ListenAndServe("0.0.0.0:17293", nil)
		log.Println(err)

		time.Sleep(10 * time.Second)
	}
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

	mac := hmac.New(sha1.New, []byte(SecretToken))
	_, _ = mac.Write(data)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature[5:]), []byte(expectedMAC)) {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("签名不一致"))
		return
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
