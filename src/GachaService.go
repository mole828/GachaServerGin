package src

import (
	"time"
)

//type GachaDatabase struct {
//	Users  *qmgo.Collection
//	Gachas *qmgo.Collection
//}
//
//func newDataBase(mongoUri string) (GachaDatabase, error) {
//	ctx := context.Background()
//	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: mongoUri})
//	if err != nil {
//		return GachaDatabase{}, err
//	}
//	data := GachaDatabase{
//		Users:  client.Database("moles").Collection("doctors"),
//		Gachas: client.Database("moles").Collection("gachas"),
//	}
//	return data, nil
//}

type GachaService struct {
	api  ArknightsApi
	data GachaData
}

func (s GachaService) updateUser(user User) error {
	apiUser, err := s.api.GetUser(user.Token)
	if err != nil {
		return err
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
			return err
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
			}
		}
		page += 1
	}
	return nil
}

func (s GachaService) fetchAll() {
	for _, user := range s.data.GetUsers() {
		err := s.updateUser(user)
		if err != nil {
			Logger.Error(err)
		}
		time.Sleep(time.Minute)
	}
}

func NewGachaService(data GachaData) *GachaService {

	return &GachaService{
		api:  ArknightsApi{},
		data: data,
	}
}
