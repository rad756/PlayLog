package logic

import (
	"sort"
	"strings"
)

type Alpha struct {
	Name string
	Kind string
}

type Kind struct {
	Slice []string
}

type AlphaSlice struct {
	Slice []Alpha
}

func NewAlpha(Name string, Kind string) *Alpha {
	return &Alpha{Name: Name, Kind: Kind}
}

func NewAlphaSlice() *AlphaSlice {
	return &AlphaSlice{}
}

func (as *AlphaSlice) AddAlpha(Name string, Kind string, MyApp MyApp, tabName string) *AlphaSlice {
	as.Slice = append(as.Slice, *NewAlpha(Name, Kind))

	as.order()

	SaveAlphaSlice(tabName, MyApp, *as)

	return as
}

func (as *AlphaSlice) DeleteAlpha(id int, MyApp MyApp, tabName string) {
	if id >= 0 {
		as.Slice = append(as.Slice[:id], as.Slice[id+1:]...)
	}

	SaveAlphaSlice(tabName, MyApp, *as)
}

func (as *AlphaSlice) order() *AlphaSlice {
	sort.Slice(as.Slice, func(i, j int) bool {
		return strings.ToLower(as.Slice[i].Name) < strings.ToLower(as.Slice[j].Name)
	})

	return as
}

func ReturnTestKind() *Kind {
	return &Kind{Slice: []string{"GB", "PC", "PS1", "Xbox"}}
}

func (k *Kind) AddKind(name string, fileName string, MyApp MyApp) *Kind {
	k.Slice = append(k.Slice, name)

	k.orderKind()

	SaveAlphaKind(fileName, MyApp, *k)

	return k
}

func (k *Kind) DeleteKind(id int, fileName string, MyApp MyApp) *Kind {
	if id >= 0 {
		k.Slice = append(k.Slice[:id], k.Slice[id+1:]...)
	}

	SaveAlphaKind(fileName, MyApp, *k)

	return k
}

func (k *Kind) orderKind() *Kind {
	sort.Slice(k.Slice, func(i, j int) bool {
		return strings.ToLower(k.Slice[i]) < strings.ToLower(k.Slice[j])
	})

	return k
}
