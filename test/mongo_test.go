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
