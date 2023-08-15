package src

import "errors"

type GachaService struct {
	api           ArknightsApi
	data          GachaData
	UpdateChannel chan string
	analyst       Analyst
}

func (s GachaService) updateUser(user User) (int, error) {
	count := 0
	apiUser, err := s.api.FindUser(user.Token)
	if err != nil {
		var responseDataStatusError ResponseDataStatusError[User]
		switch {
		case errors.As(err, &responseDataStatusError):
			switch responseDataStatusError.Status {
			case 3:
				return 0, errors.New(responseDataStatusError.Msg)
			}
			break
		default:
			return count, err
		}
	}
	if apiUser.NickName == "" {
		Logger.Errorf("nickName is empty, user: %+v", user)
		return 0, nil
	}
	if apiUser.NickName != user.NickName {
		s.data.UpdateName(user.Uid, apiUser.NickName)
		Logger.Infof("uid: %s change name, %s -> %s", user.Uid, user.NickName, apiUser.NickName)
	}

	page := 1
	needNextPage := true
	for needNextPage {
		paginationGacha, err := s.api.GetGacha(user.Token, user.ChannelMasterId, page)
		if err != nil {
			return count, err
		}
		gachas := paginationGacha.List
		if len(gachas) == 0 {
			break
		}
		for _, gacha := range gachas {
			gacha.Uid = user.Uid
			if s.data.HasGacha(gacha.Uid, gacha.Ts) {
				needNextPage = false
			} else {
				s.data.AddGacha(gacha)
				count += 1
			}
		}
		page += 1
	}
	return count, nil
}

func (s GachaService) task() {
	defer Logger.Error("task end")
	for {
		uid := <-s.UpdateChannel
		user := s.data.GetUser(uid)
		count, err := s.updateUser(user)
		if err != nil {
			Logger.WithError(err).Infof("update user %+v", user)
		}
		if count > 0 {
			Logger.Infof("user:%s, update: %d", user.NickName, count)
			s.analyst.Analyze(user.Uid)
		}
	}
}

func NewGachaService(data GachaData, analyst Analyst) *GachaService {
	service := &GachaService{
		api:           ArknightsApi{},
		data:          data,
		UpdateChannel: make(chan string),
		analyst:       analyst,
	}
	go service.task()
	return service
}
