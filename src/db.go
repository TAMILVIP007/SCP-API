package src

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	bannedinfo *mongo.Collection
	tokensinfo *mongo.Collection
)

func init() {
	log.Println("Setting up database...")
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(Envars.DbUrl))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	bannedinfo = client.Database("scpfoundation").Collection("bannedinfo")
	tokensinfo = client.Database("scpfoundation").Collection("tokensinfo")
}

func AddNewBan(u *BannedInfo) error {
	filter := bson.M{"user_id": u.UserId}
	update := bson.M{"$set": u}
	_, err := bannedinfo.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	return err
}

func CheckBan(u *BannedInfo) bool {
	err := bannedinfo.FindOne(context.Background(), u).Decode(&u)
	return err == nil
}

func GetBanReason(userid string) string {
	var u BannedInfo
	err := bannedinfo.FindOne(context.Background(), bson.M{"user_id": userid}).Decode(&u)
	if err != nil {
		return ""
	}
	return u.Reason
}

func AddNewToken(u *TokensInfo) error {
	filter := bson.M{"user_id": u.UserId}
	update := bson.M{"$set": u}
	_, err := tokensinfo.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	return err
}

func FetchTokenRole(token string) *TokensInfo {
	var u TokensInfo
	err := tokensinfo.FindOne(context.Background(), bson.M{"token": token}).Decode(&u)
	if err != nil {
		return nil
	}
	return &u
}

func RemoveBanById(userid string) error {
	_, err := bannedinfo.DeleteOne(context.Background(), bson.M{"user_id": userid})
	return err
}
