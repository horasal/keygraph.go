/*
KeyGraph word filter
Zhai Hongjie (c)2012
All rights reserved.
*/

package kparser

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

type listfilter struct {
	list  map[string]int
	count int
}

func (filter listfilter) Has(token string) bool {
	_, ok := filter.list[token]
	return ok
}

func (filter *listfilter) Add(token string) {
	if !filter.Has(token) {
		filter.list[token] = 1
	}
}

func (filter *listfilter) Initialize(language string) {
	filter.list = make(map[string]int)
	if language == "" {
		return
	}
	f, err := os.Open(language + ".stop")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()
	buffer := bufio.NewReader(f)
	i := 0
	for token, err := buffer.ReadString('\n'); token != ""; i++ {
		if err != nil {
			if err != io.EOF {
				log.Fatal(err.Error())
			}
			break
		}
		tokenlist := strings.Split(token, "\t")
		filter.list[strings.ToUpper(strings.TrimSpace(tokenlist[0]))] = 1
		token, err = buffer.ReadString('\n')
	}
	log.Printf("stop word %d readed.\n", i)
}

func NewFilter() *listfilter {
	return new(listfilter)
}
