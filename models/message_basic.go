package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// 根据房间identity获取聊天记录
func GetMessageBasicByRootIdentity(rootIdentity string, pageIndex, pageSize int64) ([]*MessageBasic, error) {
	var messageBasics []*MessageBasic

	offset := (pageIndex - 1) * pageSize

	cursor, err := Mongo.Collection(MessageBasic{}.CollectionName()).
		Find(context.Background(), bson.M{"root_identity": rootIdentity},
			// 分页查询
			&options.FindOptions{
				Limit: &pageSize,               // Limit
				Skip:  &offset,                 // Offset
				Sort:  bson.M{"create_at": -1}, // Sort
			})

	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		mb := &MessageBasic{}
		cursor.Decode(mb)
		messageBasics = append(messageBasics, mb)
	}

	return messageBasics, err
}
