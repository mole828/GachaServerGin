package tools

import (
	"errors"
	"fmt"
	"testing"
)

func TestCache11e(t *testing.T) {
	c11e := Cache11e[int, string](func(i int) (string, error) {
		if i == 2 {
			return "", errors.New("i=2")
		}
		return fmt.Sprintf("i:%d", i), nil
	})
	for i := range []int{1, 2, 3} {
		t.Log(c11e(i))
	}
}

func TestCache12(t *testing.T) {
	c12 := Cache12(func(i int) (int, int) {
		return i, i + 1
	})
	for i := range []int{1, 2, 3} {
		t.Log(c12(i))
	}
}
