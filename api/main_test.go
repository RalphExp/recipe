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

func TestUpdateRecipeHandler(t *testing.T) {
	r := SetupRouter()
	r.PUT("/api/v1/recipes/:id", recipesHandler.UpdateRecipeHandler)

	recipe := models.Recipe{
		ID:   "603fa1203e8ae186e9427f98",
		Name: "Mic's Yorkshire Puds",
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

func TestDeleteRecipeHandler(t *testing.T) {
	r := SetupRouter()
	r.DELETE("/api/v1/recipes/:id", recipesHandler.DeleteRecipeHandler)
	req, _ := http.NewRequest("DELETE", "/api/v1/recipes/0000000000", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.Bytes()
	var payload map[string]string
	json.Unmarshal(body, &payload)

	assert.Equal(t, payload["message"], "Recipe has been deleted")
}
