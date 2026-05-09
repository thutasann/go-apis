package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thutasann/restaurant-api/database"
	"github.com/thutasann/restaurant-api/helpers"
	"github.com/thutasann/restaurant-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// mongo user collection
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// Get all users
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage

		matchStage := bson.D{{Key: "$match", Value: bson.D{}}}
		projectStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{
				{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
			}},
		}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			projectStage,
		})

		if err != nil {
			helpers.Error(c, "failed to get users", 0, err)
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		helpers.Success(c, "get users success", allUsers[0])
	}
}

// Get User by Id
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Param("user_id")

		var user models.User

		err := userCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: userId}}).Decode(&user)
		if err != nil {
			helpers.Error(c, "failed to get user with userId", 0, err)
			return
		}
		helpers.Success(c, "get user with userId success", user)
	}
}

// Sign Up the user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		// convert the json data
		if err := c.BindJSON(&user); err != nil {
			helpers.Error(c, "failed to sign up the user", 0, err)
			return
		}

		// validate the data based on use struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			helpers.Error(c, "validation error", 400, validationErr)
			return
		}

		// check if the email has already been used by another user
		emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			helpers.Error(c, "user already existed", 0, err)
			return
		}

		if emailCount > 0 {
			helpers.Error(c, "user with this email already existed", 0, err)
			return
		}

		// hash password
		password := hashPassword(*user.Password)
		user.Password = &password

		// check if the phone no. has already been used by another user
		phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			helpers.Error(c, "user with this phone number already existed", 0, err)
			return
		}

		if phoneCount > 0 {
			helpers.Error(c, "user with this phone number already existed", 0, err)
			return
		}

		// create some extra details for the user object - created_at, updated_at
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// generate token and refresh token (generate all tokens functions from the helper)
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_Token = &refreshToken

		// insert this new user into the user collection
		resultInsert, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created: %s", insertErr)
			helpers.Error(c, msg, 0, err)
			return
		}

		// return status OK and send the result back
		helpers.Success(c, "sign up success", resultInsert)
	}
}

// Sign In the user
func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		// convert the login data from client which is in JSON to golang readable format
		if err := c.BindJSON(&user); err != nil {
			helpers.Error(c, "json bind error", 0, err)
			return
		}

		// find a user with that email and see if that user even exists
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			helpers.Error(c, "user not found", 404, err)
			return
		}

		// then you will verify the password
		passwordIsValid, msg := verifyPassword(*user.Password, *foundUser.Password)
		if !passwordIsValid {
			helpers.Error(c, msg, 400, err)
			return
		}

		// if all goes well, then you'll generate tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		// update tokens - token and refresh token
		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		// return statusOK
		helpers.Success(c, "user sign in success", foundUser)
	}
}

// Private: hash password
func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// Private: verify passwrod
func verifyPassword(userPassword string, providePassword string) (bool, string) {

	err := bcrypt.CompareHashAndPassword([]byte(providePassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or password is incorrect: %s", err)
		check = false
	}

	return check, msg
}
