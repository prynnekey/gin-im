package models

import "context"

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

func InsertOneMessageBasic(messageBasic *MessageBasic) error {
	_, err := Mongo.Collection(MessageBasic{}.CollectionName()).InsertOne(context.Background(), messageBasic)
	return err
}
