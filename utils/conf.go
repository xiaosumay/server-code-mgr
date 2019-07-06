package utils

import (
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

		Repositories[section.Name()] = *val
	}
	//当section空的时候的，一级配置不需要
	delete(Repositories, "DEFAULT")
}
