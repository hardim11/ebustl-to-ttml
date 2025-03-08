package ttmlgenerate

import "fmt"

type DebugMap struct {
	Total map[byte]int
}

func (m *DebugMap) AddElement(id byte) {
	if val, ok := m.Total[id]; ok {
		//do something here
		m.Total[id] = val + 1
	} else {
		m.Total[id] = 1
	}
}

func (m *DebugMap) ToString() string {
	res := "\tID\t\tCount\n"
	res = res + "       ========================\n"
	for id, count := range m.Total {
		res = res + fmt.Sprintf("\t%d\t=\t%d\n", id, count)
	}
	return res
}

func DebugMapNew() DebugMap {
	res := DebugMap{}
	res.Total = make(map[byte]int)
	return res
}
