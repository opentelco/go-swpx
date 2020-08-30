/*
 * Copyright (c) 2020. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

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
