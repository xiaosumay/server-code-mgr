package github

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/xiaosumay/server-code-mgr/utils"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

func getAuth(key string) (*gitssh.PublicKeys, error) {

	if _, err := os.Stat(key); err != nil {
		key = fmt.Sprintf("/var/www/.ssh/%s", key)
	}

	priKey, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(priKey)
	if err != nil {
		return nil, err
	}

	auth := &gitssh.PublicKeys{
		User:   gitssh.DefaultUsername,
		Signer: signer,
	}
	auth.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	return auth, nil
}

func cloneRepos(repoName, url string, rep utils.Repo) {
	localPath := utils.DefaultValue(rep.Path, "/home/wwwroot/"+repoName)

	log.Println(localPath)

	if _, err := os.Stat(localPath); err == nil {
		log.Printf("项目 %s 已存在\n", repoName)
		return
	}

	auth, err := getAuth(rep.Key)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = git.PlainClone(localPath, false, &git.CloneOptions{
		Auth:     auth,
		URL:      url,
		Progress: os.Stdout,
		Tags:     git.AllTags,
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("下载成功！")
}
