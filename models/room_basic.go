package models

type RoomBasic struct {
	Identity     string `bson:"identity"`
	UserIdentity string `bson:"user_identity"`
	Number       string `bson:"number"`
	Name         string `bson:"name"`
	Info         string `bson:"info"`
	CreateAt     int64  `bson:"create_at"`
	UpdateAt     int64  `bson:"update_at"`
}

func (RoomBasic) CollectionName() string {
	return "room_basic"
}
