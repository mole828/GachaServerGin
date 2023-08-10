package src

import (
	"encoding/json"
	"testing"
)

var token = "j+GRJigVftyXQXmz/sWvgXEE"

func TestArknightsApi_GetUser(t *testing.T) {
	user, err := ArknightsApi{}.GetUser(token)
	if err != nil {
		t.Error(err)
	}
	t.Log(user)
}

func TestArknightsApi_GetGacha(t *testing.T) {
	t.Run("getGacha", func(t *testing.T) {
		body, err := getGachaBody(token, 1, 1)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(string(body))
		var data ResponseData[PaginationData[Gacha]]
		err = json.Unmarshal(body, &data)
		if err != nil {
			t.Error(err)
			return
		}
		t.Logf("%+v", data)
	})
	t.Run("GetGacha", func(t *testing.T) {
		gacha, err := ArknightsApi{}.GetGacha(token, 1, 11)
		if err != nil {
			t.Error(err)
		}
		t.Logf("gacha: %+v", gacha)
	})
}
