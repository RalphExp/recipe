package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipe struct {
	//swagger:ignore
	ID           string    `json:"id" bson:"_id"`
	Name         string    `json:"name" bson:"name"`
	Tags         []string  `json:"tags" bson:"tags"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients"`
	Instructions []string  `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time `json:"publishedAt" bson:"publishedAt"`
}

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

	/* insert new authorization data */
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	if err = collection.Drop(ctx); err != nil {
		log.Fatal(err)
	}

	h := sha256.New()
	for username, password := range users {
		h.Reset()
		io.Copy(h, strings.NewReader(password))
		collection.InsertOne(ctx, bson.M{
			"username": username,
			"password": string(hex.EncodeToString(h.Sum(nil))),
		})
	}

	/* insert recipes */
	recipes := make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	err = json.Unmarshal([]byte(file), &recipes)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	collection.Drop(ctx)

	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}
