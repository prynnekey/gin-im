package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type RoomBasic struct {
	Identity     string `bson:"identity"`
	UserIdentity string `bson:"user_identity"`
	Number       string `bson:"number"`
	Name         string `bson:"name"`
	Info         string `bson:"info"`
	CreateAt     int64  `bson:"create_at"`
	UpdateAt     int64  `bson:"update_at"`
}

type RoomBasicSimple struct {
	Name string `json:"name"`
	Info string `json:"info"`
}

func (RoomBasic) CollectionName() string {
	return "room_basic"
}

func InsertOneRoomBasic(roomBasic *RoomBasic) error {
	_, err := Mongo.Collection(RoomBasic{}.CollectionName()).
		InsertOne(context.Background(), roomBasic)

	return err
}

func DeleteRoomBasicByRoomIdentity(identity string) error {
	_, err := Mongo.Collection(RoomBasic{}.CollectionName()).
		DeleteOne(context.Background(), bson.D{{Key: "identity", Value: identity}})

	return err
}

func GetRoomBasicByRoomNumber(number string) (*RoomBasic, error) {
	rb := &RoomBasic{}
	err := Mongo.Collection(RoomBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{Key: "number", Value: number}}).Decode(rb)

	return rb, err
}

func GetRoomBasicByRoomIdentity(identities string) (*RoomBasic, error) {
	rb := &RoomBasic{}
	err := Mongo.Collection(RoomBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{Key: "identity", Value: identities}}).Decode(rb)

	return rb, err
}
