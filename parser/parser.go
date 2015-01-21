/*
KeyGraph word parser for TreeTagger
Zhai Hongjie (c)2012
All rights reserved.
*/

package kparser

import (
	"parser/keygraph"
	"parser/mecab"
	"parser/temp"
	"parser/treetagger"
	"strings"
)

type Token interface {
	Next() (string, error)
	Open(filename string, charset string)
	IsEOS(token string) bool
	Close()
	Reset()
	Propery() string
}

func NewParser(identifier string) Token {
	switch strings.ToUpper(strings.TrimSpace(identifier)) {
	case "MECAB":
		{
			return mecab.NewParser()
		}
	case "TREETAGGER":
		{
			return treetagger.NewParser()
		}
	case "KEYGRAPH":
		{
			return keygraph.NewParser()
		}
	case "TEST":
		{
			return temp.NewParser()
		}
	}
	return nil
}
