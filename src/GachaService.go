package src

type GachaService struct {
	api           ArknightsApi
	data          GachaData
	UpdateChannel chan string
}

func (s GachaService) updateUser(token string) (int, error) {
	apiUser, err := s.api.GetUser(token)
	apiUser.Token = token
	count := 0
	user := s.data.GetUser(apiUser.Uid)
	if err != nil {
		return count, err
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
			if s.data.HasGacha(gacha.Uid, gacha.Ts) {
				needNextPage = false
			} else {
				gacha.Uid = user.Uid
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
		token := <-s.UpdateChannel
		count, err := s.updateUser(token)
		if err != nil {
			Logger.Error(err)
		}
		if count > 0 {
			Logger.Infof("token:%s, update: %d", token, count)
		}
	}
}

func NewGachaService(data GachaData) *GachaService {
	service := &GachaService{
		api:           ArknightsApi{},
		data:          data,
		UpdateChannel: make(chan string),
	}
	go service.task()
	return service
}
