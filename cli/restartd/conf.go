package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Conf struct {
	User     string   `yaml:"user"`
	Services []string `yaml:"services"`
}

func ReadConf(raw []byte, conf *Conf) error {
	err := yaml.Unmarshal(raw, conf)
	return err
}

func ReadConfFolder(confFolder string) ([]*Conf, error) {
	confPaths, err := filepath.Glob(filepath.Join(confFolder, "*.yml"))
	if err != nil {
		return nil, err
	}
	confs := make([]*Conf, len(confPaths))
	for n, confPath := range confPaths {
		file, err := os.Open(confPath)
		if err != nil {
			return nil, err
		}
		raw, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		conf := Conf{}
		err = ReadConf(raw, &conf)
		fmt.Println(conf)
		if err != nil {
			return nil, err
		}
		confs[n] = &conf
	}
	return confs, nil
}
