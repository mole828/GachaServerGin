package src

import (
	"encoding/json"
	"errors"
	"testing"
)

var token = "j+GRJigVftyXQXmz/sWvgXEE"

func TestArknightsApi_GetUser(t *testing.T) {
	token = "48I0GjeojP2M01tjdnMcQafS"
	user, err := ArknightsApi{}.FindUser(token)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%+v", user)
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

func TestError(t *testing.T) {
	e := ResponseDataStatusError[any]{
		ResponseData[any]{
			Msg:    "err msg",
			Data:   "",
			Status: 3,
		},
	}
	var responseDataStatusError ResponseDataStatusError[any]
	switch {
	case errors.As(e, &responseDataStatusError):
		t.Log("is status error")
	}
	t.Log(errors.As(e, &ResponseDataStatusError[any]{}))
}
