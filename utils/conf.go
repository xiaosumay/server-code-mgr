package utils

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type Repo struct {
	Path       string `ini:"path,omitempty"`
	Key        string `ini:"key,omitempty"`
	Script     string `ini:"script,omitempty"`
	Branch     string `ini:"branch,omitempty"`
	RemotePath string `ini:"remote_path,omitempty"`
}

func ParseConfig(configPath string) {
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalln("请提供配置文件")
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalln("配置文件出错2")
	}

	for _, section := range cfg.Sections() {
		val := new(Repo)

		err = section.MapTo(val)
		if err != nil {
			log.Fatalf("配置文件出错3: %v\n", err)
		}

		val.Path = DefaultValue(val.Path, fmt.Sprintf("/var/www/html/%s", section.Name()))
		val.RemotePath = DefaultValue(val.RemotePath, fmt.Sprintf("git@github.com/MLTechMy/%s.git", section.Name()))
		val.Branch = DefaultValue(val.Branch, "master")
		val.Key = DefaultValue(val.Key, fmt.Sprintf("/var/www/.ssh/%s", section.Name()))

		Repositories[section.Name()] = *val
	}
	//当section空的时候的，一级配置不需要
	delete(Repositories, "DEFAULT")
}
