package handlers

import (
	"context"
	"net/http"
	"time"

	"cinema-booking/internal/middleware"
	"cinema-booking/internal/models"
	"cinema-booking/internal/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthHandler struct {
	Mongo    *services.MongoService
	JWTSecret string
	Firebase *services.FirebaseAuth
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Demo  bool   `json:"demo"`
	Role  string `json:"role"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Demo mode login
	role := models.RoleUser
	if req.Role == "ADMIN" {
		role = models.RoleAdmin
	}

	now := time.Now()
	providerSub := "demo:" + req.Email

	filter := bson.M{"provider_sub": providerSub}
	update := bson.M{
		"$set": bson.M{
			"email":      req.Email,
			"name":       req.Name,
			"role":       role,
			"updated_at": now,
		},
		"$setOnInsert": bson.M{
			"_id":           primitive.NewObjectID(),
			"auth_provider": "demo",
			"provider_sub":  providerSub,
			"created_at":    now,
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var user models.User
	err := h.Mongo.Collection("users").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := middleware.GenerateJWT(user.ID.Hex(), user.Role, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID.Hex(),
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
		"token": token,
	})
}

type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.Firebase == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Firebase auth not configured"})
		return
	}

	claims, err := h.Firebase.VerifyIDToken(req.IDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Google token"})
		return
	}

	now := time.Now()
	providerSub := "google:" + claims.UserID

	filter := bson.M{"provider_sub": providerSub}
	update := bson.M{
		"$set": bson.M{
			"email":      claims.Email,
			"name":       claims.Name,
			"updated_at": now,
		},
		"$setOnInsert": bson.M{
			"_id":           primitive.NewObjectID(),
			"auth_provider": "google",
			"provider_sub":  providerSub,
			"role":          models.RoleUser,
			"created_at":    now,
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var user models.User
	err = h.Mongo.Collection("users").FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := middleware.GenerateJWT(user.ID.Hex(), user.Role, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID.Hex(),
			"email": user.Email,
			"name":  user.Name,
			"role":  user.Role,
		},
		"token": token,
	})
}
