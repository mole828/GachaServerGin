package implements

import (
	"testing"
)

func TestNewMongoData(t *testing.T) {
	data, err := NewMongoData()
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("find", func(t *testing.T) {
		t.Logf("%+v", data.GetUser("69023059"))
	})
	t.Run("find all", func(t *testing.T) {
		t.Logf("%+v", data.GetUsers())
	})
}
