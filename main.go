package main

import (
	"fmt"
	"sort"
	"strings"
)

//go:generate go run gen-item.go data.go utils.go

type itemResult struct {
	itemName string
	flag     string
}

func (i itemResult) String() string {
	return fmt.Sprintf("%s%s", i.itemName, i.flag)
}

type personResult struct {
	pName string
	flag  string
}

func (p personResult) String() string {
	return fmt.Sprintf("%s%s", p.pName, p.flag)
}

type data struct {
	people []person
	items  []itemInfo
}

type itemCount struct {
	item   itemInfo
	weight int
	pCount int
	pNames []string
}

type usedItemInfo struct {
	name   string
	pCount int
}

func (i itemCount) String() string {
	return fmt.Sprintf("{%s,%d}", i.item, i.pCount)
}

type itemCountSlice []*itemCount

func (ics itemCountSlice) Len() int {
	return len(ics)
}

func (ics itemCountSlice) Less(i, j int) bool {
	icI := ics[i]
	icJ := ics[j]
	if icI.pCount == 0 && icJ.pCount > 0 {
		return false
	} else if icI.pCount > 0 && icJ.pCount == 0 {
		return true
	} else if icI.pCount == 0 && icJ.pCount == 0 {
		return icI.item.Idx < icJ.item.Idx
	}
	icIW := icI.weight
	icJW := icJ.weight

	if icIW == icJW {
		if icI.pCount == icJ.pCount {
			return icI.item.Idx < icJ.item.Idx
		}
		return icI.pCount > icJ.pCount
	} else if icIW > icJW {
		return true
	} else {
		return false
	}
}

func (ics itemCountSlice) Swap(i, j int) {
	ics[i], ics[j] = ics[j], ics[i]
}

func analyze(d data) {
	itemInfoMap := make(map[string]itemInfo)
	bstICMap := make(map[string]*itemCount)
	betterICMap := make(map[string]*itemCount)
	betterICs := itemCountSlice{}
	normalICMap := make(map[string]*itemCount)
	normalICs := itemCountSlice{}
	for _, item := range d.items {
		if item.Have == "H" {
			itemInfoMap[item.Name] = item
			weight := len(item.Gold + item.Big)
			bestLen := len(item.BestFor)
			if bestLen > 0 {
				ic := &itemCount{
					item:   item,
					weight: weight,
					pCount: bestLen,
					pNames: item.BestFor,
				}
				bstICMap[item.Name] = ic
			}
			betterLen := len(item.BetterFor)
			if betterLen > 0 {
				ic := &itemCount{
					item:   item,
					weight: weight,
					pCount: betterLen,
					pNames: item.BetterFor,
				}
				betterICMap[item.Name] = ic
				betterICs = append(betterICs, ic)
			}
			normalLen := len(item.NormalFor)
			if normalLen > 0 {
				ic := &itemCount{
					item:   item,
					weight: weight,
					pCount: betterLen,
					pNames: item.NormalFor,
				}
				normalICMap[item.Name] = ic
				normalICs = append(normalICs, ic)
			}
		}
	}

	var existPeople []person
	var notExistNames []string
	personMap := make(map[string]person)
	for _, p := range d.people {
		if !strings.HasSuffix(p.Name, "*") {
			existPeople = append(existPeople, p)
			personMap[p.Name] = p
		} else {
			notExistNames = append(notExistNames, p.Name)
		}
	}
	delPNameSlice(betterICs, notExistNames)
	delPNameSlice(normalICs, notExistNames)
	lenExistPeople := len(existPeople)
	fmt.Println("NotExistPeople", notExistNames, lenExistPeople)
	fmt.Println("ExistPeople", existPeople, lenExistPeople)
	resultPersonItemMap := make(map[string]itemResult)
	resultItemPeopleMap := make(map[string][]personResult)

	//for _, ic := range bstICMap {
	//	fmt.Println("********")
	//	fmt.Println(ic.item.Name, ic.pNames)
	//	betterIC, ok := betterICMap[ic.item.Name]
	//	if ok {
	//		fmt.Println(betterIC.pCount, betterIC.pNames)
	//	} else {
	//		fmt.Println("Better Nil")
	//	}
	//}

	for _, p := range existPeople {
		icBest, ok := bstICMap[p.Best]
		if !ok {
			continue
		}
		useBest := true
		// icBest.weight一定是N
		icBetter, bOk := betterICMap[p.Best]
		// 如果较好的不存在, 或者较好的数量小于2, 比较该人的Better项, 如果better项中的权重高,并且数量高的, 就使用better里的.
		if !bOk || icBetter.pCount < 2 {
			for _, itemName := range p.BetterSlice() {
				if ic, ok2 := betterICMap[itemName]; ok2 {
					if ic.weight > 0 && ic.pCount > 2 {
						useBest = false
						break
					}
				}
			}
		}

		if useBest {
			var delPNames []string
			for _, pName := range icBest.pNames {
				if _, ok2 := resultPersonItemMap[pName]; !ok2 {
					itemName := icBest.item.Name
					resultPersonItemMap[pName] = itemResult{itemName, icBest.item.weight() + "++"}
					resultItemPeopleMap[itemName] = append(resultItemPeopleMap[itemName], personResult{pName, "++"})
					delPNames = append(delPNames, pName)
				}
			}
			delPNameSlice(betterICs, delPNames)
			delPNameSlice(normalICs, delPNames)
			//delPNames = []string{}
			//if bOk {
			//	for _, pName := range icBetter.pNames {
			//		if _, ok2 := resultPersonItemMap[pName]; !ok2 {
			//			resultPersonItemMap[pName] = itemResult{itemName, icBest.item.weight() + "+"}
			//			resultItemPeopleMap[itemName] = append(resultItemPeopleMap[itemName], pName+"+")
			//			delPNames = append(delPNames, pName)
			//		}
			//	}
			//}
			//delPNameSlice(betterICs, delPNames)
			//delPNameSlice(normalICs, delPNames)
			//delPNames = []string{}
			//if nIc, ok := normalICMap[p.Best]; ok {
			//	for _, pName := range nIc.pNames {
			//		if _, ok2 := resultPersonItemMap[pName]; !ok2 {
			//			resultPersonItemMap[pName] = itemResult{itemName, icBest.item.weight()}
			//			resultItemPeopleMap[itemName] = append(resultItemPeopleMap[itemName], pName)
			//			delPNames = append(delPNames, pName)
			//		}
			//	}
			//}
			//delPNameSlice(betterICs, delPNames)
			//delPNameSlice(normalICs, delPNames)
		}
	}
	fmt.Printf("#After Best People: %s %d/%d\n", resultPersonItemMap, len(resultPersonItemMap), lenExistPeople)
	fmt.Printf("#After Best Item: %s,%d\n", resultItemPeopleMap, len(resultItemPeopleMap))
	fmt.Println("--------------------------")
	iCount := 0
	for {
		sort.Sort(betterICs)
		if len(resultPersonItemMap) == len(existPeople) || len(betterICs) == 0 {
			break
		}
		iCount++
		ic := betterICs[0]
		if ic.pCount == 0 {
			break
		}
		var delPNames []string
		for _, pName := range ic.pNames {
			if _, ok := resultPersonItemMap[pName]; !ok {
				itemName := ic.item.Name
				resultPersonItemMap[pName] = itemResult{itemName, ic.item.weight() + "+"}
				resultItemPeopleMap[itemName] = append(resultItemPeopleMap[itemName], personResult{pName, "+"})
				delPNames = append(delPNames, pName)
			}
		}
		delPNameSlice(betterICs, delPNames)
		delPNameSlice(normalICs, delPNames)
	}
	fmt.Printf("#After Better People 1: %s %d/%d\n", resultPersonItemMap, len(resultPersonItemMap), lenExistPeople)
	fmt.Printf("#After Better Item 1: %s,%d\n", resultItemPeopleMap, len(resultItemPeopleMap))
	fmt.Println("--------------------------")

	var usedItemSlice []usedItemInfo
	for itemName, pNames := range resultItemPeopleMap {
		usedItemSlice = append(usedItemSlice, usedItemInfo{itemName, len(pNames)})
	}
	sort.Slice(usedItemSlice, func(i, j int) bool {
		return usedItemSlice[i].pCount > usedItemSlice[j].pCount
	})
	var delPNames []string
	for _, p := range existPeople {
		if _, ok := resultPersonItemMap[p.Name]; ok {
			continue
		}
		for _, usedItem := range usedItemSlice {
			itemName := usedItem.name
			if inStringArray(p.NormalSlice(), itemName) {
				_itemInfo := itemInfoMap[itemName]
				resultPersonItemMap[p.Name] = itemResult{itemName, _itemInfo.weight() + ""}
				resultItemPeopleMap[itemName] = append(resultItemPeopleMap[itemName], personResult{p.Name, ""})
				delPNames = append(delPNames, p.Name)
				break
			}
		}
	}
	delPNameSlice(normalICs, delPNames)
	fmt.Printf("#After Better People 2: %s %d/%d\n", resultPersonItemMap, len(resultPersonItemMap), lenExistPeople)
	fmt.Printf("#After Better Item 2: %s,%d\n", resultItemPeopleMap, len(resultItemPeopleMap))
	fmt.Println("--------------------------")
	iCount = 0
	for {
		iCount++
		sort.Sort(normalICs)
		if len(resultPersonItemMap) == len(existPeople) || len(normalICs) == 0 {
			break
		}
		ic := normalICs[0]
		if ic.pCount == 0 {
			break
		}
		var _delPNames []string
		for _, pName := range ic.pNames {
			if _, ok := resultPersonItemMap[pName]; !ok {
				itemName := ic.item.Name
				resultPersonItemMap[pName] = itemResult{itemName, ic.item.weight()}
				resultItemPeopleMap[itemName] = append(resultItemPeopleMap[itemName], personResult{pName, ""})
				_delPNames = append(_delPNames, pName)
			}
		}
		delPNameSlice(normalICs, _delPNames)
	}

	// 检查分配结果, 如果一个的可以再分配到其他地方, 就进行分配.
	var usedItemSliceEsc []usedItemInfo
	var usedItemSliceDesc []usedItemInfo
	for itemName, pNames := range resultItemPeopleMap {
		usedItemSliceEsc = append(usedItemSliceEsc, usedItemInfo{itemName, len(pNames)})
		usedItemSliceDesc = append(usedItemSliceDesc, usedItemInfo{itemName, len(pNames)})
	}
	sort.Slice(usedItemSliceEsc, func(i, j int) bool {
		return usedItemSliceEsc[i].pCount < usedItemSliceEsc[j].pCount
	})
	sort.Slice(usedItemSliceDesc, func(i, j int) bool {
		return usedItemSliceDesc[i].pCount > usedItemSliceDesc[j].pCount
	})
	for _, usedItemEsc := range usedItemSliceEsc {
		resultPersonSlice := resultItemPeopleMap[usedItemEsc.name]
		if len(resultPersonSlice) == 1 {
			rp := resultPersonSlice[0]
			if rp.flag == "++" {
				continue
			}
			pName := rp.pName
			p := personMap[pName]
			for _, usedItemDesc := range usedItemSliceDesc {
				toItemName := usedItemDesc.name
				if _, ok := resultItemPeopleMap[toItemName]; !ok {
					continue
				}
				if toItemName == usedItemEsc.name {
					continue
				}
				item := itemInfoMap[toItemName]
				if inStringArray(p.BetterSlice(), toItemName) {
					resultPersonItemMap[pName] = itemResult{item.Name, item.weight() + "+"}
					resultItemPeopleMap[toItemName] = append(resultItemPeopleMap[toItemName], personResult{pName, "+"})
					delete(resultItemPeopleMap, usedItemEsc.name)
					break
				}
				if rp.flag == "+" {
					continue
				}
				if inStringArray(p.NormalSlice(), toItemName) {
					resultPersonItemMap[pName] = itemResult{item.Name, item.weight()}
					resultItemPeopleMap[toItemName] = append(resultItemPeopleMap[toItemName], personResult{pName, ""})
					delete(resultItemPeopleMap, usedItemEsc.name)
					break
				}
			}
		}
	}

	//fmt.Println("After Normal:", resultPersonItemMap, len(resultPersonItemMap), iCount)
	//fmt.Println("After Normal:", resultItemPeopleMap, len(resultItemPeopleMap))
	fmt.Println("==============================")

	var peopleNoItemSlice []string
	for idx, p := range existPeople {
		ir, ok := resultPersonItemMap[p.Name]
		if ok {
			fmt.Printf("#%d %s: %s", idx+1, p.Name, ir)
		} else {
			fmt.Printf("#%d %s:", idx+1, p.Name)
			peopleNoItemSlice = append(peopleNoItemSlice, p.Name)
		}
		fmt.Println()
	}
	fmt.Println("=================")
	//
	idx := 0
	for _, item := range d.items {
		if pNames, ok := resultItemPeopleMap[item.Name]; ok {
			idx++
			fmt.Printf("#%d %s: %s\n", idx, item.Name+item.weight(), pNames)
		}
	}
	fmt.Println("没有分配物品的:", peopleNoItemSlice)
}

func delPNameSlice(ics itemCountSlice, pNames []string) {
	for _, ic := range ics {
		for _, delName := range pNames {
			ic.pNames = delPName(ic.pNames, delName)
		}
		ic.pCount = len(ic.pNames)
	}
}

func delPName(pNames []string, pName string) []string {
	var nPNames []string

	for _, name := range pNames {
		if name == pName {
			continue
		} else {
			nPNames = append(nPNames, name)
		}
	}
	return nPNames
	// 要实验 append 之后地址是否发生变化
	//idx := -1
	//for i, name := range pNames {
	//	if name == pName {
	//		idx = i
	//		break
	//	}
	//}
	//if idx < 0 {
	//	return pNames
	//}
	// 这种情况下append之后的地址不发生变化, 会发生数据乱的情况
	// 如 p1: [1,2,3,4,5]
	// p2 := p1
	// 从p2删除2,
	// p2: [1,3,4,5]
	// p1: [1,3,4,5,5]
	//pNames = append(pNames[:idx], pNames[idx+1:]...)
	//return pNames
}

func main() {
	genItems()
	fmt.Println("====================================")
	d := data{}
	hobbiesFile := "./g-hobbies-fix.yaml"
	readYamlFile(hobbiesFile, &(d.people))

	itemFile := "./g-item-list.yaml"
	readYamlFile(itemFile, &(d.items))

	//itemFile := "./g-item-flow.yaml"
	//var itemInfoFlowSlice []itemInfoFlow
	//readYamlFile(itemFile, &itemInfoFlowSlice)
	//for _, i := range itemInfoFlowSlice {
	//	fmt.Println(i)
	//}
	analyze(d)
}
