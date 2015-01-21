/*
tfidf-kmeans Main priject file
Zhai Hongjie (c)2012
All rights reserved.
*/

package main

import (
	"flag"
	"fmt"
	"kmeans"
	"log"
	"parser"
	"tfidf"
)

type data struct {
	name  []string
	value []float32
}

func (datas data) Len() int {
	return len(datas.value)
}

func (datas data) Distance(i, j int) float64 {
	return float64(datas.value[i] - datas.value[j])
}

func main() {
	input := flag.String("i", "", "input file")
	sfilter := flag.String("f", "english", "filter")
	charset := flag.String("c", "", "charset")
	parser := flag.String("p", "treetagger", "set parser")
	k := flag.Int("k", 2, "number of class")
	flag.Parse()
	fmt.Println("***************************")
	fmt.Println("*　　      TF-IDF2        *")
	fmt.Println("*   zhai hongjie (c)2012  *")
	fmt.Println("*   All rights reserved   *")
	fmt.Println("***************************")
	if *input == "" {
		fmt.Println("usage: ./tfkmeans [options]")
		flag.PrintDefaults()
		return
	}
	pas := kparser.NewParser(*parser)
	if pas == nil {
		log.Fatalf("can not create parser for %s", *parser)
	}
	pas.Open(*input, *charset)
	defer pas.Close()
	filter := kparser.NewFilter()
	filter.Initialize(*sfilter)
	tfidf := ktfidf.CalculateTFIDF(pas, filter)
	datas := new(data)
	datas.name = make([]string, len(tfidf))
	datas.value = make([]float32, len(tfidf))
	count := 0
	for i, v := range tfidf {
		datas.name[count] = i
		datas.value[count] = v
		count++
	}
	classes := kmeans.Kmeans(datas, *k, 100)
	for i := 0; i < len(classes); i++ {
		fmt.Printf("class %d with center %s:\n", i, datas.name[classes[i].Center])
		s := ""
		for j := 0; j < len(classes[i].Classes); j++ {
			s += datas.name[classes[i].Classes[j]] + " "
		}
		fmt.Println(s)
	}
	log.Println("finished")
}
