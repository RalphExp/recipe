package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/ralphexp/recipes-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

func (handler *RecipesHandler) getRecipes(c *gin.Context) ([]models.Recipe, error) {
	val, err := handler.redisClient.Get("recipes").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			return nil, err
		}
		defer cur.Close(handler.ctx)

		recipes := make([]models.Recipe, 0)
		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}

		data, _ := json.Marshal(recipes)
		handler.redisClient.Set("recipes", string(data), 0)
		return recipes, nil
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return nil, err
	} else {
		log.Printf("Request to Redis")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		return recipes, nil
	}
}

// swagger:operation GET /index recipes listRecipes
// Returns the web page of recipes
// ---
// produces:
// - application/html
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *RecipesHandler) IndexHandler(c *gin.Context) {
	recipes, err := handler.getRecipes(c)
	if recipes != nil {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"recipes": recipes})
	} else {
		c.HTML(http.StatusInternalServerError, "404.html", gin.H{"error": err.Error()})
	}
}

// swagger:operation GET /api/v1/recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	recipes, err := handler.getRecipes(c)
	if recipes != nil {
		c.JSON(http.StatusOK, recipes)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// swagger:operation POST /api/v1/recipes recipes newRecipe
// Create a new recipe
// ---
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'400':
//	    description: Invalid input
func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if recipe.ID == "" {
		recipe.ID = primitive.NewObjectID().Hex()
	}
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}

	log.Println("Remove data from Redis")
	handler.redisClient.Del("recipes")
	c.JSON(http.StatusOK, recipe)
}

// swagger:operation PUT /api/v1/recipes/{id} recipes updateRecipe
// Update an existing recipe
// ---
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// produces:
// - application/json
// responses:
//
//	'200':
//	    description: Successful operation
//	'400':
//	    description: Invalid input
//	'404':
//	    description: Invalid recipe ID
func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := handler.collection.UpdateOne(handler.ctx, bson.M{
		"_id": id,
	}, bson.D{{"$set", bson.D{
		{"_id", id},
		{"name", recipe.Name},
		{"steps", recipe.Steps},
		{"ingredients", recipe.Ingredients},
		{"imageURL", recipe.ImageURL},
	}}})

	fmt.Printf("%v\n", res)

	if res.MatchedCount == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Recipe not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// swagger:operation DELETE /api/v1/recipes/{id} recipes deleteRecipe
// Delete an existing recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: ID of the recipe
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
//	'404':
//	    description: Invalid recipe ID
func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	res, err := handler.collection.DeleteOne(handler.ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if res.DeletedCount > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Recipe not found"})
	}
}

// swagger:operation GET /api/v1/recipes/{id} recipes
// Get one recipe
// ---
// produces:
// - application/json
// parameters:
//   - name: id
//     in: path
//     description: recipe ID
//     required: true
//     type: string
//
// responses:
//
//	'200':
//	    description: Successful operation
func (handler *RecipesHandler) GetOneRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("id = %s\n", id)
	// objectId, _ := primitive.ObjectIDFromHex(id)
	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"_id": id,
	})
	var recipe models.Recipe
	err := cur.Decode(&recipe)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.HTML(http.StatusOK, "recipe.tmpl", gin.H{"recipe": recipe})
}
