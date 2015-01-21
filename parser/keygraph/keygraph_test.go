package keygraph_test

import (
	"parser/keygraph"
	"testing"
)

func TestNewParser(t *testing.T) {
	pas := keygraph.NewParser()
	if pas == nil {
		t.Fatal("can not create parser")
	}
}
