package kparser_test

import (
	"parser"
	"testing"
)

func TestNewParser(t *testing.T) {
	pas := kparser.NewParser("treetagger")
	t.Log("testing treetagger parser..")
	if pas == nil {
		t.Fatal("can not create parser")
	}
	pas = kparser.NewParser("mecab")
	t.Log("testing mecab parser..")
	if pas == nil {
		t.Fatal("can not create parser")
	}
	pas = kparser.NewParser("keygraph")
	t.Log("testing keygraph parser..")
	if pas == nil {
		t.Fatal("can not create parser")
	}
	emotion := kparser.NewEmotion()
	t.Log("testing emotion parser..")
	if emotion == nil {
		t.Fatal("can not create parser")
	}
	filter := kparser.NewFilter()
	t.Log("testing filter..")
	if filter == nil {
		t.Fatal("can not create filter")
	}
}
