package kmeans

import (
	"log"
	"math"
	"math/rand"
	"time"
)

type datas interface {
	Distance(i, j int) float64
	Len() int
}

type Class struct {
	Center  int
	Classes []int
}

func classify(data datas, Classes []*Class) {
	for i := 0; i < len(Classes); i++ {
		Classes[i].Classes = []int{}
	}
	for i := 0; i < data.Len(); i++ {
		min := 0
		for dis, j := math.Inf(1), 0; j < len(Classes); j++ {
			if t := math.Abs(data.Distance(i, Classes[j].Center)); dis > t {
				dis = t
				min = j
			}
		}
		Classes[min].Classes = append(Classes[min].Classes, i)
	}
}

func center(data datas, Classes []*Class) bool {
	changed := false
	for j := 0; j < len(Classes); j++ {
		k := 0.0
		for i := 0; i < len(Classes[j].Classes); i++ {
			k += data.Distance(Classes[j].Center, Classes[j].Classes[i])
		}
		if len(Classes[j].Classes) == 0 {
			k = 0
		} else {
			k = k / float64(len(Classes[j].Classes))
		}
		c := Classes[j].Center
		length := math.Abs(k)
		for i := 0; i < len(Classes[j].Classes); i++ {
			if t := math.Abs(k - data.Distance(c, Classes[j].Classes[i])); t < length {
				Classes[j].Center = Classes[j].Classes[i]
				length = t
			}
		}
		if c != Classes[j].Center && math.Abs(data.Distance(c, Classes[j].Center)) > math.Abs(data.Distance(c, c)) {
			changed = changed || true
			//log.Printf("class %d : %d to %d", j, c, Classes[j].Center)
		} else {
			changed = changed || false
		}
	}
	return changed
}

func randomize() {
	rand.Seed(time.Second.Nanoseconds())
}

func Kmeans(data datas, k, maxloop int) []*Class {
	if k > data.Len() {
		log.Printf("Can not classify %d in %d Classes\n", data.Len(), k)
		return nil
	}
	resc := make([]*Class, k)
	for i := 0; i < k; i++ {
		resc[i] = new(Class)
		resc[i].Center = i
		resc[i].Classes = make([]int, 1)
		resc[i].Classes[0] = i
	}
	if k == data.Len() {
		return resc
	}
	randomize()
	for i := 0; i < k; i++ {
		resc[i].Center = rand.Intn(data.Len()/k) + data.Len()/k*i
	}
	changed := true
	count := 0
	for changed {
		classify(data, resc)
		changed = center(data, resc)
		count++
		if maxloop >= 0 {
			if count > maxloop {
				break
			}
		}
	}
	log.Printf("kmeans: looped %d times", count)
	return resc
}
