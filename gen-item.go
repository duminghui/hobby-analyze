// go:build ignore
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

type kindMap struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items,flow"`
}

type gKindCustom struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items,flow"`
	Fix   []string `yaml:"fix,flow,omitempty"`
}

func (kc gKindCustom) ItemSlice() []string {
	if len(kc.Fix) > 0 {
		return kc.Fix
	}
	return kc.Items
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
	saveYamlFile(sKindCustomFile, kindCustomSlice)

	sKindMapFile := "./s-kind-map.yaml"
	var kindMapSlice []kindMap
	readYamlFile(sKindMapFile, &kindMapSlice)
	saveYamlFile(sKindMapFile, kindMapSlice)

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

	// 种类转换处理的Map
	// kindNameMap[所有,各种,xxx]=所有
	kindNameMap := make(map[string]string)
	kindNameMapCountMap := make(map[string]int)
	for _, km := range kindMapSlice {
		for _, name := range km.Items {
			kindNameMap[name] = km.Name
		}
		kindNameMapCountMap[km.Name] = len(km.Items)
	}

	// 这个里面的Key已经去重了.
	kcItemCountMap := make(map[string]int)
	for _, kc := range kindCustomSlice {
		addValue := 1
		if count, ok := kindNameMapCountMap[kc.Name]; ok {
			addValue = count
		}
		for _, itemName := range kc.Items {
			kcItemCountMap[itemName] = kcItemCountMap[itemName] + addValue
		}
	}

	//处理kind custom中的哪些物品可以使用
	//best better在person中出现的次数
	bItemInPCountMap := make(map[string]int)
	allItemInPCountMap := make(map[string]int)

	for _, p := range personSlice {
		bItemInPCountMap[p.Best]++
		allItemInPCountMap[p.Best]++
		for _, itemName := range p.Better {
			if rKcName, ok := kindNameMap[itemName]; ok {
				bItemInPCountMap[rKcName]++
				allItemInPCountMap[rKcName]++
			} else {
				bItemInPCountMap[itemName]++
				allItemInPCountMap[itemName]++
			}
		}
		for _, itemName := range p.Normal {
			if rKcName, ok := kindNameMap[itemName]; ok {
				allItemInPCountMap[rKcName]++
			} else {
				allItemInPCountMap[itemName]++
			}
		}
	}
	// 如果没在s-kind-item中出现, 则删除
	// 如果在person出现,则不删除
	// 如果没在person中出现, 查看他在kc中出现的次数, >1不删除
	// 如果=1,查看他所在的种类在person中的情况, 如果在better中出现不删除, 否则查看全部的出现情况, 如果出现次数>1不删除

	var delKcItemSlice []string
	processItemMap := make(map[string]int)
	for _, kc := range kindCustomSlice {
		for _, itemName := range kc.Items {
			if _, ok := processItemMap[itemName]; ok {
				continue
			}
			processItemMap[itemName] = 1
			if _, ok := allItemIdxMap[itemName]; !ok {
				delKcItemSlice = append(delKcItemSlice, itemName)
				continue
			}
			if _, ok := allItemInPCountMap[itemName]; ok {
				continue
			}
			if count := kcItemCountMap[itemName]; count > 1 {
				continue
			}
			if _, ok := bItemInPCountMap[kc.Name]; ok {
				continue
			}
			if count := allItemInPCountMap[kc.Name]; count > 1 {
				continue
			}
			delKcItemSlice = append(delKcItemSlice, itemName)
		}
	}
	// kind-custom的物品没在s-kind-item中
	var kcItemNotExistSlice []string
	kcItemMap := make(map[string]gKindCustom)
	for idx := range kindCustomSlice {
		kc := &kindCustomSlice[idx]
		var fix []string
		isFix := false
		for _, itemName := range kc.Items {
			if _, idxOk := allItemIdxMap[itemName]; !idxOk {
				kcItemNotExistSlice = appendSingle(kcItemNotExistSlice, itemName)
			}
			if !inStringArray(delKcItemSlice, itemName) {
				fix = append(fix, itemName)
			} else {
				isFix = true
			}
		}
		if isFix {
			kc.Fix = fix
		}
		kcItemMap[kc.Name] = *kc
	}

	// person的物品没在s-kind-item中
	var pItemNotExistSlice []string
	// person中使用的物品
	var pItemInfoSlice []*itemInfo
	pItemInfoMap := make(map[string]*itemInfo)

	// 从kind custom取物品
	itemSliceFromKindCustom := func(kcName string) []string {
		if rKcName, ok := kindNameMap[kcName]; ok {
			kcName = rKcName
		}
		if kc, ok := kcItemMap[kcName]; ok {
			return kc.ItemSlice()
		}
		return []string{}
	}

	genInfo := func(itemName string, personName string, flag string) {
		info, ok := pItemInfoMap[itemName]
		if !ok {
			idx, idxOk := allItemIdxMap[itemName]
			if !idxOk {
				pItemNotExistSlice = appendSingle(pItemNotExistSlice, itemName)
				idx = len(allItemIdxMap) + len(pItemNotExistSlice)
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
			pItemInfoSlice = append(pItemInfoSlice, info)
			pItemInfoMap[itemName] = info
		}
		if flag == "best" {
			info.BestFor = append(info.BestFor, personName)
		} else if flag == "b" {
			info.BetterFor = append(info.BetterFor, personName)
		} else if flag == "n" {
			info.NormalFor = append(info.NormalFor, personName)
		}
	}

	for pIdx, p := range personSlice {
		p.Idx = pIdx + 1
		if p.Best != "" {
			genInfo(p.Best, p.Name, "best")
		}
		isFix := false
		var betterFix []string
		for _, itemName := range p.Better {
			if kcItems := itemSliceFromKindCustom(itemName); len(kcItems) > 0 {
				betterFix = append(betterFix, kcItems...)
				isFix = true
			} else {
				betterFix = append(betterFix, itemName)
			}
		}
		if isFix {
			p.BetterFix = betterFix
		}
		for _, itemName := range betterFix {
			genInfo(itemName, p.Name, "b")
		}
		isFix = false
		var normalFix []string
		for _, itemName := range p.Normal {
			if kcItems := itemSliceFromKindCustom(itemName); len(kcItems) > 0 {
				normalFix = append(normalFix, kcItems...)
				isFix = true
			} else {
				normalFix = append(normalFix, itemName)
			}
		}
		if isFix {
			p.NormalFix = normalFix
		}
		for _, itemName := range normalFix {
			genInfo(itemName, p.Name, "n")
		}
	}

	sort.Slice(pItemInfoSlice, func(i, j int) bool {
		return pItemInfoSlice[i].Idx < pItemInfoSlice[j].Idx
	})
	for idx, item := range pItemInfoSlice {
		item.Idx = idx + 1
	}

	itemData := gItemData{}
	outItemFile := "./g-item.yaml"
	for _, item := range pItemInfoSlice {
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
		if item.Big == "B" {
			itemData.Big = append(itemData.Big, itemName)
		}
	}
	itemData.PNotExist = pItemNotExistSlice
	itemData.KCNotExist = kcItemNotExistSlice
	saveYamlFile(outItemFile, itemData)

	outKindCustomFile := "./g-kind-custom.yaml"
	saveYamlFile(outKindCustomFile, kindCustomSlice)

	outKindItemFile := "./g-kind-item.yaml"
	var gOutKindItemSlice []kindItem
	for _, ki := range kindItemSlice {
		var items []string
		var golds []string
		var big []string
		for _, itName := range ki.Items {
			itNameR := strings.Replace(itName, "*", "", -1)
			if _, ok := pItemInfoMap[itNameR]; ok {
				items = append(items, itName)
			}
		}
		for _, itName := range ki.Golds {
			if _, ok := pItemInfoMap[itName]; ok {
				golds = append(golds, itName)
			}
		}
		for _, itName := range ki.Big {
			if _, ok := pItemInfoMap[itName]; ok {
				big = append(big, itName)
			}
		}
		newKi := kindItem{
			KindLv1: ki.KindLv1,
			KindLv2: ki.KindLv2,
			Items:   items,
			Golds:   golds,
			Big:     big,
		}
		gOutKindItemSlice = append(gOutKindItemSlice, newKi)
	}
	saveYamlFile(outKindItemFile, gOutKindItemSlice)

	outItemListFile := "./g-item-list.yaml"
	saveYamlFile(outItemListFile, pItemInfoSlice)

	var itemInfoFlowSlice []itemInfoFlow
	for _, item := range pItemInfoSlice {
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
			if item, ok := pItemInfoMap[iName]; ok {
				iIdx = item.Idx
			}
			jIdx := 0
			if item, ok := pItemInfoMap[jName]; ok {
				jIdx = item.Idx
			}

			return iIdx < jIdx
		})
		sort.Slice(p.BetterFix, func(i, j int) bool {
			iName := p.BetterFix[i]
			jName := p.BetterFix[j]
			return pItemInfoMap[iName].Idx < pItemInfoMap[jName].Idx
		})
		sort.Slice(p.Normal, func(i, j int) bool {
			iName := p.Normal[i]
			jName := p.Normal[j]
			iIdx := 0
			if item, ok := pItemInfoMap[iName]; ok {
				iIdx = item.Idx
			}
			jIdx := 0
			if item, ok := pItemInfoMap[jName]; ok {
				jIdx = item.Idx
			}

			return iIdx < jIdx
		})
		sort.Slice(p.NormalFix, func(i, j int) bool {
			iName := p.NormalFix[i]
			jName := p.NormalFix[j]
			return pItemInfoMap[iName].Idx < pItemInfoMap[jName].Idx
		})
	}
	saveYamlFile(fixFile, personSlice)
	fmt.Println("s-kind-item重复出现的: ", duplicateNameSlice)
	fmt.Println("s-hobbies没在s-kind-item中", len(pItemNotExistSlice), pItemNotExistSlice)
	fmt.Println("s-kind-custom没在s-kind-item中", len(kcItemNotExistSlice), kcItemNotExistSlice)
	fmt.Println("s-kind-custom待删除的:", len(delKcItemSlice), delKcItemSlice)

}
