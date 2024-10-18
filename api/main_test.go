package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ralphexp/recipes-api/models"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestListRecipesHandler(t *testing.T) {
	r := SetupRouter()
	r.GET("/api/v1/recipes", recipesHandler.ListRecipesHandler)
	req, _ := http.NewRequest("GET", "/api/v1/recipes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.Bytes()
	var recipes []models.Recipe
	json.Unmarshal(body, &recipes)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.GreaterOrEqual(t, len(recipes), 9)
}

func TestDeleteRecipeHandler(t *testing.T) {
	r := SetupRouter()
	r.DELETE("/api/v1/recipes/:id", recipesHandler.DeleteRecipeHandler)
	req, _ := http.NewRequest("DELETE", "/api/v1/recipes/603fa1203e8ae186e9427f98", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.Bytes()
	var payload map[string]string
	json.Unmarshal(body, &payload)

	assert.Equal(t, payload["message"], "Recipe has been deleted")
}

func TestInsertRecipeHandler(t *testing.T) {
	r := SetupRouter()
	r.POST("/api/v1/recipes", recipesHandler.NewRecipeHandler)

	recipe := models.Recipe{
		ID:   "603fa1203e8ae186e9427f98",
		Name: "Mic's Yorkshire Puds2",
		Ingredients: []models.Ingredient{{
			Quantity: "200g",
			Name:     "plain flour",
			Type:     "Baking",
		}, {
			Quantity: "3",
			Name:     "eggs",
			Type:     "Dairy",
		}, {
			Quantity: "300ml",
			Name:     "milk",
			Type:     "Dairy",
		}, {
			Quantity: "3 tbsp",
			Name:     "vegetable oil",
			Type:     "Condiments",
		}},
		Steps: []string{
			"Put the flour and some seasoning into a large bowl.",
			"Stir in eggs, one at a time.",
			"Whisk in milk until you have a smooth batter.",
			"Chill in the fridge for at least 30 minutes.",
			"Heat oven to 220C/gas mark 7.",
			"Pour the oil into the holes of a 8-hole muffin tin.",
			"Heat tin in the oven for 5 minutes.",
			"Ladle the batter mix into the tin.",
			"Bake for 30 minutes until well browned and risen.",
		},
		ImageURL: "/assets/images/yorkshire_pudding.jpg",
	}

	jsonValue, _ := json.Marshal(recipe)
	req, _ := http.NewRequest("POST", "/api/v1/recipes",
		bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.Bytes()
	var payload map[string]string
	json.Unmarshal(body, &payload)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateRecipeHandler(t *testing.T) {
	r := SetupRouter()
	r.PUT("/api/v1/recipes/:id", recipesHandler.UpdateRecipeHandler)

	recipe := models.Recipe{
		ID:   "603fa1203e8ae186e9427f98",
		Name: "Mic's Yorkshire Puds",
		Ingredients: []models.Ingredient{{
			Quantity: "200g",
			Name:     "plain flour",
			Type:     "Baking",
		}, {
			Quantity: "3",
			Name:     "eggs",
			Type:     "Dairy",
		}, {
			Quantity: "300ml",
			Name:     "milk",
			Type:     "Dairy",
		}, {
			Quantity: "3 tbsp",
			Name:     "vegetable oil",
			Type:     "Condiments",
		}},
		Steps: []string{
			"Put the flour and some seasoning into a large bowl.",
			"Stir in eggs, one at a time.",
			"Whisk in milk until you have a smooth batter.",
			"Chill in the fridge for at least 30 minutes.",
			"Heat oven to 220C/gas mark 7.",
			"Pour the oil into the holes of a 8-hole muffin tin.",
			"Heat tin in the oven for 5 minutes.",
			"Ladle the batter mix into the tin.",
			"Bake for 30 minutes until well browned and risen.",
		},
		ImageURL: "/assets/images/yorkshire_pudding.jpg",
	}

	jsonValue, _ := json.Marshal(recipe)
	req, _ := http.NewRequest("PUT", "/api/v1/recipes/603fa1203e8ae186e9427f98",
		bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.Bytes()
	var payload map[string]string
	json.Unmarshal(body, &payload)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, payload["message"], "Recipe has been updated")
}
