package ktfidf

import (
	"bufio"
	"container/heap"
	"log"
	"strconv"
)

type kfilter interface {
	Has(token string) bool
}

type kparser interface {
	Next() (string, error)
	IsEOS(token string) bool
	Reset()
}

type heapitem struct {
	Key   string
	Value float32
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

func CalculateTFIDF(token kparser, filter kfilter) map[string]float32 {
	tf := calculateTF(token, filter)
	token.Reset()
	idf := calculateIDF(token, filter)
	result := make(map[string]float32)
	for i, v := range tf {
		if _, ok := idf[i]; !ok {
			log.Printf("tfIDF:missing word %s.", i)
			continue
		}
		result[i] = float32(v) * idf[i]
	}
	log.Printf("tfIDF: total %d words.", len(result))
	return result
}

func CalculateRIDF(token kparser, filter kfilter) map[string]float32 {
	idf := calculateIDF(token, filter)
	token.Reset()
	poisson := calculatePoisson(token, filter)
	result := make(map[string]float32)
	for i, v := range idf {
		if _, ok := poisson[i]; !ok {
			log.Printf("RIDF:missing word: %s", i)
		}
		result[i] = v + poisson[i]
	}
	log.Printf("RIDF: total %d words.", len(result))
	return result
}

func Save(value map[string]float32, w *bufio.Writer) {
	itemheap := make(itemqueue, 0, len(value))
	for v, l := range value {
		heap.Push(&itemheap, &heapitem{Key: v, Value: l})
	}
	for itemheap.Len() > 0 {
		item := heap.Pop(&itemheap).(*heapitem)
		if item == nil {
			log.Println("unexpected item : nil")
			break
		}
		w.WriteString(item.Key + "\t" + strconv.FormatFloat(float64(item.Value), 'e', -1, 32) + "\n")
	}
	w.Flush()
	log.Printf("save %d items", len(value))
}
