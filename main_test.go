package main

import (
	"fmt"
	"sort"
	"testing"
)

type tmp struct {
	name  string
	best  string
	count int
	idx   int
}

type tmpSlice []tmp

func (t tmpSlice) Len() int {
	return len(t)
}

func (t tmpSlice) Less(i, j int) bool {
	ti := t[i]
	tj := t[j]
	if (ti.best != "" && tj.best != "") || (ti.best == "" && tj.best == "") {
		if ti.count == tj.count {
			return ti.idx < tj.idx
		}
		return ti.count > tj.count
	} else if ti.best != "" {
		if ti.count >= 2 || ti.count == tj.count {
			return true
		} else {
			return ti.count > tj.count
		}
	} else {
		if tj.count >= 2 || ti.count == tj.count {
			return false
		} else {
			return ti.count > tj.count
		}
	}
	//if t[i].pCount == t[j].pCount {
	//	lenI := len(t[i].best)
	//	lenJ := len(t[j].best)
	//	if lenI == 0 && lenJ == 0 {
	//		return t[i].idx < t[j].idx
	//	} else if lenI > 0 && lenJ > 0 {
	//		return t[i].idx < t[j].idx
	//	} else if lenI > 0 {
	//		return true
	//	} else {
	//		return false
	//	}
	//} else {
	//	return t[i].pCount > t[j].pCount
	//}
}

func (t tmpSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

type intSlice []int

func (is intSlice) Len() int {
	return len(is)
}

func (is intSlice) Less(i, j int) bool {
	fmt.Println(is, i, j)
	return is[i] < is[j]
}

func (is intSlice) Swap(i, j int) {
	is[i], is[j] = is[j], is[i]
}

func Test_itemCountSlice(t *testing.T) {
	t1 := tmpSlice{
		{"摩卡咖啡", "*", 1, 0},
		{"咖啡", "", 2, 1},
		{"卡布奇诺", "", 1, 2},
		{"香蕉", "", 1, 3},
		{"蓝色羽毛", "", 1, 4},
		{"松茸", "", 1, 5},
		{"紫阳花", "", 2, 6},
		{"祖母绿", "", 1, 7},
		{"含羞草沙拉", "", 3, 8},
		{"大米", "", 2, 9},
		{"番茄沙拉", "", 3, 10},
		{"土豆沙拉", "*", 2, 11},
		{"番茄汁", "", 2, 12},
		{"大豆", "", 1, 13},
		{"西瓜", "", 2, 14},
		{"蜜瓜", "", 1, 15},
		{"牛奶", "", 1, 16},
		{"鸡蛋", "", 1, 17},
		{"番茄", "", 4, 18},
		{"土豆", "", 3, 19},
		{"土豆泥", "*", 2, 20},
		{"草莓", "", 5, 21},
		{"红色的布", "", 1, 22},
		{"小萝ト", "", 2, 23},
		{"洋葱", "", 2, 24},
		{"小麦", "", 1, 25},
		{"虾虎鱼", "", 1, 26},
		{"红鳟", "", 1, 27},
		{"虹鳟", "", 1, 28},
		{"沙丁鱼", "", 2, 29},
		{"秋刀鱼", "", 2, 30},
		{"小龙虾", "", 1, 31},
		{"海胆", "", 2, 32},
		{"蚬贝", "", 1, 33},
		{"花蛤", "", 1, 34},
		{"扇贝", "", 1, 35},
		{"蒸面包果", "", 1, 36},
		{"淡紫色羽扇豆", "", 1, 37},
		{"热牛奶", "", 1, 38},
		{"粉钻", "", 1, 39},
		{"白金", "", 1, 40},
		{"向日葵", "", 1, 41},
		{"鸡蛋花", "", 1, 42},
		{"雏菊", "", 1, 43},
	}
	fmt.Println(t1)
	sort.Sort(t1)
	fmt.Println(t1)

	ti := intSlice{20, 1, 2, 40, 10, 90, 21, 25, 33, 22, 89}
	sort.Sort(ti)
	fmt.Println(ti)
}
