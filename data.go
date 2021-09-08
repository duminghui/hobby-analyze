package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type kindItem struct {
	KindLv1 string   `yaml:"kindLv1"`
	KindLv2 string   `yaml:"kindLv2"`
	Items   []string `yaml:"items,flow"`
	Golds   []string `yaml:"golds,flow,omitempty"`
	Big     []string `yaml:"big,flow,omitempty"`
}

type itemInfoFlow struct {
	I []interface{} `yaml:"i,flow"`
}

type itemInfo struct {
	Name      string   `yaml:"name"`
	Idx       int      `yaml:"idx"`
	Have      string   `yaml:"have"`
	Gold      string   `yaml:"gold"`
	Big       string   `yaml:"big"`
	BestFor   []string `yaml:"bestFor,flow,omitempty"`
	BetterFor []string `yaml:"betterFor,flow"`
	NormalFor []string `yaml:"normalFor,flow"`
}

func (i itemInfo) String() string {
	weight := fmt.Sprintf("%s%s", i.Gold, i.Big)
	return fmt.Sprintf("{#%d,%s,%s,%s}", i.Idx, i.Name, i.Have, weight)
}

type person struct {
	Name      string   `yaml:"name"`
	Idx       int      `yaml:"idx,omitempty"`
	Best      string   `yaml:"best,omitempty"`
	Better    []string `yaml:"better,flow"`
	BetterFix []string `yaml:"betterFix,flow,omitempty"`
	Normal    []string `yaml:"normal,flow,omitempty"`
	NormalFix []string `yaml:"normalFix,flow,omitempty"`
}

func (p person) BetterSlice() []string {
	if len(p.BetterFix) > 0 {
		return p.BetterFix
	}
	return p.Better
}

func (p person) NormalSlice() []string {
	if len(p.NormalFix) > 0 {
		return p.NormalFix
	}
	return p.Normal
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
