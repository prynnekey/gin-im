package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRoom struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	RoomType     int    `bson:"room_type"` // 房间类型 【1-私聊,2-群聊】
	CreateAt     int64  `bson:"create_at"`
	UpdateAt     int64  `bson:"update_at"`
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

// 判断用户是否是好友
func JudgeUserIsFriend(userIdentity, friendUserIdentity string) (bool, error) {
	// 如果他们属于同一个房间且房间类型为私聊，则为好友
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{{Key: "user_identity", Value: userIdentity}, {Key: "room_type", Value: 1}})
	if err != nil {
		return false, err
	}

	// 获取用户房间集合
	roomIdentities := make([]string, 0)
	for cursor.Next(context.Background()) {
		userRoom := &UserRoom{}
		if err = cursor.Decode(userRoom); err != nil {
			return false, err
		}
		roomIdentities = append(roomIdentities, userRoom.RoomIdentity)
	}

	count, err := Mongo.Collection(UserRoom{}.CollectionName()).
		CountDocuments(context.Background(),
			bson.D{
				{Key: "user_identity", Value: friendUserIdentity},
				{Key: "room_identity", Value: bson.D{{Key: "$in", Value: roomIdentities}}}})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func InsertOneUserRoom(userRoom *UserRoom) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).
		InsertOne(context.Background(), userRoom)

	return err
}

// 根据自己id和好友id获取房间id
func GetRoomIdentityByUserIdentities(userIdentity, friendUserIdentity string) (string, error) {
	// 如果他们属于同一个房间且房间类型为私聊，则为好友
	cursor, err := Mongo.Collection(UserRoom{}.CollectionName()).
		Find(context.Background(), bson.D{{Key: "user_identity", Value: userIdentity}, {Key: "room_type", Value: 1}})
	if err != nil {
		return "", err
	}

	// 获取用户房间集合
	roomIdentities := make([]string, 0)
	for cursor.Next(context.Background()) {
		userRoom := &UserRoom{}
		if err = cursor.Decode(userRoom); err != nil {
			return "", err
		}
		roomIdentities = append(roomIdentities, userRoom.RoomIdentity)
	}

	ur := &UserRoom{}
	err = Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), bson.D{
			{Key: "user_identity", Value: friendUserIdentity},
			{Key: "room_type", Value: 1},
			{Key: "room_identity", Value: bson.D{{Key: "$in", Value: roomIdentities}}},
		}).Decode(ur)

	if err != nil {
		return "", err
	}

	return ur.RoomIdentity, nil
}

// 根据房间id删除房间
func DeleteUserRoomByRoomIdentity(roomIdentity string) error {
	_, err := Mongo.Collection(UserRoom{}.CollectionName()).
		DeleteMany(context.Background(), bson.D{
			{Key: "room_identity", Value: roomIdentity},
			{Key: "room_type", Value: 1},
		})

	return err
}
