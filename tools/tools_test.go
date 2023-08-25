package tools

import (
	"testing"
	"time"
)

func TestLo(t *testing.T) {
	begin := time.Unix(0, 0)
	now := time.Now()
	t.Log(now.Sub(begin) > 0)
	t.Log(begin.Sub(now) > 0)
}

type person struct {
	Name string
}

func TestString(t *testing.T) {
	var s string
	t.Logf("\"%s\"", s)

	t.Logf("\"%s\"", person{}.Name)
}
