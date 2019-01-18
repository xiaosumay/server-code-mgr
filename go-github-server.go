package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Repo struct {
	Path   string `ini:"path"`
	Key    string `ini:"key,omitempty"`
	Script string `ini:"cmd,omitempty"`
	Branch string `ini:"branch,omitempty"`
}

var (
	configPath   = flag.String("config", "", "Name of repo to create in authenticated user's GitHub account.")
	SecretToken  = []byte("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	repositories = make(map[string]Repo)
)

func main() {
	parseConfig()

	http.HandleFunc("/", HandleFunc)
	err := http.ListenAndServe("0.0.0.0:17293", nil)
	log.Fatal(err)
}

func parseConfig() {
	flag.Parse()

	if _, err := os.Stat(*configPath); err != nil {
		log.Fatal("请提供配置文件")
	}

	cfg, err := ini.Load(*configPath)
	if err != nil {
		log.Fatal("配置文件出错2")
	}

	for _, section := range cfg.Sections() {
		val := new(Repo)

		err = section.MapTo(val)
		if err != nil {
			log.Fatalf("配置文件出错3: %v", err)
		}

		repositories[section.Name()] = *val
	}
	//当section空的时候的，一级配置不需要
	delete(repositories, "DEFAULT")
}

func HandleFunc(writer http.ResponseWriter, request *http.Request) {
	data, err := ioutil.ReadAll(request.Body)

	if err != nil {
		log.Panicln(err)
	}

	signature := request.Header.Get("X-Hub-Signature")

	if len(signature) == 0 {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("无效签名"))
		return
	}

	mac := hmac.New(sha1.New, SecretToken)
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
		var ping PingPayload
		err = json.Unmarshal(data, &ping)
		if err == nil {
			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte("连接成功"))
			return
		}
	case "push":
		PushEvent(data, writer)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("无效操作"))
	}
}

func PushEvent(data []byte, writer http.ResponseWriter) {
	var push PushPayload
	err := json.Unmarshal(data, &push)
	if err != nil {
		log.Println(err)
		writer.Write([]byte("不可识别操作"))
		return
	}

	repoName := push.Repository.Name

	if repo, ok := repositories[repoName]; ok {

		ref := "refs/heads/" + DefaultValue(repo.Branch, "master")
		if push.Ref != ref {
			writer.Write([]byte("非指定分支，不更新"))
			return
		} else {
			go runCommand(repoName, repo)

			writer.Write([]byte("更新成功"))
		}
	} else {
		log.Println(repoName + " 不存在！")
		return
	}
}

func runCommand(repoName string, rep Repo) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	script := DefaultValue(rep.Script, dir+"/github_update.sh")

	if _, err := os.Stat(script); err != nil {
		log.Println(err)
		return
	}

	cmd := exec.Command("zsh", script)
	cmd.Env = append(cmd.Env, "BRANCH="+Quote(rep.Branch), "WORK_PATH="+Quote(rep.Path), "REPOS="+Quote(repoName))
	if 0 != len(rep.Key) {
		key := rep.Key
		if _, err := os.Stat(rep.Key); err != nil {
			key = strings.Join([]string{os.Getenv("HOME"), ".ssh", rep.Key}, string(os.PathSeparator))
		}

		cmd.Env = append(cmd.Env, "GIT_SSH_COMMAND=ssh -v -i "+Quote(key))
	}

	log.Println(strings.Join(cmd.Env, " "))

	data, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
	}

	log.Println(string(data))
}

