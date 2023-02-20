package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/prynnekey/gin-im/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestFindOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database("im")
	ub := models.UserBasic{}
	err2 := db.Collection("user_basic").FindOne(context.Background(), bson.D{}).Decode(&ub)
	if err2 != nil {
		t.Fatal(err2)
	}

	fmt.Printf("ub: %v\n", ub)
}

// 通过房间号获取用户列表
func TestFind(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}

	db := client.Database("im")
	cursor, err := db.Collection("user_room").Find(context.Background(), bson.D{})
	if err != nil {
		t.Fatal(err)
	}

	// 存储房间结果
	urs := make([]*models.UserRoom, 0)
	for cursor.Next(context.Background()) {
		ur := &models.UserRoom{}
		err := cursor.Decode(ur)
		if err != nil {
			t.Fatal(err)
		}
		urs = append(urs, ur)
	}

	// 遍历打印
	for _, ur := range urs {
		fmt.Printf("user_room: %v\n", ur)
	}
}
