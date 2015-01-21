package treetagger_test

import (
	"parser/treetagger"
	"testing"
)

func TestNewParser(t *testing.T) {
	pas := treetagger.NewParser()
	if pas == nil {
		t.Fatal("can not create parser")
	}
}
