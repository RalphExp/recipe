package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	users := map[string]string{
		"admin":    "fCRmh4Q2J7Rseqkz",
		"packt":    "RE4zfHB35VPtTkbT",
		"ralphexp": "L3nSFRcZzNQ67bcc",
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	h := sha256.New()

	for username, password := range users {
		h.Reset()
		io.Copy(h, strings.NewReader(password))
		collection.InsertOne(ctx, bson.M{
			"username": username,
			"password": string(hex.EncodeToString(h.Sum(nil))),
		})
	}
}
