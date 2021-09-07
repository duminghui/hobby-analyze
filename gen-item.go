//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type gItemData struct {
	Hobbies    []string `yaml:"hobbies,flow"`
	Have       []string `yaml:"have,flow"`
	Golds      []string `yaml:"golds,flow"`
	Big        []string `yaml:"big,flow"`
	PNotExist  []string `yaml:"pNotExist,flow"`
	KCNotExist []string `yaml:"kcNotExist,flow"`
}

type gKindCustom struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items,flow"`
}

func saveYamlFile(outFilePath string, in interface{}) {
	fixOutBytes, _ := yaml.Marshal(in)
	err := ioutil.WriteFile(outFilePath, fixOutBytes, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("保存文件:%s", outFilePath)
	fmt.Println()
}

func inStringArray(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func appendSingle(arr []string, str string) []string {
	if !inStringArray(arr, str) {
		arr = append(arr, str)
	}
	return arr
}

func stringJoin(array interface{}, seq string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", seq, -1)
}

func main() {

	sDataFile := "./s-hobbies.yaml"
	var personSlice []*person
	readYamlFile(sDataFile, &personSlice)
	saveYamlFile(sDataFile, personSlice)

	sKindItemFile := "./s-kind-item.yaml"
	var kindItemSlice []kindItem
	readYamlFile(sKindItemFile, &kindItemSlice)
	saveYamlFile(sKindItemFile, kindItemSlice)

	sKindCustomFile := "./s-kind-custom.yaml"
	var kindCustomSlice []gKindCustom
	readYamlFile(sKindCustomFile, &kindCustomSlice)
	saveYamlFile(sKindCustomFile, &kindCustomSlice)

	kcItemMap := make(map[string][]string)
	for _, kc := range kindCustomSlice {
		kcItemMap[kc.Name] = kc.Items
	}

	// 从所有物品中提取使用的物品列表,并提取出存在的物品列表,在使用的物品中存在所有物品中不存在的物品列表
	var duplicateNameSlice []string
	allItemIdxMap := make(map[string]int)
	allItemHaveMap := make(map[string]bool)

	goldItemMap := make(map[string]string)
	bigItemMap := make(map[string]string)

	_idx := 0
	for _, kind := range kindItemSlice {
		for _, itemName := range kind.Items {
			have := true
			if strings.HasSuffix(itemName, "*") {
				have = false
				itemName = strings.Replace(itemName, "*", "", -1)
			}

			allItemHaveMap[itemName] = have
			if _, ok := allItemIdxMap[itemName]; ok {
				duplicateNameSlice = append(duplicateNameSlice, itemName)
			} else {
				_idx++
				allItemIdxMap[itemName] = _idx
			}
		}
		for _, itemName := range kind.Golds {
			goldItemMap[itemName] = "G"
		}
		for _, itemName := range kind.Big {
			bigItemMap[itemName] = "B"
		}
	}

	// person的物品没在s-kind-item中
	var pItemNotExistSlice []string
	// person中使用的物品
	var gItemInfoSlice []*itemInfo
	gItemInfoMap := make(map[string]*itemInfo)
	haveIdxMap := make(map[string]int)

	genInfo := func(itemName string, personName string, flag string) {
		info, ok := gItemInfoMap[itemName]
		if !ok {
			idx, idxOk := allItemIdxMap[itemName]
			if !idxOk {
				pItemNotExistSlice = appendSingle(pItemNotExistSlice, itemName)
				idx = len(allItemIdxMap) + len(pItemNotExistSlice)
			} else {
				haveIdxMap[itemName] = 1
			}
			haveFlag := allItemHaveMap[itemName]
			have := "N"
			if haveFlag {
				have = "H"
			}
			gold := goldItemMap[itemName]
			big := bigItemMap[itemName]
			info = &itemInfo{
				Name: itemName,
				Idx:  idx,
				Have: have,
				Gold: gold,
				Big:  big,
			}
			gItemInfoSlice = append(gItemInfoSlice, info)
			gItemInfoMap[itemName] = info
		}
		if flag == "best" {
			info.BestFor = append(info.BestFor, personName)
		} else if flag == "b" {
			info.BetterFor = append(info.BetterFor, personName)
		} else if flag == "n" {
			info.NormalFor = append(info.NormalFor, personName)
		}
	}

	// kind-custom的物品没在s-kind-item中
	var kcItemNotExistSlice []string
	// 检查kind-custom的物品是否都在person中
	for pIdx, p := range personSlice {
		p.Idx = pIdx + 1
		if p.Best != "" {
			genInfo(p.Best, p.Name, "best")
		}
		var betterFix []string
		for _, itemName := range p.Better {
			if kcItems, ok := kcItemMap[itemName]; ok && len(kcItems) > 0 {
				for _, kcItemName := range kcItems {
					if _, idxOk := allItemIdxMap[kcItemName]; !idxOk {
						kcItemNotExistSlice = appendSingle(kcItemNotExistSlice, kcItemName)
					}
				}
				betterFix = append(betterFix, kcItems...)
			} else {
				betterFix = append(betterFix, itemName)
			}
		}
		p.BetterFix = betterFix
		for _, itemName := range betterFix {
			genInfo(itemName, p.Name, "b")
		}
		var normalFix []string
		for _, itemName := range p.Normal {
			if kcItems, ok := kcItemMap[itemName]; ok && len(kcItems) > 0 {
				for _, kcItemName := range kcItems {
					if _, idxOk := allItemIdxMap[kcItemName]; !idxOk {
						kcItemNotExistSlice = appendSingle(kcItemNotExistSlice, kcItemName)
					}
				}
				normalFix = append(normalFix, kcItems...)
			} else {
				normalFix = append(normalFix, itemName)
			}
		}
		p.NormalFix = normalFix
		for _, itemName := range normalFix {
			genInfo(itemName, p.Name, "n")
		}
	}

	sort.Slice(gItemInfoSlice, func(i, j int) bool {
		return gItemInfoSlice[i].Idx < gItemInfoSlice[j].Idx
	})
	for idx, item := range gItemInfoSlice {
		item.Idx = idx + 1
	}

	itemData := gItemData{}
	outItemFile := "./g-item.yaml"
	for _, item := range gItemInfoSlice {
		itemName := item.Name
		if item.Have != "H" {
			itemName += "*"
		} else {
			itemData.Have = append(itemData.Have, itemName)
		}
		itemData.Hobbies = append(itemData.Hobbies, itemName)
		if item.Gold == "G" {
			itemData.Golds = append(itemData.Golds, itemName)
		}
		if item.Big == "G" {
			itemData.Big = append(itemData.Big, itemName)
		}
	}
	itemData.PNotExist = pItemNotExistSlice
	itemData.KCNotExist = kcItemNotExistSlice
	saveYamlFile(outItemFile, itemData)

	outItemListFile := "./g-item-list.yaml"
	saveYamlFile(outItemListFile, gItemInfoSlice)

	var itemInfoFlowSlice []itemInfoFlow
	for _, item := range gItemInfoSlice {
		info := itemInfoFlow{[]interface{}{item.Idx, item.Name, item.Have, item.Gold, item.Big, item.BestFor, item.BetterFor, item.NormalFor}}
		itemInfoFlowSlice = append(itemInfoFlowSlice, info)
	}
	saveYamlFile("./g-item-flow.yaml", itemInfoFlowSlice)

	fixFile := "./g-hobbies-fix.yaml"
	for _, p := range personSlice {
		sort.Slice(p.Better, func(i, j int) bool {
			iName := p.Better[i]
			jName := p.Better[j]
			iIdx := 0
			if item, ok := gItemInfoMap[iName]; ok {
				iIdx = item.Idx
			}
			jIdx := 0
			if item, ok := gItemInfoMap[jName]; ok {
				jIdx = item.Idx
			}

			return iIdx < jIdx
		})
		sort.Slice(p.BetterFix, func(i, j int) bool {
			iName := p.BetterFix[i]
			jName := p.BetterFix[j]
			return gItemInfoMap[iName].Idx < gItemInfoMap[jName].Idx
		})
		sort.Slice(p.Normal, func(i, j int) bool {
			iName := p.Normal[i]
			jName := p.Normal[j]
			iIdx := 0
			if item, ok := gItemInfoMap[iName]; ok {
				iIdx = item.Idx
			}
			jIdx := 0
			if item, ok := gItemInfoMap[jName]; ok {
				jIdx = item.Idx
			}

			return iIdx < jIdx
		})
		sort.Slice(p.NormalFix, func(i, j int) bool {
			iName := p.NormalFix[i]
			jName := p.NormalFix[j]
			return gItemInfoMap[iName].Idx < gItemInfoMap[jName].Idx
		})
	}
	saveYamlFile(fixFile, personSlice)
	fmt.Println("s-kind-item重复出现的: ", duplicateNameSlice)
	fmt.Println("s-hobbies没在s-kind-item中", len(pItemNotExistSlice), pItemNotExistSlice)
	fmt.Println("s-kind-custom没在s-kind-item中", len(kcItemNotExistSlice), kcItemNotExistSlice)

}
