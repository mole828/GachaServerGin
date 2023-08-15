package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type ArknightsApi struct{}

type ResponseData[T any] struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   T      `json:"data"`
}

type ResponseDataStatusError[T any] struct {
	ResponseData[T]
}

func (r ResponseDataStatusError[T]) Error() string {
	return r.Msg
}

type User struct {
	Uid             string `json:"uid"`
	NickName        string `json:"nickName"`
	Token           string `json:"token"`
	Guest           int    `json:"guest"`
	ChannelMasterId int    `json:"channelMasterId"`
}

func closeResponseBody(body io.ReadCloser) {
	err := body.Close()
	if err != nil {
		Logger.Error(err)
	}
}

func (ArknightsApi) GetUser(token string, channelMasterId int) (User, error) {
	var postData map[string]interface{}
	if channelMasterId == 1 {
		postData = map[string]interface{}{
			"appId":           1,
			"channelMasterId": 1,
			"channelToken": map[string]string{
				"token": token,
			},
		}
	} else if channelMasterId == 2 {
		postData = map[string]interface{}{
			"token": token,
		}
	}

	postJSON, _ := json.Marshal(postData)
	postResponse, err := http.Post("https://as.hypergryph.com/u8/user/info/v1/basic", "application/json", bytes.NewBuffer(postJSON))
	if err != nil {
		return User{}, err
	}
	defer closeResponseBody(postResponse.Body)
	postResponseBody, err := io.ReadAll(postResponse.Body)
	if err != nil {
		return User{}, err
	}
	var responseData ResponseData[User]
	err = json.Unmarshal(postResponseBody, &responseData)
	if err != nil {
		Logger.Error(string(postResponseBody))
		return User{}, err
	}
	if responseData.Status != 0 {
		return User{}, ResponseDataStatusError[User]{
			ResponseData: responseData,
		}
	}
	data := responseData.Data
	if data.NickName == "" || data.Uid == "" {
		return data, errors.New(fmt.Sprintf("empty user data: %s", postResponseBody))
	}
	data.Token = token
	return data, nil
}

func (r ArknightsApi) FindUser(token string) (User, error) {
	user1, err1 := r.GetUser(token, 1)
	user2, err2 := r.GetUser(token, 2)
	if err1 != nil && err2 != nil {
		if errors.As(err1, &ResponseDataStatusError[User]{}) &&
			errors.As(err2, &ResponseDataStatusError[User]{}) {
			return User{}, err2
		}
		return User{}, errors.Join(err1, err2)
	}
	var user User
	if user1.Uid == "" {
		user = user2
	} else {
		user = user1
	}
	return user, nil
}

type Pagination struct {
	Current  int `json:"current"`
	Total    int `json:"total"`
	PageSize int `json:"pageSize"`
}

type PaginationData[T any] struct {
	List       []T `json:"list"`
	Pagination `json:"pagination"`
}

type Char struct {
	IsNew  bool   `json:"isNew"`
	Name   string `json:"name"`
	Rarity int    `json:"rarity"`
}

type Gacha struct {
	Chars    []Char `json:"chars"`
	Pool     string `json:"pool"`
	Ts       int    `json:"ts"`
	Uid      string `json:"uid"`      // add flied
	NickName string `json:"nickName"` // add flied
}

func getGachaBody(token string, channelId int, page int) ([]byte, error) {
	u, err := url.Parse("https://ak.hypergryph.com/user/api/inquiry/gacha")
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("token", token)
	q.Add("channelId", strconv.Itoa(channelId))

	u.RawQuery = q.Encode()

	// 执行 GET 请求
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Error(err)
		}
	}(resp.Body)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func (ArknightsApi) GetGacha(token string, channelId int, page int) (PaginationData[Gacha], error) {
	body, err := getGachaBody(token, channelId, page)
	var responseData ResponseData[PaginationData[Gacha]]
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return PaginationData[Gacha]{}, err
	}
	return responseData.Data, nil
}
