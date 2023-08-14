package src

type GachaData interface {
	AddUser(user User)
	GetUser(uid string) User
	GetUsers() []User
	GetUsersWithToken() []User
	UpdateToken(uid string, token string)
	UpdateName(uid string, name string)

	AddGacha(gacha Gacha)
	GetGacha(uid string, ts int) Gacha
	HasGacha(uid string, ts int) bool
	GetGachasByPage(uid string, page int, pageSize int) PaginationData[Gacha]
}
