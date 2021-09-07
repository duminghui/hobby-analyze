package main

import (
	"fmt"
	"sort"
	"strings"
)

//go:generate go run gen-item.go data.go

type data struct {
	people []person
	items  []itemInfo
}

type itemCount struct {
	item   itemInfo
	weight string
	pCount int
	pNames []string
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
	icIW := len(icI.weight)
	if icI.weight == "N" {
		icIW = 0
	}
	icJW := len(icJ.weight)
	if icJ.weight == "N" {
		icJW = 0
	}

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
	bstICMap := make(map[string]*itemCount)
	betterICMap := make(map[string]*itemCount)
	betterICs := itemCountSlice{}
	normalICMap := make(map[string]*itemCount)
	normalICs := itemCountSlice{}
	for _, item := range d.items {
		if item.Have == "H" {
			weight := item.Gold + item.Big
			if weight == "" {
				weight = "N"
			}
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
	for _, p := range d.people {
		if !strings.HasSuffix(p.Name, "*") {
			existPeople = append(existPeople, p)
		} else {
			delete(bstICMap, p.Best)
			notExistNames = append(notExistNames, p.Name)
		}
	}
	delPNameSlice(betterICs, notExistNames)
	delPNameSlice(normalICs, notExistNames)
	fmt.Println("ExistPeople", existPeople, len(existPeople))
	resultPersonItemMap := make(map[string]string)
	resultItemPeopleMap := make(map[string][]string)

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
			for _, itemName := range p.BetterFix {
				if ic, ok2 := betterICMap[itemName]; ok2 {
					if ic.weight != "N" && ic.pCount > 2 {
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
					resultPersonItemMap[pName] = icBest.item.Name + "++"
					resultItemPeopleMap[icBest.item.Name] = append(resultItemPeopleMap[icBest.item.Name], pName+"++")
					delPNames = append(delPNames, pName)
				}
			}
			delPNameSlice(betterICs, delPNames)
			delPNameSlice(normalICs, delPNames)
			//delPNames = []string{}
			//if bOk {
			//	for _, pName := range icBetter.pNames {
			//		if _, ok2 := resultPersonItemMap[pName]; !ok2 {
			//			resultPersonItemMap[pName] = icBest.item.Name
			//			resultItemPeopleMap[icBest.item.Name] = append(resultItemPeopleMap[icBest.item.Name], pName)
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
			//			resultPersonItemMap[pName] = icBest.item.Name
			//			resultItemPeopleMap[icBest.item.Name] = append(resultItemPeopleMap[icBest.item.Name], pName)
			//			delPNames = append(delPNames, pName)
			//		}
			//	}
			//}
			//delPNameSlice(betterICs, delPNames)
			//delPNameSlice(normalICs, delPNames)
		}
	}
	fmt.Println("Best:", resultPersonItemMap, len(resultPersonItemMap))
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
				resultPersonItemMap[pName] = ic.item.Name + "+"
				resultItemPeopleMap[ic.item.Name] = append(resultItemPeopleMap[ic.item.Name], pName+"+")
				delPNames = append(delPNames, pName)
			}
		}
		delPNameSlice(betterICs, delPNames)
		delPNameSlice(normalICs, delPNames)
	}
	fmt.Println("Better:", resultPersonItemMap, len(resultPersonItemMap), iCount)
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
		var delPNames []string
		for _, pName := range ic.pNames {
			if _, ok := resultPersonItemMap[pName]; !ok {
				resultPersonItemMap[pName] = ic.item.Name
				resultItemPeopleMap[ic.item.Name] = append(resultItemPeopleMap[ic.item.Name], pName)
				delPNames = append(delPNames, pName)
			}
		}
		delPNameSlice(normalICs, delPNames)
	}
	fmt.Println("Normal:", resultPersonItemMap, len(resultPersonItemMap), iCount)
	fmt.Println("==============================")

	for idx, p := range existPeople {
		fmt.Printf("%d, %s: %s", idx+1, p.Name, resultPersonItemMap[p.Name])
		fmt.Println()
	}
	fmt.Println("=================")
	//
	idx := 0
	for _, item := range d.items {
		if pNames, ok := resultItemPeopleMap[item.Name]; ok {
			idx++
			fmt.Println(idx, item.Name, pNames)
		}
	}
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
