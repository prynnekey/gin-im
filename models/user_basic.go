package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type UserBasic struct {
	Identidy  string `json:"identidy"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Nickname  string `json:"nickname"`
	Gender    int    `json:"gender"`
	Email     string `json:"email"`
	Avator    string `json:"avator"`
	Create_at int64  `json:"create___at"`
	Update_at int64  `json:"update___at"`
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
