package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type UserBasic struct {
	Identity string `bson:"identity"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Nickname string `bson:"nickname"`
	Gender   int    `bson:"gender"`
	Email    string `bson:"email"`
	Avatar   string `bson:"avatar"`
	CreateAt int64  `bson:"create_at"`
	UpdateAt int64  `bson:"update_at"`
}

type UserInfo struct {
	Identity string `bson:"identity"`
	Username string `bson:"username"`
	Nickname string `bson:"nickname"`
	Gender   int    `bson:"gender"`
	Avatar   string `bson:"avatar"`
	IsFriend bool   `json:"is_friend"`
}

func (UserBasic) CollectionName() string {
	return "user_basic"
}

// 通过用户名和密码获取用户信息
func GetUserBasicByUsernameAndPassword(username, password string) (*UserBasic, error) {
	ub := &UserBasic{}
	// 查询数据
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{Key: "username", Value: username}, {Key: "password", Value: password}}).
		Decode(ub)

	return ub, err
}

// 通过用户名获取用户信息
func GetUserBasicByUsername(username string) (count int64, err error) {
	// 查询数据
	return Mongo.Collection(UserBasic{}.CollectionName()).
		CountDocuments(context.Background(), bson.D{{Key: "username", Value: username}})
}

// 将输入插入到MongoDB
func InsertOneUserBasic(ub *UserBasic) error {
	// 插入数据
	_, err := Mongo.Collection(UserBasic{}.CollectionName()).
		InsertOne(context.Background(), ub)

	return err
}

// 根据Identity查询用户数据
func GetUserBasicByIdentity(identity string) (*UserBasic, error) {
	ub := &UserBasic{}
	// 查询数据
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{Key: "identity", Value: identity}}).
		Decode(ub)

	return ub, err
}

func GetUserInfoByUsername(username string) (*UserInfo, error) {
	ui := &UserInfo{}
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), bson.D{{Key: "username", Value: username}}).
		Decode(ui)

	return ui, err
}
