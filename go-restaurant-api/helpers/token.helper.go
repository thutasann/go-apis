package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/thutasann/restaurant-api/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// JWT Signed Details
type SignedDetails struct {
	Email              string // User's email name
	First_name         string // User's first name
	Last_name          string // User's last name
	Uid                string // User's user id
	jwt.StandardClaims        // JWT's standard claims
}

// mongo user collection
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// JWT Secret key from .env
var SECRET_KEY string = os.Getenv("SECRET_KEY")

// Generate token, refrsh tokens
func GenerateAllTokens(email, first_name, last_name, user_id string) (signedToken, refreshToke string, err error) {

	claims := &SignedDetails{
		Email:      email,
		First_name: first_name,
		Last_name:  last_name,
		Uid:        user_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
	}

	refreshToken, refreshError := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if refreshError != nil {
		log.Panic(err)
	}

	return token, refreshToken, err
}

// Update token, refresh tokens
func UpdateAllTokens(token, refrshToken, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: token})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: refrshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	if err != nil {
		log.Panic(err)
		return
	}
}

// Validate Token
func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		msg = "Failed to parse token: " + err.Error()
		return nil, msg
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok || !token.Valid {
		msg = "Token is invalid"
		return nil, msg
	}

	if claims.ExpiresAt < time.Now().Unix() {
		msg = "Token is expired"
		return nil, msg
	}

	return claims, ""
}
