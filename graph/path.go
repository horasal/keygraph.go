/*
KeyGraph path searching
Zhai Hongjie (c)2012
All rights reserved.
*/

package kgraph

type EdgeSearch struct {
	visitState []int
}

func (es *EdgeSearch) initialize(graph BGraph) {
	es.visitState = make([]int, len(graph.Linestable))
}

func (es *EdgeSearch) hasPath(graph BGraph, i, j int) bool {
	if graph.Linestable[i][j] != 0 || i == j {
		return true
	}
	for k := 0; k < len(graph.Linestable[i]); k++ {
		if graph.Linestable[i][k] != 0 && es.visitState[k] == 0 {
			es.visitState[k] = 1
			return es.hasPath(graph, k, j)
		}
	}
	return false
}

func (es *EdgeSearch) HasMultiPath(graph BGraph, i, j int) (result bool) {
	es.initialize(graph)
	if i == j || graph.Linestable[i][j] == 0 || !es.hasPath(graph, i, j) {
		return false
	}
	es.reset()
	graph.Linestable[i][j] = 0
	result = es.hasPath(graph, i, j)
	graph.Linestable[i][j] = 1
	return
}

func (es *EdgeSearch) HasPath(graph BGraph, i, j int) bool {
	es.initialize(graph)
	return es.hasPath(graph, i, j)
}

func (es *EdgeSearch) reset() {
	for i := 0; i < len(es.visitState); i++ {
		es.visitState[i] = 0
	}
}
