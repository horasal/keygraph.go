/*
KeyGraph visualize unit
Zhai Hongjie (c)2012
All rights reserved.
*/

package kgraph

import (
	"bufio"
	"container/heap"
	"github.com/ajstarks/svgo"
	"io"
	"log"
	"math"
	"strconv"
)

const rate = 3

func pointPosition(r *int, width, height, count int) (int, int) {
	count++
	rr := *r
	if 2*rr*(rate+1)*2*rr*(rate+1)*count > width*height {
		*r = int(math.Pow(float64(width*height/count), 0.5)) / (rate + 1) / 2
		rr = *r
	}
	rw := rr*(rate+1)*count%width + rr
	rh := rr * (rate + 1) * count / width
	return rw, (rh+1)*rr*(rate+1) + rr
}

func lengthBetween(x1, y1, x2, y2 int) int {
	return int(math.Pow(math.Pow(math.Abs(float64(x1-x2)), 2)+math.Pow(math.Abs(float64(y1-y2)), 2), 0.5)) + 1
}

func DrawGraph(graph BGraph, w io.Writer, width, height, r int) {
	if width < len(graph.Linestable)*5 {
		width = len(graph.Linestable) * 5
	}
	if height < len(graph.Linestable)*5 {
		height = len(graph.Linestable) * 5
	}
	log.Printf("save picture with size %d*%d", width, height)
	canvas := svg.New(w)
	canvas.Start(width, height)
	for i := 0; i < len(graph.Linestable); i++ {
		x, y := pointPosition(&r, width, height, i)
		canvas.Circle(x, y, r, "fill:none;stroke:black")
	}
	for i := 0; i < len(graph.Linestable); i++ {
		for j := i + 1; j < len(graph.Linestable[i]); j++ {
			if graph.Linestable[i][j] > 0 {
				x1, y1 := pointPosition(&r, width, height, i)
				x2, y2 := pointPosition(&r, width, height, j)
				canvas.Line(x1, y1, x2, y2, "fill:none;stroke:black")
			}
		}
	}
	canvas.End()
}

func Save(status wordStatus, propertfilter string, w *bufio.Writer) {
	itemheap := make(itemqueue, 0, len(status.Word))
	for v, l := range status.Word {
		heap.Push(&itemheap, &heapitem{Key: v, Value: l})
	}
	count := 0
	for itemheap.Len() > 0 {
		item := heap.Pop(&itemheap).(*heapitem)
		if item == nil {
			log.Println("unexpected item : nil")
			break
		}
		if status.properies[item.Key].prop == "" || propertfilter == "" || status.properies[item.Key].prop == propertfilter {
			w.WriteString(item.Key + "\t" + strconv.Itoa(item.Value) + "\n")
			count++
		}
	}
	w.Flush()
	log.Printf("%d item saved \n", count)
}
