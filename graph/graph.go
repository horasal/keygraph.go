/*
KeyGraph prepare datas
Zhai Hongjie (c)2012
All rights reserved.
*/

package kgraph

import (
	"container/heap"
	"log"
	"parser"
)

type Filter interface {
	Has(token string) bool
	Add(token string)
	Initialize(language string)
}

type heapitem struct {
	Key   string
	Value int
	index int
}

type itemqueue []*heapitem

func (queue itemqueue) Len() int { return len(queue) }

func (queue itemqueue) Less(i, j int) bool {
	return queue[i].Value > queue[j].Value
}

func (queue itemqueue) Swap(i, j int) {
	queue[i], queue[j] = queue[j], queue[i]
	queue[i].index = i
	queue[j].index = j
}

func (queue *itemqueue) Push(x interface{}) {
	a := *queue
	n := len(a)
	a = a[0 : n+1]
	item := x.(*heapitem)
	item.index = n
	a[n] = item
	*queue = a
}

func (pq *itemqueue) Pop() interface{} {
	a := *pq
	n := len(a)
	if n < 1 {
		return nil
	}
	item := a[n-1]
	item.index = -1
	*pq = a[0 : n-1]
	return item
}

type propery struct {
	prop string
}

type wordStatus struct {
	Word      map[string]int
	Emotion   map[string]float32
	properies map[string]propery
}

func GetPeopry(token kparser.Token, filter Filter, wordlist *wordStatus) {
	wordlist.properies = make(map[string]propery)
	for word, err := token.Next(); err == nil; word, err = token.Next() {
		if filter.Has(word) || token.IsEOS(word) {
			continue
		}
		if _, ok := wordlist.properies[word]; !ok {
			wordlist.properies[word] = propery{prop: token.Propery()}
		}
	}
}

func getmaxnum(maxnum, wordcount int) int {
	if maxnum < 1 {
		maxnum = wordcount / 1000
	}
	if maxnum < 15 {
		maxnum = 15
	}
	if wordcount < maxnum {
		maxnum = wordcount / 3
	}
	return maxnum
}

func GenHighFreq(token kparser.Token, filter Filter, maxnum int) map[string]int {
	var (
		wordlist map[string]int = make(map[string]int)
		highfreq map[string]int = make(map[string]int)
	)
	for word, err := token.Next(); err == nil; word, err = token.Next() {
		if filter.Has(word) || token.IsEOS(word) {
			continue
		}
		if _, ok := wordlist[word]; ok {
			wordlist[word]++
		} else {
			wordlist[word] = 0
		}
	}
	maxnum = getmaxnum(maxnum, len(wordlist))
	itemheap := make(itemqueue, 0, len(wordlist))
	for v, l := range wordlist {
		heap.Push(&itemheap, &heapitem{Key: v, Value: l})
	}
	for i := 0; i < maxnum; i++ {
		item := heap.Pop(&itemheap).(*heapitem)
		if item == nil {
			break
		}
		highfreq[item.Key] = item.Value
	}
	log.Printf("total %d words,%d words selected.\n", len(wordlist), maxnum)
	return highfreq
}

func SumCo(token kparser.Token, filter Filter, graph BGraph) wordStatus {
	wordlist := make(map[string]int)
	wordemotion := make(map[string]float32)
	search := new(EdgeSearch)
	for word, err := token.Next(); err == nil; word, err = token.Next() {
		sentword := make([]string, 0)
		keyword := make([]string, 0)
		sent := make(map[string]int)
		key := make(map[string]int)
		for ; !token.IsEOS(word) && err == nil; word, err = token.Next() {
			if filter.Has(word) || token.IsEOS(word) {
				continue
			}
			if _, ok := graph.Nametable[word]; ok {
				has := false
				for i := 0; i < len(keyword); i++ {
					if search.HasPath(graph, graph.Nametable[word], graph.Nametable[keyword[i]]) {
						has = true
						key[keyword[i]]++
						break
					}
				}
				if !has {
					keyword = append(keyword, word)
					key[word] = 1
				}
			} else {
				sentword = append(sentword, word)
				if _, ok := sent[word]; ok {
					sent[word]++
				} else {
					sent[word] = 1
				}
			}
		}
		if len(keyword) > 0 {
			for i := 0; i < len(sentword); i++ {
				if _, ok := wordlist[sentword[i]]; !ok {
					wordlist[sentword[i]] = len(keyword) * sent[sentword[i]]
				} else {
					wordlist[sentword[i]] += len(keyword) * sent[sentword[i]]
				}
			}
			for i := 0; i < len(keyword); i++ {
				if _, ok := wordlist[keyword[i]]; !ok {
					wordlist[keyword[i]] = (len(keyword) - 1) * key[keyword[i]]
				} else {
					wordlist[keyword[i]] += (len(keyword) - 1) * key[keyword[i]]
				}
			}
		}
	}
	log.Printf("words in result: %d", len(wordlist))
	return wordStatus{Word: wordlist, Emotion: wordemotion}
}
