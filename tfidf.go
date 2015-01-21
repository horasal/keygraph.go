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
	"log"
	"os"
	"parser"
	"tfidf"
)

func main() {
	input := flag.String("i", "", "input file")
	output := flag.String("o", "", "output file")
	sfilter := flag.String("f", "english", "filter")
	charset := flag.String("c", "", "charset")
	parser := flag.String("p", "treetagger", "set parser")
	r := flag.Bool("r", false, "RIDF")
	flag.Parse()
	fmt.Println("***************************")
	fmt.Println("*　　      TF-IDF         *")
	fmt.Println("*   zhai hongjie (c)2012  *")
	fmt.Println("*   All rights reserved   *")
	fmt.Println("***************************")
	if *input == "" {
		fmt.Println("usage: ./tfidf [options]")
		flag.PrintDefaults()
		return
	}
	if *output == "" {
		*output = *input + ".tfidf.txt"
	}
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
	var tfidf map[string]float32
	if *r {
		tfidf = ktfidf.CalculateRIDF(pas, filter)
	} else {
		tfidf = ktfidf.CalculateTFIDF(pas, filter)
	}
	ktfidf.Save(tfidf, buf)
	log.Println("finished")
}
