////go:build ignore
//// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

type sBreedRule struct {
	PrefixKey []string         `yaml:"prefixKey,flow"`
	Exception []breedException `yaml:"exception"`
}

type breedException struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items,flow"`
}

type sItemData struct {
	Exists    []string `yaml:"exists,flow"`
	ExistsFix []string `yaml:"existsFix,flow"`
	Golds     []string `yaml:"golds,flow"`
	Big       []string `yaml:"big,flow"`
}

type gBreed struct {
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

func stringSliceItemNotInMap(src []string, dstMap map[string]int) ([]string, []string) {
	var notInSlice []string
	notInMap := make(map[string]int)
	for _, s := range src {
		if _, ok := dstMap[s]; !ok {
			notInSlice = append(notInSlice, s)
			notInMap[s] = 1
		}
	}
	var fixSlice []string
	for _, s := range src {
		if _, ok := notInMap[s]; !ok {
			fixSlice = append(fixSlice, s)
		}
	}
	return notInSlice, fixSlice
}

func inStringArray(arr []string, str string) bool {
	for _, s := range arr {
		if s == str {
			return true
		}
	}
	return false
}

func main() {

	sDataFile := "./s-hobbies.yaml"
	var personSlice []*person
	readYamlFile(sDataFile, &personSlice)

	// Format File
	//formatHobbiesFile := "./s-hobbies-f.yaml"
	//saveYamlFile(formatHobbiesFile, personSlice)
	saveYamlFile(sDataFile, personSlice)

	sItemFile := "./s-item.yaml"
	itemData := sItemData{}
	readYamlFile(sItemFile, &itemData)

	sBreedRuleFile := "./s-breed-rule.yaml"
	breedRule := sBreedRule{}
	readYamlFile(sBreedRuleFile, &breedRule)
	// Format File
	saveYamlFile(sBreedRuleFile, breedRule)

	sBreedFile := "./s-breed.yaml"
	var breedData []gBreed
	readYamlFile(sBreedFile, &breedData)
	// Format File
	saveYamlFile(sBreedFile, &breedData)

	breedNameSliceMap := make(map[string][]string)
	var breedNameSlice []string
	for _, b := range breedData {
		breedNameSliceMap[b.Name] = b.Items
		breedNameSlice = append(breedNameSlice, b.Name)
	}

	exceptionBreedMap := make(map[string]int)
	for _, ex := range breedRule.Exception {
		for _, iName := range ex.Items {
			key := fmt.Sprintf("%s:%s", ex.Name, iName)
			exceptionBreedMap[key] = 1
		}
	}

	// 获取所有物品信息及品种信息
	itemNameMap := make(map[string]int)

	for pIdx, p := range personSlice {
		p.Idx = pIdx + 1
		var itemNames []string
		if p.Best != "" {
			itemNames = append(itemNames, p.Best)
		}
		itemNames = append(itemNames, p.Better...)
		itemNames = append(itemNames, p.Normal...)
		for _, itemName := range itemNames {
			hasKey := false
			for _, key := range breedRule.PrefixKey {
				if strings.Contains(itemName, key) {
					hasKey = true
					break
				}
			}
			if hasKey {
				if _, ok := breedNameSliceMap[itemName]; !ok {
					breedNameSlice = append(breedNameSlice, itemName)
					breedNameSliceMap[itemName] = []string{}
				}
				continue
			}
			itemNameMap[itemName] = 1
		}
	}
	// 添加品种相关的物品, 如果不在存在列表中, 不添加
	itemExistMap := make(map[string]int)
	for _, itemName := range itemData.Exists {
		itemExistMap[itemName] = 1
	}

	for _, breed := range breedData {
		for _, itemName := range breed.Items {
			//if _, ok := itemExistMap[itemName]; ok {
			itemNameMap[itemName] = 1
			//}
		}
	}

	for itemName := range itemNameMap {
		if _, ok := breedNameSliceMap[itemName]; ok {
			continue
		}
		for _, breedName := range breedNameSlice {
			rBreedName := breedName
			for _, key := range breedRule.PrefixKey {
				rBreedName = strings.Replace(rBreedName, key, "所有", -1)
			}
			allIndex := strings.LastIndex(rBreedName, "所有")
			if strings.Contains(rBreedName, "所有") {
				rBreedName = rBreedName[allIndex+len("所有"):]
			}
			_, ok1 := exceptionBreedMap[fmt.Sprintf("%s:%s", rBreedName, itemName)]
			_, ok2 := exceptionBreedMap[fmt.Sprintf("%s:%s", breedName, itemName)]
			if ok1 || ok2 {
				continue
			}

			if rBreedName == "果汁" {
				rBreedName = "汁"
			}
			if strings.HasSuffix(itemName, rBreedName) {
				if !inStringArray(breedNameSliceMap[breedName], itemName) {
					breedNameSliceMap[breedName] = append(breedNameSliceMap[breedName], itemName)
				}
			}
		}
	}

	itemNameFixMap := make(map[string]int)
	for idx, p := range personSlice {
		if p.Best != "" {
			itemNameFixMap[p.Best] = 1
		}
		var betterFix []string
		for _, itemName := range p.Better {
			itemNames, _ := breedNameSliceMap[itemName]
			if len(itemNames) > 0 {
				for _, itemName1 := range itemNames {
					if !inStringArray(betterFix, itemName1) {
						betterFix = append(betterFix, itemName1)
					}
					itemNameFixMap[itemName1] = 1
				}
			} else {
				if !inStringArray(betterFix, itemName) {
					betterFix = append(betterFix, itemName)
				}
				itemNameFixMap[itemName] = 1
			}
		}
		personSlice[idx].BetterFix = betterFix

		var normalFix []string
		for _, itemName := range p.Normal {
			itemNames, _ := breedNameSliceMap[itemName]
			if len(itemNames) > 0 {
				for _, itemName1 := range itemNames {
					itemNameFixMap[itemName1] = 1
					if !inStringArray(betterFix, itemName1) && !inStringArray(normalFix, itemName1) {
						normalFix = append(normalFix, itemName1)
					}
				}
			} else {
				if !inStringArray(betterFix, itemName) && !inStringArray(normalFix, itemName) {
					normalFix = append(normalFix, itemName)
				}
				itemNameFixMap[itemName] = 1
			}
		}
		personSlice[idx].NormalFix = normalFix
	}

	_, itemExistSlice := stringSliceItemNotInMap(itemData.Exists, itemNameFixMap)
	_, golds := stringSliceItemNotInMap(itemData.Golds, itemNameFixMap)
	_, bigSlice := stringSliceItemNotInMap(itemData.Big, itemNameFixMap)
	itemData.ExistsFix = itemExistSlice
	itemData.Golds = golds
	itemData.Big = bigSlice

	existFixMap := make(map[string]int)
	goldMap := make(map[string]string)
	bigMap := make(map[string]string)
	for _, name := range itemData.ExistsFix {
		existFixMap[name] = 1
	}
	for _, name := range itemData.Golds {
		goldMap[name] = "G"
	}
	for _, name := range itemData.Big {
		bigMap[name] = "B"
	}
	idxMap := make(map[string]int)
	for idx, name := range itemData.Exists {
		name = strings.Replace(name, "*", "", -1)
		idxMap[name] = idx + 1
	}

	var itemInfoSlice []*itemInfo
	itemMap := make(map[string]*itemInfo)

	genNameList := func(itemNames []string, personName string, personBest string, addFlag string) {
		for _, itemName := range itemNames {
			info, ok := itemMap[itemName]
			if !ok {

				exitFlag := existFixMap[itemName]
				exist := "N"
				if exitFlag > 0 {
					exist = "E"
				}
				big := bigMap[itemName]
				gold := goldMap[itemName]
				idx := idxMap[itemName]
				info = &itemInfo{
					Name:  itemName,
					Idx:   idx,
					Exist: exist,
					Big:   big,
					Gold:  gold,
				}
				itemInfoSlice = append(itemInfoSlice, info)
				itemMap[itemName] = info
			}
			if personBest == itemName {
				info.BestFor = append(info.BestFor, personName)
			}
			if addFlag == "b" {
				info.BetterFor = append(info.BetterFor, personName)
			} else if addFlag == "n" {
				info.NormalFor = append(info.NormalFor, personName)
			}
		}
	}

	for pIdx, p := range personSlice {
		p.Idx = pIdx + 1
		if p.Best != "" {
			genNameList([]string{p.Best}, p.Name, p.Best, "")
		}
		genNameList(p.BetterFix, p.Name, p.Best, "b")
		genNameList(p.NormalFix, p.Name, p.Best, "n")
	}

	idxFix := len(idxMap) + 1
	for _, item := range itemInfoSlice {
		if item.Idx == 0 {
			item.Idx = idxFix
			idxFix++
		}
	}

	sort.Slice(itemInfoSlice, func(i int, j int) bool {
		return itemInfoSlice[i].Idx < itemInfoSlice[j].Idx
	})

	arrayJoin := func(array interface{}, seq string) string {
		return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", seq, -1)
	}

	var names []string
	var notExistNames []string
	count := 0
	for _, info := range itemInfoSlice {
		name := info.Name
		if info.Exist != "E" {
			count++
			name = name + "*"
			notExistNames = append(notExistNames, name)
		}
		names = append(names, name)
	}
	allItemNames := arrayJoin(names, ",")

	saveYamlFile("./itemNames.txt", allItemNames)

	notExistNameStr := arrayJoin(notExistNames, ",")
	saveYamlFile("./itemNames-no.txt", notExistNameStr)

	saveYamlFile(sItemFile, itemData)

	outItemFile := "./g-item.yaml"
	saveYamlFile(outItemFile, itemInfoSlice)

	var itemInfoFlowSlice []itemInfoFlow
	for _, item := range itemInfoSlice {
		info := itemInfoFlow{[]interface{}{item.Idx, item.Name, item.Exist, item.Gold, item.Big, item.BestFor, item.BetterFor, item.NormalFor}}
		itemInfoFlowSlice = append(itemInfoFlowSlice, info)
	}
	saveYamlFile("./g-item-flow.yaml", itemInfoFlowSlice)

	breedFile := "./g-breed.yaml"
	var breedSlice []gBreed
	for _, breedName := range breedNameSlice {
		breedItemSlice := breedNameSliceMap[breedName]
		//if len(breedItemSlice) > 0 {
		breedSlice = append(breedSlice, gBreed{breedName, breedItemSlice})
		//}
	}

	saveYamlFile(breedFile, breedSlice)

	fixFile := "./g-hobbies-fix.yaml"
	for _, p := range personSlice {
		sort.Slice(p.Better, func(i, j int) bool {
			iName := p.Better[i]
			jName := p.Better[j]
			iIdx := 0
			if item, ok := itemMap[iName]; ok {
				iIdx = item.Idx
			}
			jIdx := 0
			if item, ok := itemMap[jName]; ok {
				jIdx = item.Idx
			}

			return iIdx < jIdx
		})
		sort.Slice(p.BetterFix, func(i, j int) bool {
			iName := p.BetterFix[i]
			jName := p.BetterFix[j]
			return itemMap[iName].Idx < itemMap[jName].Idx
		})
		sort.Slice(p.Normal, func(i, j int) bool {
			iName := p.Normal[i]
			jName := p.Normal[j]
			iIdx := 0
			if item, ok := itemMap[iName]; ok {
				iIdx = item.Idx
			}
			jIdx := 0
			if item, ok := itemMap[jName]; ok {
				jIdx = item.Idx
			}

			return iIdx < jIdx
		})
		sort.Slice(p.NormalFix, func(i, j int) bool {
			iName := p.NormalFix[i]
			jName := p.NormalFix[j]
			return itemMap[iName].Idx < itemMap[jName].Idx
		})
	}
	saveYamlFile(fixFile, personSlice)

	var notInPerson []string
	for existName := range idxMap {
		if _, ok := itemMap[existName]; !ok {
			notInPerson = append(notInPerson, existName)
		}
	}
	fmt.Println("s-item存在,s-hobbies不存在", len(notInPerson), notInPerson)
	var notInExist []string
	for iName := range itemMap {
		if _, ok := idxMap[iName]; !ok {
			notInExist = append(notInExist, iName)
		}
	}
	fmt.Println("s-hobbies存在,s-item不存在", len(notInExist), notInExist)

}
