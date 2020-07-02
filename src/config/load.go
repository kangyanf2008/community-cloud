package config

import (
	"community-cloud/utils"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/zookeeper/v2"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func LoadConfigAndSetDefault() error {
	var configfile string
	var rootdir, vardir string

	flag.StringVar(&configfile, "conf", "", "conf file path")
	flag.StringVar(&rootdir, "root_dir", "", "root dir")
	flag.StringVar(&vardir, "var_dir", "", "var dir")

	flag.Parse()

	//community
	if configfile == "" {
		configfile = "/etc/community.toml"
	}

	//分环境
	env := os.Getenv("profiles.active")
	if env != "" {
		dir, fileName := filepath.Split(configfile)
		fileNameArr := strings.Split(fileName, ".")
		if len(fileNameArr) == 2 {
			tem := dir + fileNameArr[0] + "_" + env + "." + fileNameArr[1]
			if utils.FileExist(tem) {
				configfile = tem
			}
		}
	}
	var conf Config
	if _, err := toml.DecodeFile(configfile, &conf); err != nil {
		return errors.New("load toml conf file fail:" + err.Error())
	}

	//注册注册中心
	if len(conf.RegisterCenter.Address) > 0 {
		r := zookeeper.NewRegistry(func(op *registry.Options) {
			op.Addrs = conf.RegisterCenter.Address
			op.Context = context.Background()
			if conf.RegisterCenter.Timeout > 0 {
				op.Timeout = time.Second * time.Duration(conf.RegisterCenter.Timeout)
			} else {
				op.Timeout = time.Second * 5
			}
		})
		conf.RegisterCenter.register = r
	}

	if conf.Base.WebAddr == "" {
		conf.Base.WebAddr = ":8080"
	}
	if conf.Base.GRPCAddr == "" {
		conf.Base.GRPCAddr = ":8181"
	}
	if conf.Base.ServiceName == "" {
		conf.Base.ServiceName = conf.LogConf.Project
	}
	if rootdir != "" {
		conf.Base.RootDir = rootdir
	}
	if vardir != "" {
		conf.Base.VarDir = vardir
	}
	if conf.Base.VarDir == "" {
		conf.Base.VarDir = conf.Base.RootDir
	}

	conf.LogConf.LogDir = conf.Base.VarDir + "/" + conf.LogConf.LogDir

	fmt.Printf("conf is :%+v\n", conf)

	config = &conf

	return nil
}
