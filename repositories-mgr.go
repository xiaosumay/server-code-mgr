// Copyright 2018 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The newrepo command utilizes go-github as a cli tool for
// creating new repositoriesDo. It takes an auth Token as
// an enviroment variable and creates the new repo under
// the account affiliated with that Token.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/google/go-github/github"
	"github.com/jessevdk/go-flags"
	"golang.org/x/oauth2"
)

var (
	Owner   string
	Token   string
	version string
	Secret  string
)

type CommandOptions struct {
	Create bool `short:"a" long:"add" description:"新建"`
	List   bool `short:"l" long:"list" description:"查看"`
	Delete bool `short:"d" long:"delete" description:"删除"`
}

type ObjectOptions struct {
	Repositories  bool `long:"repo" description:"仓库"`
	Collaborators bool `long:"user" description:"合作者"`
	DeployKey     bool `long:"deploy" description:"服务器部署key"`
	WebHook       bool `long:"hook" description:"更新钩子"`
}

type CommonOptions struct {
	Token string `long:"token" description:"账户Token"`
	Owner string `long:"owner" description:"账户名"`
	Name  string `long:"name" description:"仓库名称"`
}

type RepoOptions struct {
	Description string `long:"desc" description:"描述"`
	Private     bool   `long:"private" description:"是否私有仓库"`
	Since       int64  `long:"since" description:"虽然有这个字段，但我也不知道干啥用的"`
}

type InviteOptions struct {
	GithubId string `long:"user-id" description:"被邀请人github的用户ID"`
}

type DevKeyOptions struct {
	Id       int64  `long:"deploy-id" description:"远程Key的唯一ID"`
	Key      string `long:"key" description:"ssh的public key 字符串或文件名"`
	Title    string `long:"title" description:"key的标题"`
	ReadOnly bool   `long:"read-only" description:"此key是否只读"`
}

type WebHookOptions struct {
	Secret string `long:"secret" description:"web hook Secret"`
	Ip     string `long:"ip" description:"更新钩子触发的服务器IP"`
	Id     int64  `long:"hook-id" description:"钩子的唯一ID"`
}

type Options struct {
	CommonOpts  CommonOptions  `group:"Common Options"`
	CommandOpts CommandOptions `group:"Command Options"`
	ObjectOpts  ObjectOptions  `group:"Object Options"`
	RepoOpts    RepoOptions    `group:"Repositories Options"`
	InviteOpts  InviteOptions  `group:"Invite Options"`
	DevKeyOpts  DevKeyOptions  `group:"DeployKey Options"`
	WebHookOpts WebHookOptions `group:"WebHook Options"`

	Version bool `long:"version" description:"版本信息"`
}

func DefaultValues(val string, fallback ...string) string {
	if val == "" && len(fallback) > 0 {
		return DefaultValues(fallback[0], fallback[1:]...)
	}

	return val
}

func fatalln(a ...interface{}) {
	fmt.Println(a...)
	os.Exit(1)
}

func main() {

	var opts Options
	var parser = flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	Token = DefaultValues(
		opts.CommonOpts.Token,
		os.Getenv("GITHUB_AUTH_TOKEN"),
		Token,
	)

	Owner = DefaultValues(
		opts.CommonOpts.Owner,
		os.Getenv("GITHUB_AUTH_OWNER"),
		Owner,
	)

	Secret = DefaultValues(
		opts.WebHookOpts.Secret,
		os.Getenv("GITHUB_AUTH_SECRET"),
		Secret,
	)

	if Token == "" || Owner == "" {
		fmt.Println("Unauthorized: No Token present")
		os.Exit(1)
	}

	if opts.ObjectOpts.Repositories {
		repositoriesDo(opts)
	} else if opts.ObjectOpts.Collaborators {
		collaborators(opts)
	} else if opts.ObjectOpts.DeployKey {
		deployKey(opts)
	} else if opts.ObjectOpts.WebHook {
		webHook(opts)
	} else {
		fatalln("你需要指定一个 Object Options")
	}
}

func getClient(token string) (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return client, ctx
}

func repositoriesDo(opts Options) {

	client, ctx := getClient(Token)

	if opts.CommandOpts.Create {
		options := opts.RepoOpts

		if opts.CommonOpts.Name == "" {
			fatalln("No name: New repos must be given a name")
		}

		if !options.Private {
			fmt.Println("必须开启 --private 此项才能是创建私有仓库")
		}

		r := &github.Repository{
			Name:        &opts.CommonOpts.Name,
			Private:     &options.Private,
			Description: &options.Description,
		}
		repo, _, err := client.Repositories.Create(ctx, "", r)
		if err != nil {
			fatalln(err)
		}

		fmt.Printf("Successfully created new repo: %v\n", repo.GetSSHURL())
	} else if opts.CommandOpts.List {
		reps, resp, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{
			Visibility: "all",
		})

		if err != nil || resp.StatusCode != 200 {
			fatalln(err)
		}

		for _, rep := range reps {
			fmt.Printf("%s: %s\n", rep.GetName(), rep.GetSSHURL())
		}
	} else if opts.CommandOpts.Delete {
		_, err := client.Repositories.Delete(ctx, Owner, opts.CommonOpts.Name)

		if err != nil {
			fatalln(err)
		}

		fmt.Println("删除" + opts.CommonOpts.Name + "成功!")
	} else {
		fatalln("你需要指定一个 Command Options")
	}
}

func collaborators(opts Options) {
	client, ctx := getClient(Token)

	if opts.CommandOpts.Create {
		invite := opts.InviteOpts

		if opts.CommonOpts.Name == "" {
			fatalln("No name: New repos must be given a name")
		}

		_, err := client.Repositories.AddCollaborator(ctx, Owner,
			opts.CommonOpts.Name,
			invite.GithubId,
			&github.RepositoryAddCollaboratorOptions{
				Permission: "push",
			})

		if err != nil {
			fatalln(err)
		}

		fmt.Println("https://github.com/MLTechMy/" + opts.CommonOpts.Name + "/invitations")
	} else if opts.CommandOpts.List {
		if opts.CommonOpts.Name == "" {
			fatalln("No name: New repos must be given a name")
		}

		users, _, err := client.Repositories.ListCollaborators(ctx, Owner,
			opts.CommonOpts.Name,
			&github.ListCollaboratorsOptions{
				Affiliation: "all",
			})

		if err != nil {
			fatalln(err)
		}

		for _, user := range users {
			fmt.Printf("%s: %s\n", user.GetLogin(), user.GetURL())
		}

	} else if opts.CommandOpts.Delete {
		if opts.CommonOpts.Name == "" {
			fatalln("No name: New repos must be given a name")
		}

		_, err := client.Repositories.RemoveCollaborator(ctx, Owner,
			opts.CommonOpts.Name,
			opts.InviteOpts.GithubId,
		)

		if err != nil {
			fatalln(err)
		}

		fmt.Println("删除合作者:" + opts.InviteOpts.GithubId)
	} else {
		fatalln("你需要指定一个 Command Options")
	}
}

func deployKey(opts Options) {
	client, ctx := getClient(Token)

	if opts.CommandOpts.Create {
		options := opts.DevKeyOpts

		if !options.ReadOnly {
			fmt.Println("必须开启 --read-only 此项才是最安全的")
			options.ReadOnly = true
		}

		if opts.CommonOpts.Name == "" {
			fatalln("No name: New repos must be given a name")
		}

		if _, ok := os.Stat(options.Key); ok == nil {
			tmpKey, err := ioutil.ReadFile(options.Key)
			if err == nil {
				options.Key = string(tmpKey)
			}
		}

		key, _, err := client.Repositories.CreateKey(ctx, Owner,
			opts.CommonOpts.Name,
			&github.Key{
				Key:      &options.Key,
				Title:    &options.Title,
				ReadOnly: &options.ReadOnly,
			})

		if err != nil {
			fatalln(err)
		}

		fmt.Printf("key[%d] %s install ok\n", + key.GetID(), key.GetTitle())
	} else if opts.CommandOpts.List {

		if opts.CommonOpts.Name == "" {
			fatalln("No name: New repos must be given a name")
		}

		keys, _, err := client.Repositories.ListKeys(ctx, Owner,
			opts.CommonOpts.Name,
			&github.ListOptions{},
		)

		if err != nil {
			fatalln(err)
		}

		for _, key := range keys {
			fmt.Printf("[%d]: %s\n", key.GetID(), key.GetTitle())
		}

	} else if opts.CommandOpts.Delete {
		if opts.CommonOpts.Name == "" {
			fatalln("No name: New repos must be given a name")
		}

		_, err := client.Repositories.DeleteKey(ctx, Owner, opts.CommonOpts.Name, opts.DevKeyOpts.Id)

		if err != nil {
			fatalln(err)
		}

		fmt.Println("删除DevKey成功！")
	} else {
		fatalln("你需要指定一个 Command Options")
	}
}

func webHook(opts Options) {
	client, ctx := getClient(Token)

	if opts.CommonOpts.Name == "" {
		fatalln("No name: New repos must be given a name")
	}

	if opts.CommandOpts.Create {
		active := true
		hookInfo := github.Hook{
			Config: map[string]interface{}{
				"url":          "http://" + opts.WebHookOpts.Ip + ":17293",
				"content_type": "json",
				"Secret":       Secret,
				"insecure_ssl": "0",
			},
			Events: []string{"push"},
			Active: &active,
		}
		hook, _, err := client.Repositories.CreateHook(ctx, Owner, opts.CommonOpts.Name, &hookInfo)

		if err != nil {
			fatalln(err)
		}

		fmt.Println("hook 安装: " + strconv.FormatBool(*hook.Active))
	} else if opts.CommandOpts.List {
		hooks, _, err := client.Repositories.ListHooks(ctx, Owner, opts.CommonOpts.Name, &github.ListOptions{})
		if err != nil {
			fatalln(err)
		}

		for _, hook := range hooks {
			fmt.Printf("%10d, %s\n", hook.GetID(), hook.Config["url"])
		}
	} else if opts.CommandOpts.Delete {
		_, err := client.Repositories.DeleteHook(ctx, Owner, opts.CommonOpts.Name, opts.WebHookOpts.Id)

		if err != nil {
			fatalln(err)
		}

		fmt.Println("删除hook成功!")
	} else {
		fatalln("你需要指定一个 Command Options")
	}
}
