package zipper

import (
	"testing"
)

func TestContains(t *testing.T) {
	test := []string{
		"hello",
		"world",
	}
	exists := contains(test, "hello")
	if !exists {
		t.Errorf("contains was incorrect")
	}
	exists = contains(test, "invalid")
	if exists {
		t.Errorf("contains was incorrect")
	}
}