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

func ReadConfFolder(conf_folder string) ([]*Conf, error) {
	conf_paths, err := filepath.Glob(filepath.Join(conf_folder, "*.ya?ml"))
	if err != nil {
		return nil, err
	}
	confs := make([]*Conf, len(conf_paths))
	for n, conf_path := range conf_paths {
		file, err := os.Open(conf_path)
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
