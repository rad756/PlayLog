package logic

import (
	"sort"
	"strconv"
	"strings"
)

type Beta struct {
	Name     string
	Count    int
	SubCount int
	Finished bool
}

type BetaSlice struct {
	Slice []Beta
}

func NewBeta(Name string, Count int, SubCount int, Finished bool) *Beta {
	return &Beta{Name: Name, Count: Count, SubCount: SubCount, Finished: Finished}
}

func NewBetaSlice() *BetaSlice {
	return &BetaSlice{}
}

func (bs *BetaSlice) AddBeta(Name string, Count string, SubCount string, Finished bool, MyApp *MyApp, tabName string) *BetaSlice {
	CountInt, _ := strconv.Atoi(Count)
	SubCountInt, _ := strconv.Atoi(SubCount)
	bs.Slice = append(bs.Slice, *NewBeta(Name, CountInt, SubCountInt, Finished))

	bs.order()

	SaveBetaSlice(tabName, MyApp, *bs)

	return bs
}

func (bs *BetaSlice) order() *BetaSlice {
	sort.Slice(bs.Slice, func(i, j int) bool {
		return strings.ToLower(bs.Slice[i].Name) < strings.ToLower(bs.Slice[j].Name)
	})

	return bs
}

func (bs *BetaSlice) DeleteBeta(id int, MyApp *MyApp, tabName string) {
	if id >= 0 {
		bs.Slice = append(bs.Slice[:id], bs.Slice[id+1:]...)
	}

	SaveBetaSlice(tabName, MyApp, *bs)
}

func MoreBeta(s string) string {
	x, err := strconv.Atoi(s)

	if s == "" || err != nil {
		return "1"
	} else {
		return strconv.Itoa(x + 1)
	}
}

func LessBeta(s string) string {
	x, err := strconv.Atoi(s)

	if s == "" || err != nil {
		return "0"
	} else {
		return strconv.Itoa(x - 1)
	}
}

func (bs *BetaSlice) CountCurrent() int {
	x := 0

	for i := 0; i < len(bs.Slice); i++ {
		if !bs.Slice[i].Finished {
			x++
		}
	}

	return x
}

func (bs *BetaSlice) CountFinished() int {
	x := 0

	for i := 0; i < len(bs.Slice); i++ {
		if bs.Slice[i].Finished {
			x++
		}
	}

	return x
}
