package handlers

import (
	"context"
	"net/http"

	"cinema-booking/internal/models"
	"cinema-booking/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieHandler struct {
	Mongo *services.MongoService
}

func (h *MovieHandler) ListMovies(c *gin.Context) {
	var movies []models.Movie
	cursor, err := h.Mongo.Collection("movies").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch movies"})
		return
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &movies); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode movies"})
		return
	}

	if movies == nil {
		movies = []models.Movie{}
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var movie models.Movie
	err = h.Mongo.Collection("movies").FindOne(context.Background(), bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}
