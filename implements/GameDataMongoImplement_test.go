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
		t.Logf("%+v", data.GetUser("63059"))
	})
	t.Run("find all", func(t *testing.T) {
		t.Logf("%+v", data.GetUsers())
	})
	t.Run("update", func(t *testing.T) {
		uid := "452653239"
		t.Logf("%+v", data.GetUser(uid))
		data.UpdateToken(uid, "")
		t.Logf("%+v", data.GetUser(uid))
	})
}
