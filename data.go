package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type itemInfoFlow struct {
	I []interface{} `yaml:"i,flow"`
}

type itemInfo struct {
	Name      string   `yaml:"name"`
	Idx       int      `yaml:"idx"`
	Exist     string   `yaml:"exist"`
	Gold      string   `yaml:"gold"`
	Big       string   `yaml:"big"`
	BestFor   []string `yaml:"bestFor,flow,omitempty"`
	BetterFor []string `yaml:"betterFor,flow"`
	NormalFor []string `yaml:"normalFor,flow"`
}

func (i itemInfo) String() string {
	weight := fmt.Sprintf("%s%s", i.Gold, i.Big)
	return fmt.Sprintf("{#%d,%s,%s,%s}", i.Idx, i.Name, i.Exist, weight)
}

type person struct {
	Name      string   `yaml:"name"`
	Idx       int      `yaml:"idx,omitempty"`
	Best      string   `yaml:"best"`
	Better    []string `yaml:"better,flow"`
	BetterFix []string `yaml:"betterFix,flow,omitempty"`
	Normal    []string `yaml:"normal,flow,omitempty"`
	NormalFix []string `yaml:"normalFix,flow,omitempty"`
}

func (p person) String() string {
	return fmt.Sprintf("%d,%s", p.Idx, p.Name)
}

func readYamlFile(filePath string, dest interface{}) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, dest)
	if err != nil {
		panic(err)
	}
}
