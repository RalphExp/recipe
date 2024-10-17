package models

import (
	"time"
)

// swagger:parameters recipes newRecipe
type Recipe struct {
	//swagger:ignore
	ID          string       `json:"id" bson:"_id"`
	Name        string       `json:"name" bson:"name"`
	Tags        []string     `json:"tags" bson:"tags"`
	Ingredients []Ingredient `json:"ingredients" bson:"ingredients"`
	Steps       []string     `json:"steps" bson:"steps"`
	PublishedAt time.Time    `json:"publishedAt" bson:"publishedAt"`
	ImageURL    string       `json:"imageURL" bson:"imageURL"`
}

type Ingredient struct {
	Quantity string `json:"quantity" bson:"quantity"`
	Name     string `json:"name" bson:"name"`
	Type     string `json:"type" bson:"type"`
}
