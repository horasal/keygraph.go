/*
User classify Main priject file
Zhai Hongjie (c)2012
All rights reserved.
*/

package main

import (
	"flag"
	"fmt"
	"kmeans"
	"log"
	"os"
	"parser"
	"strconv"
)

type user struct {
	username string
	keyword  []string
	value    []float64
}

type data struct {
	users []user
}

func (datas data) Len() int {
	return len(datas.users)
}

func (datas data) Distance(i, j int) float64 {
	count := 0.0
	for i1 := 0; i1 < len(datas.users[i].keyword); i1++ {
		for j1 := 0; j1 < len(datas.users[j].keyword); j1++ {
			if datas.users[i].keyword[i1] == datas.users[j].keyword[j1] {
				count += datas.users[i].value[i1] * datas.users[j].value[j1]
				break
			}
		}
	}
	if i < j {
		count = -count
	}
	return 1.0 / float64(count)
}

func fillusers(maxword int, parser kparser.Token, base string, fi []os.FileInfo, datas *data, charset string) {
	datas.users = make([]user, len(fi))
	for i, v := range fi {
		parser.Open(base+"/"+v.Name(), charset)
		parser.Reset()
		count := 0
		total := 0
		datas.users[i].keyword = make([]string, maxword)
		datas.users[i].value = make([]float64, maxword)
		datas.users[i].username = v.Name()
		for word, err := parser.Next(); err == nil; count++ {
			if count >= maxword {
				break
			}
			datas.users[i].keyword[count] = word
			a, err := strconv.Atoi(parser.Propery())
			datas.users[i].value[count] = float64(a)
			if err != nil {
				datas.users[i].value[count] = -1
				total++
				log.Println(err.Error())
			}
			if total < int(datas.users[i].value[count]) {
				total = int(datas.users[i].value[count])
			}
			word, err = parser.Next()
		}
		for j, _ := range datas.users[i].value {
			datas.users[i].value[j] /= float64(total)
		}
		if count < maxword {
			log.Printf("not enough keyword at %s : %d vs %d", v.Name(), maxword, count)
		}
		parser.Close()
	}
}

func main() {
	input := flag.String("i", "", "input file")
	charset := flag.String("c", "", "charset")
	parser := flag.String("p", "keygraph", "set parser")
	k := flag.Int("k", 2, "number of class")
	maxword := flag.Int("m", 50, "max word to compare")
	flag.Parse()
	if *input == "" {
		fmt.Println("usage: ./userclassify [options]")
		flag.PrintDefaults()
		return
	}
	pas := kparser.NewParser(*parser)
	if pas == nil {
		log.Fatalf("can not create parser for %s", *parser)
	}
	list, err := os.Open(*input)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer list.Close()
	if a, err := list.Stat(); err != nil || !a.IsDir() {
		log.Fatalf("link %s is not dir", *input)
	}
	fi, err := list.Readdir(-1)
	if err != nil {
		log.Fatal(err.Error())
	}
	datas := new(data)
	fillusers(*maxword, pas, *input, fi, datas, *charset)
	classes := kmeans.Kmeans(datas, *k, 100)
	for i := 0; i < len(classes); i++ {
		fmt.Printf("class %d with center %s:\n", i, datas.users[classes[i].Center].username)
		s := ""
		for j := 0; j < len(classes[i].Classes); j++ {
			s += datas.users[classes[i].Classes[j]].username + " "
		}
		fmt.Println(s)
		fmt.Println()
	}
	log.Println("finished")
}
