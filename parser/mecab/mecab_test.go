package mecab_test

import (
	"parser/mecab"
	"testing"
)

func TestNewParser(t *testing.T) {
	pas := mecab.NewParser()
	if pas == nil {
		t.Fatal("can not create parser")
	}
}
