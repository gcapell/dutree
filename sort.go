package main

// Boiler-plate for sorting DisplayNodes

type DisplayNodeSlice []DisplayNode

func (p DisplayNodeSlice) Len() int {
	return len(p)
}

func (p DisplayNodeSlice) Less(i, j int) bool {
	return p[i].size > p[j].size
}

func (p DisplayNodeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
