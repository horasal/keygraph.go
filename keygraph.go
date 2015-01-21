/*
KeyGraph Main priject file
Zhai Hongjie (c)2012
All rights reserved.
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"graph"
	"log"
	"os"
	"parser"
)

func printtable(base kgraph.BGraph) {
	fmt.Printf("%8s\t", "TBL")
	for i := 0; i < len(base.Linestable); i++ {
		for v, j := range base.Nametable {
			if i == j {
				fmt.Printf("%8s\t", v)
			}
		}
	}
	fmt.Println("")
	for i := 0; i < len(base.Linestable); i++ {
		for v, j := range base.Nametable {
			if i == j {
				fmt.Printf("%8s\t", v)
			}
		}
		for j := 0; j < len(base.Linestable[i]); j++ {
			fmt.Printf("%8d\t", base.Linestable[i][j])
		}
		fmt.Println("")
	}
}

func main() {
	input := flag.String("i", "", "input file")
	output := flag.String("o", "", "output file")
	emotionfile := flag.String("e", "", "emotion")
	sfilter := flag.String("f", "english", "filter")
	charset := flag.String("c", "", "charset")
	highfreqNum := flag.Int("n", 15, "high frequency words number")
	threshold := flag.Int("t", 100, "threshold")
	showtable := flag.Bool("s", false, "show table content")
	auto := flag.Bool("a", false, "auto adjust parameter")
	parser := flag.String("p", "treetagger", "set parser")
	propery := flag.String("x", "", "set filter for propery")
	flag.Parse()
	fmt.Println("***************************")
	fmt.Println("*　　     KeyGraph        *")
	fmt.Println("*   zhai hongjie (c)2012  *")
	fmt.Println("*   All rights reserved   *")
	fmt.Println("***************************")
	if *input == "" {
		fmt.Println("usage: ./keygraph [options]")
		flag.PrintDefaults()
		return
	}
	if *output == "" {
		*output = *input + ".key.txt"
	}
	if *auto {
		*highfreqNum = 0
		*threshold = -1
	}
	if *emotionfile != "" {
		emotion := kparser.NewEmotion()
		emotion.Initialize(*emotionfile)
	}
	log.Printf("\nPropery:%s\nParser:%s\nFilter:%s\n",*propery,*parser,*sfilter)
	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	pas := kparser.NewParser(*parser)
	if pas == nil {
		log.Fatalf("can not create parser for %s", *parser)
	}
	pas.Open(*input, *charset)
	defer pas.Close()
	filter := kparser.NewFilter()
	filter.Initialize(*sfilter)

	highfreq := kgraph.GenHighFreq(pas, filter, *highfreqNum)
	if *showtable {
		for v, i := range highfreq {
			log.Printf("word:%s times:%d \n", v, i)
		}
	}
	pas.Reset()
	base := kgraph.BuildGraph(highfreq, pas)
	if *showtable {
		printtable(*base)
	}
	kgraph.InitializeBaseEmotion(base, highfreq, nil, pas)
	kgraph.SelectEdge(base, nil, 3, *threshold)
	if *showtable {
		printtable(*base)
	}
	pas.Reset()
	result := kgraph.SumCo(pas, filter, *base)
	pas.Reset()
	log.Printf("save file to %s", *output)
	kgraph.GetPeopry(pas, filter, &result)
	kgraph.Save(result, *propery, buf)
	log.Println("finished")
	//kgraph.DrawGraph(*base, os.Stdout, 500, 500, 15)
	/*	for i, v := range base.Emotiontable {
			fmt.Println(i, "\t", v)
		}
				for i := 0; i < len(base.Linestable); i++ {
				for j := 0; j < len(base.Linestable[i]); j++ {
					base1 := *base
					fmt.Print(kgraph.HasPath(base1, i, j), "\t")
				}
				fmt.Println("")
			}
			printtable(*base)*/

}
