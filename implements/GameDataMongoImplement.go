package implements

import (
	"GachaServerGin/src"
	"GachaServerGin/tools"
	"context"
	"errors"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type GameDataMongoImplement struct {
	client *qmgo.Client
	users  *qmgo.Collection
	gachas *qmgo.Collection

	getUser func(uid string) (src.User, error)
}

func NewMongoData() (GameDataMongoImplement, error) {
	client, err := qmgo.NewClient(context.Background(), &qmgo.Config{Uri: "mongodb://localhost:27017"})
	if err != nil {
		return GameDataMongoImplement{}, err
	}
	impl := GameDataMongoImplement{
		client: client,
		users:  client.Database("moles").Collection("doctors"),
		gachas: client.Database("moles").Collection("gachas"),
	}
	count, err := impl.users.Find(context.Background(), bson.M{}).Count()
	if err != nil {
		return GameDataMongoImplement{}, err
	}
	src.Logger.Infof("users count: %d", count)
	count, err = impl.gachas.Find(context.Background(), bson.M{}).Count()
	if err != nil {
		return GameDataMongoImplement{}, err
	}
	src.Logger.Infof("gachas count: %d", count)
	impl.getUser = func(uid string) (src.User, error) {
		var (
			user src.User
			errB error
		)
		errB = impl.users.Find(context.Background(), bson.M{"uid": uid}).One(&user)
		if err != nil {
			return src.User{}, errB
		}
		return user, nil
	}
	impl.getUser = tools.Cache11e(impl.getUser)
	if err != nil {
		return GameDataMongoImplement{}, err
	}
	return impl, nil
}

func (g GameDataMongoImplement) AddUser(user src.User) {
	one, err := g.users.InsertOne(context.Background(), user)
	if err != nil {
		return
	}
	src.Logger.Infof("user \"%s\" insert _id:%s", user.NickName, one.InsertedID)
}

func (g GameDataMongoImplement) GetUser(uid string) src.User {
	user, err := g.getUser(uid)
	if err != nil {
		src.Logger.Errorf("GetUser(%s) has error: %e", uid, err)
	}
	return user
}

func (g GameDataMongoImplement) GetUsers() []src.User {
	var users []src.User
	err := g.users.Find(context.Background(), bson.M{}).Select(bson.M{"token": 0}).All(&users)
	if err != nil {
		src.Logger.Error(err)
		return []src.User{}
	}
	return users
}

func (g GameDataMongoImplement) GetUsersWithToken() []src.User {
	var users []src.User
	err := g.users.Find(context.Background(), bson.M{}).All(&users)
	if err != nil {
		src.Logger.Error(err)
		return []src.User{}
	}
	return users
}

func (g GameDataMongoImplement) UpdateToken(uid string, token string) {
	err := g.users.UpdateOne(context.Background(), bson.M{"uid": uid}, bson.M{"$set": bson.M{"token": token}})
	if err != nil {
		src.Logger.Error(err)
		return
	}
}

func (g GameDataMongoImplement) UpdateName(uid string, name string) {
	err := g.users.UpdateOne(context.Background(), bson.M{"uid": uid}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		src.Logger.Error(err)
		return
	}
}

func (g GameDataMongoImplement) AddGacha(gacha src.Gacha) {
	if gacha.Uid == "" {
		src.Logger.Error(errors.New("gacha loss id"))
		return
	}
	_, err := g.gachas.InsertOne(context.Background(), gacha)
	if err != nil {
		src.Logger.Error(err)
		return
	}
}

func (g GameDataMongoImplement) GetGacha(uid string, ts int) src.Gacha {
	var gacha src.Gacha
	err := g.gachas.Find(context.Background(), bson.M{"uid": uid, "ts": ts}).One(&gacha)
	if err != nil {
		src.Logger.Error(err)
		return src.Gacha{}
	}
	return gacha
}

func (g GameDataMongoImplement) HasGacha(uid string, ts int) bool {
	count, err := g.gachas.Find(context.Background(), bson.M{"uid": uid, "ts": ts}).Count()
	if err != nil {
		src.Logger.Error(err)
		return false
	}
	return count > 0
}

func (g GameDataMongoImplement) GetGachasByPage(uid string, page int, pageSize int) src.PaginationData[src.Gacha] {
	var gachas []src.Gacha
	filter := bson.M{}
	if uid != "" {
		filter["uid"] = uid
	}
	//src.Logger.Infof("GetGachasByPage: filter %+v", filter)
	query := g.gachas.Find(context.Background(), filter).Sort("-ts")
	total, err := query.Count()
	if err != nil {
		src.Logger.Error(err)
		return src.PaginationData[src.Gacha]{}
	}
	query = query.Skip(int64(page * pageSize)).Limit(int64(pageSize))
	err = query.All(&gachas)
	if err != nil {
		src.Logger.Error(err)
		return src.PaginationData[src.Gacha]{}
	}
	return src.PaginationData[src.Gacha]{
		List: gachas,
		Pagination: src.Pagination{
			Current:  page,
			Total:    int(total),
			PageSize: pageSize,
		},
	}
}
