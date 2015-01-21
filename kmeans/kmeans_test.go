package kmeans_test

import (
	"kmeans"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const testSize = 1000
const classNum = 10
const maxloop = 100

type data struct {
	datas []int
}

func (datas data) Len() int {
	return testSize
}

func (datas data) Distance(i, j int) float64 {
	return float64(datas.datas[i] - datas.datas[j])
}

func randomize() {
	rand.Seed(time.Second.Nanoseconds())
}

func filldata(datas *data) {
	randomize()
	j := 0
	for j = 0; j < classNum; j++ {
		for i := 0; i < len(datas.datas)/classNum; i++ {
			datas.datas[i+j*len(datas.datas)/classNum] = rand.Intn(testSize/classNum/2) + testSize/classNum/4
		}
	}
	for i := j * len(datas.datas) / classNum; i < len(datas.datas); i++ {
		datas.datas[i] = rand.Intn(testSize)
	}
}

func TestKmeans(t *testing.T) {
	testdata := new(data)
	testdata.datas = make([]int, testSize)
	filldata(testdata)
	t.Logf("source: %x", testdata.datas)
	t.Logf("classify into %d classes", classNum)
	result := kmeans.Kmeans(testdata, classNum, maxloop)
	for i := 0; i < len(result); i++ {
		t.Logf("class %d with center %d:", i, testdata.datas[result[i].Center])
		s := ""
		for j := 0; j < len(result[i].Classes); j++ {
			s += strconv.Itoa(testdata.datas[result[i].Classes[j]]) + " "
		}
		t.Log(s)
	}
}
