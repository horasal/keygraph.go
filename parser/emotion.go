/*
KeyGraph emotion initializer
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

type emotion struct {
	positive map[string]int
	negitive map[string]int
	count    int
}

func (e *emotion) Emotion(token string) float32 {
	if _, ok := e.positive[strings.ToUpper(strings.TrimSpace(token))]; ok {
		return 1
	}
	if _, ok := e.negitive[strings.ToUpper(strings.TrimSpace(token))]; ok {
		return -1
	}
	return 0
}

func (e *emotion) Add(token string, emo float32) {
	if emo < 0 {
		e.negitive[strings.ToUpper(strings.TrimSpace(token))] = -1
	} else {
		e.positive[strings.ToUpper(strings.TrimSpace(token))] = 1
	}
}

func (e *emotion) Initialize(identifier string) {
	e.positive = make(map[string]int)
	e.negitive = make(map[string]int)
	if identifier == "" {
		return
	}
	f, err := os.Open(identifier + ".pe")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()
	buffer := bufio.NewReader(f)
	for token, err := buffer.ReadString('\n'); token != ""; e.count++ {
		if err != nil {
			if err != io.EOF {
				log.Fatal(err.Error())
			}
			break
		}
		if strings.TrimSpace(token) == "" {
			continue
		}
		e.positive[strings.ToUpper(strings.TrimSpace(token))] = 1
		token, err = buffer.ReadString('\n')
	}
	f, err = os.Open(identifier + ".ne")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer f.Close()
	buffer = bufio.NewReader(f)
	for token, err := buffer.ReadString('\n'); token != ""; e.count++ {
		if err != nil {
			if err != io.EOF {
				log.Fatal(err.Error())
			}
			break
		}
		if strings.TrimSpace(token) == "" {
			continue
		}
		e.negitive[strings.ToUpper(strings.TrimSpace(token))] = -1
		token, err = buffer.ReadString('\n')
	}
	log.Printf("eomtion %d readed.\n", e.count)
}

func NewEmotion() *emotion {
	return new(emotion)
}
