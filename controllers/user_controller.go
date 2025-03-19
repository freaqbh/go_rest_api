package controllers

import (
	"context"
	"net/http"
	"time"

	"rest_api/database"
	"rest_api/middlewares"
	"rest_api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginReq struct {
	Username string `bson:"usernname" json:"username" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required,min=6"`
}

func GetUserByID(c *gin.Context) {
	// Get the user ID from the URL parameter
	id := c.Param("id")
	collection := database.GetCollection("users")

	// convert string id to primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
	}

	// make ctx with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find user by id
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
	}

	c.JSON(http.StatusOK, user)
}

func GetUser(c *gin.Context) {
	// make ctx with timeout
	collection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all users
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data"})
		return
	}

	defer cursor.Close(ctx)

	var users []models.User
	// Loop through the cursor and append each document to the slice of users
	for cursor.Next(ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data"})
			return
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	// get json data from post request
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = primitive.NewObjectID()

	// hash password
	if err := user.HashPassword(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengubah password"})
		return
	}

	// Insert user into database
	collection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func LoginHandler(c *gin.Context) {
	// get json data from post request
	var req LoginReq

	// bind json data to req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := database.GetCollection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// verification password
	if !user.VerifyPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username atau password salah 2"})
	}

	// generate token
	token, err := middlewares.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat token"})
		return
	}

	// send token
	c.JSON(http.StatusOK, gin.H{"token": token})
}
