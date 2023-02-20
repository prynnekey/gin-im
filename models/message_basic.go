package models

type MessageBasic struct {
	Identity     string `bson:"identity"`
	UserIdentity string `bson:"user_identity"`
	RootIdentity string `bson:"root_identity"`
	Data         string `bson:"data"`
	CreateAt     int64  `bson:"create_at"`
	UpdateAt     int64  `bson:"update_at"`
}

func (MessageBasic) CollectionName() string {
	return "message_basic"
}
