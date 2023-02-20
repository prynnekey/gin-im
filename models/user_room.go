package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRoom struct {
	UserIdentity    string `bson:"user_identity"`
	RoomIdentity    string `bson:"room_identity"`
	MessageIdentity string `bson:"message_identity"`
	CreateAt        int64  `bson:"create_at"`
	UpdateAt        int64  `bson:"update_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}

// 通过用户标识和房间标识获取用户房间
func GetUserRoomByUserIdentityAndRoomIdentity(userIdentity, RoomIdentity string) (*UserRoom, error) {
	userRoom := &UserRoom{}
	err := Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.D{{Key: "user_identity", Value: userIdentity}, {Key: "room_identity", Value: RoomIdentity}}).Decode(userRoom)

	return userRoom, err
}

// 通过房间号获取用户房间
func GetUserRoomByRoomIdentity(roomIdentity string) ([]*UserRoom, error) {
	var userRooms []*UserRoom
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{{Key: "room_identity", Value: roomIdentity}})
	if err != nil {
		return userRooms, err
	}
	if err = cursor.All(context.Background(), &userRooms); err != nil {
		return userRooms, err
	}
	return userRooms, nil
}
