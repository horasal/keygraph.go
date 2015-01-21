/*
KeyGraph process graph
Zhai Hongjie (c)2012
All rights reserved.
*/

package kgraph

import (
	"log"
	"parser"
)

type Emotion interface {
	Emotion(token string) float32
	Initialize(identifier string)
	Add(token string, emot float32)
}

type BGraph struct {
	Nametable    map[string]int
	Linestable   [][]int
	Emotiontable map[string]float32
}

func BuildGraph(baseWords map[string]int, token kparser.Token) *BGraph {
	graph := new(BGraph)
	index := 0
	nametable := make(map[string]int)
	graph.Nametable = nametable
	for i, _ := range baseWords {
		if _, ok := graph.Nametable[i]; !ok {
			graph.Nametable[i] = index
			index++
		}
	}
	log.Printf("build graph for total %d words", index)
	linestable := make([][]int, index)
	for i := range linestable {
		linestable[i] = make([]int, index)
	}
	graph.Linestable = linestable
	for word, err := token.Next(); err == nil; word, err = token.Next() {
		senttable := make(map[string]int)
		for ; !token.IsEOS(word) && err == nil; word, err = token.Next() {
			if _, ok := graph.Nametable[word]; ok {
				if _, ok := graph.Nametable[word]; ok {
					senttable[word]++
				} else {
					senttable[word] = 1
				}
			}
		}
		for v, i := range senttable {
			for u, j := range senttable {
				if i == j {
					continue
				}
				graph.Linestable[graph.Nametable[v]][graph.Nametable[u]] += i * j
			}
		}
	}
	return graph
}

func abs(i float32) float32 {
	if i < 0 {
		return -i
	}
	return i
}

func InitializeBaseEmotion(graph *BGraph, baseWord map[string]int, emotion Emotion, token kparser.Token) {
	if emotion == nil {
		return
	}
	wordcount := make(map[string]float32)
	graph.Emotiontable = make(map[string]float32)
	for v, _ := range baseWord {
		wordcount[v] = 0
	}
	for word, err := token.Next(); err == nil; word, err = token.Next() {
		senttable := make(map[string]int)
		var sentemotion float32 = 0.0
		for ; !token.IsEOS(word) && err == nil; word, err = token.Next() {
			if _, ok := baseWord[word]; ok {
				senttable[word] = 1
			}
			sentemotion += emotion.Emotion(word)
		}
		for v, _ := range senttable {
			wordcount[v] += float32(sentemotion)
		}
	}
	var k float32 = 0.0
	for _, i := range wordcount {
		if k < abs(i) {
			k = abs(i)
		}
	}
	if k == 0 {
		k = 1
	}
	for v, i := range wordcount {
		graph.Emotiontable[v] = i / k
	}
}

func SelectEdge(graph *BGraph, emotion Emotion, loopdepth, threshold int) {

	if emotion != nil {
		if loopdepth < 3 {
			loopdepth = 3
		}
		log.Printf("iterative %d times for emotion", loopdepth)
		for i := 0; i < loopdepth; i++ {
			for j, _ := range graph.Emotiontable {
				var sumEmotion float32 = 0
				var sidesum int = 0
				for k := 0; k < len(graph.Linestable[graph.Nametable[j]]); k++ {
					sidesum += graph.Linestable[graph.Nametable[j]][k]
				}
				for k, _ := range graph.Emotiontable {
					sumEmotion += graph.Emotiontable[k] * float32(graph.Linestable[graph.Nametable[j]][graph.Nametable[k]]) / float32(sidesum)
				}
				graph.Emotiontable[j] = sumEmotion
			}
		}
	}
	if threshold < 0 {
		threshold = 0
		max := 0
		for i := 0; i < len(graph.Linestable); i++ {
			for j := 0; j < len(graph.Linestable[i]); j++ {
				if graph.Linestable[i][j] > max {
					max = graph.Linestable[i][j]
				}
				threshold += graph.Linestable[i][j]
			}
		}
		threshold /= len(graph.Linestable) * len(graph.Linestable)
		if threshold == 0 {
			threshold = 100
		}
		if threshold > max {
			threshold = max
		} else {
			threshold = threshold / 2
		}
	}
	for i := 0; i < len(graph.Linestable); i++ {
		for j := 0; j < len(graph.Linestable[i]); j++ {
			if graph.Linestable[i][j] > threshold {
				graph.Linestable[i][j] = 1
			} else {
				graph.Linestable[i][j] = 0
			}
		}
	}
	search := new(EdgeSearch)
	k := 0
	for i := 0; i < len(graph.Linestable); i++ {
		for j := 0; j < len(graph.Linestable[i]); j++ {
			if graph.Linestable[i][j] != 0 && !search.HasMultiPath(*graph, i, j) {
				graph.Linestable[i][j] = 0
				k++
			}
		}
	}
	log.Printf("delete %d edges", k)
}
