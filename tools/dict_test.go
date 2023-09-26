package tools

import (
	"encoding/json"
	"github.com/samber/lo"
	"testing"
)

func TestCounter(t *testing.T) {
	counter := NewCounter[string]()
	counter.Inc("a", 5)
	marshal, err := json.Marshal(counter.Data())
	if err != nil {
		t.Error(err)
	}
	t.Log(string(marshal))
}

func TestLoCount(t *testing.T) {
	lo.CountValuesBy([]int{1, 2, 3}, func(item int) bool {
		return item%2 == 0
	})
}
