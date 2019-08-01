package github

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	. "github.com/xiaosumay/server-code-mgr/utils"
	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
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

func CloneRepos(repoName string, rep Repo) {
	localPath := DefaultValue(rep.Path, "/var/www/html/"+repoName)

	log.Println(localPath)

	if _, err := os.Stat(localPath + "/.git"); err == nil {
		log.Printf("项目 %s 已存在\n", repoName)
		return
	}

	err := os.RemoveAll(localPath)
	if err != nil {
		log.Println(err)
	}

	auth, err := getAuth(rep.Key)
	if err != nil {
		log.Println(err)
		return
	}

	r, err := git.PlainClone(localPath, false, &git.CloneOptions{
		Auth:     auth,
		URL:      rep.RemotePath,
		Progress: os.Stdout,
		Tags:     git.AllTags,
	})

	if err != nil {
		log.Println(err)
		return
	}

	w, err := r.Worktree()
	if err != nil {
		log.Println(err)
		return
	}

	if rep.Branch != "master" {
		remoteRef, err := r.Reference(
			plumbing.NewRemoteReferenceName("origin", rep.Branch),
			true,
		)
		if err != nil {
			log.Println(err)
			return
		}

		err = w.Checkout(&git.CheckoutOptions{
			Hash:   remoteRef.Hash(),
			Branch: plumbing.NewBranchReferenceName(rep.Branch),
			Create: true,
			Keep:   true,
			Force:  true,
		})

		if err != nil {
			log.Println(err)
			return
		}
	}

	log.Println("下载成功！")
}

func runCommand(repoName string, rep Repo) bool {
	if _, err := os.Stat(rep.Script); err != nil {
		rep.Script = fmt.Sprintf("/var/www/.scripts/%s", rep.Script)
		if _, err := os.Stat(rep.Script); err != nil {
			return false
		}
	}

	cmd := exec.Command("bash", rep.Script)
	cmd.Env = append(cmd.Env, "BRANCH="+Quote(rep.Branch), "WORK_PATH="+Quote(rep.Path), "REPOS="+Quote(repoName))
	if 0 != len(rep.Key) {
		key := rep.Key
		if _, err := os.Stat(key); err != nil {
			key = fmt.Sprintf("/var/www/.ssh/%s", key)
		}

		cmd.Env = append(cmd.Env, "GIT_SSH_COMMAND=ssh -v -i "+Quote(key))
	}

	log.Println(strings.Join(cmd.Env, " "))

	data, err := cmd.CombinedOutput()

	log.Println(string(data))

	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func DoReposUpdate(repoName string, rep Repo) {
	for {
		if len(rep.Script) == 0 {
			break
		}

		log.Printf("启用自定义脚本: %s\n", rep.Script)

		if runCommand(repoName, rep) {
			return
		}

		break
	}

	if _, err := os.Stat(rep.Path + "/.git"); err != nil {
		CloneRepos(repoName, rep)
	}

	r, err := git.PlainOpenWithOptions(rep.Path, &git.PlainOpenOptions{
		DetectDotGit: false,
	})

	if err != nil {
		log.Println(err)
		return
	}

	auth, err := getAuth(rep.Key)
	if err != nil {
		log.Println(err)
		return
	}

	err = r.Fetch(&git.FetchOptions{
		Auth:     auth,
		Force:    true,
		Progress: os.Stdout,
		Tags:     git.AllTags,
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("强制拉去完成")

	remoteRef, err := r.Reference(
		plumbing.NewRemoteReferenceName("origin", rep.Branch),
		true,
	)
	if err != nil {
		log.Println(err)
		return
	}

	localRef, err := r.Reference(plumbing.NewBranchReferenceName(rep.Branch), true)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(remoteRef)
	log.Println(localRef)

	if remoteRef.Hash() == localRef.Hash() {
		log.Println("已经是最新的了！")
		return
	}

	w, err := r.Worktree()
	if err != nil {
		log.Println(err)
		return
	}

	err = w.Reset(&git.ResetOptions{
		Commit: remoteRef.Hash(),
		Mode:   git.HardReset,
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("更新完成")
}
