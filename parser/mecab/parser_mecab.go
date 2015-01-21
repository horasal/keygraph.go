/*
KeyGraph word parser for MeCab
Zhai Hongjie (c)2012
All rights reserved.
*/

package mecab

import (
	"bufio"
	"code.google.com/p/mahonia"
	"errors"
	"log"
	"os"
	"strings"
)

type Parser struct {
	count          int
	offset         int
	token          []string
	extraproperies []string
}

func (pars *Parser) Next() (string, error) {
	if pars.offset >= pars.count {
		return "", errors.New("EOF")
	}
	pars.offset++
	return pars.token[pars.offset-1], nil
}

func (pars *Parser) Propery() string {
	if pars.offset >= pars.count {
		return ""
	}
	return pars.extraproperies[pars.offset-1]
}

func (pars *Parser) Open(filename string, charset string) {
	var decoder mahonia.Decoder
	pars.count = 0
	pars.offset = 0
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer file.Close()
	buffer := bufio.NewReader(file)
	if charset != "utf-8" || charset != "" {
		decoder = mahonia.NewDecoder(charset)
	} else {
		decoder = nil
	}
	if charset == "" {
		log.Printf("open file %s with charset %s \n", filename, "default(UTF-8)")
	} else {
		log.Printf("open file %s with charset %s \n", filename, charset)
	}
	pars.token = make([]string, 0)
	pars.extraproperies = make([]string, 0)
	for line, err := buffer.ReadString('\n'); err == nil; line, err = buffer.ReadString('\n') {
		if decoder != nil {
			line = decoder.ConvertString(line)
		}
		strarray := strings.Split(line, "\t")
		if s := strings.ToUpper(strings.TrimSpace(strarray[0])); s == "" {
			pars.token = append(pars.token, " ")
		} else {
			pars.token = append(pars.token, s)
		}
		if len(strarray) < 2 {
			pars.extraproperies = append(pars.extraproperies, "")
			if !pars.IsEOS(pars.token[pars.count]) {
				log.Printf("missing word at %d-%s,%s \n", pars.count, pars.token[pars.count], line)
			}
			pars.count++
			continue
		} else {
			propery := strings.Split(strarray[1], ",")
			pars.extraproperies = append(pars.extraproperies, strings.ToUpper(strings.TrimSpace(propery[0])))
		}
		pars.count++
		if pars.count != len(pars.token) {
			log.Fatalf("error token count %d vs %d", pars.count, len(pars.token))
		}
	}
}

func (pars *Parser) Close() {
	pars.count = 0
}

func (pars *Parser) Reset() {
	pars.offset = 0
}

func (pars *Parser) IsEOS(token string) bool {
	if pars.offset >= pars.count {
		return true
	}
	return token == "EOS"
}

func NewParser() *Parser {
	return new(Parser)
}
