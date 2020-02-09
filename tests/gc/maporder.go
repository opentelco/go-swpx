package main

import "sort"

type Plugin struct {
	n string
	w int
}

type byWeight []*Plugin

func (p byWeight) Less(i int, j int) bool { return p[i].Weight() < p[j].Weight() }
func (p byWeight) Len() int               { return len(p) }
func (p byWeight) Swap(i int, j int)      { p[i], p[j] = p[j], p[i] }

func (p Plugin) Weight() int {
	return p.w
}

func (p Plugin) Name() string {
	return p.n
}

func NewPlugin(w int, n string) *Plugin {
	return &Plugin{w: w, n: n}
}

func main() {
	providers := make([]*Plugin, 0)
	dp := NewPlugin(1, "default")
	sp := NewPlugin(20, "ssab")
	tp := NewPlugin(30, "telia")
	zp := NewPlugin(10, "zitius")
	providers = append(providers, dp, sp, tp, zp)

	sort.Sort(byWeight(providers))
	for _, b := range providers {
		println(b.Name(), b.Weight())
	}

}
