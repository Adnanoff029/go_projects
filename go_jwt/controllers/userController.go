package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Adnanoff029/go_jwt/database"
	helper "github.com/Adnanoff029/go_jwt/helpers"
	"github.com/Adnanoff029/go_jwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("Error loading the env file")
	}
	cost, _ := strconv.Atoi(os.Getenv("COST"))
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPassword)
}

func VerifyPassword(userPassword string, typedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(typedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("Incorrect email or password.")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ctxt, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": validationErr.Error(),
			})
			return
		}

		count, err := userCollection.CountDocuments(ctxt, bson.M{
			"email": user.Email,
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occured while validating the email",
			})
			log.Panic(err)
			return
		}

		pass := HashPassword(*user.Password)
		user.Password = &pass

		count, err = userCollection.CountDocuments(ctxt, bson.M{
			"phone": user.Phone,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error occured while validating the phone number.",
			})
			log.Panic(err)
			return
		}
		if count > 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "This email / phone number already exists.",
			})
			return
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken
		resutlInsertionNumber, insertionErr := userCollection.InsertOne(ctxt, user)
		if insertionErr != nil {
			msg := fmt.Sprintf("User was not registered.")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}
		ctx.JSON(http.StatusOK, resutlInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		err := userCollection.FindOne(ctxt, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		if passwordIsValid != true {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": msg,
			})
			return
		}
		if foundUser.Email == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "User not found.",
			})
			return
		}
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctxt, bson.M{
			"user_id": foundUser.User_id,
		}).Decode(&foundUser)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := helper.CheckUserType(ctx, "ADMIN"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(ctx.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}
		startIdx := (page - 1) * recordPerPage
		startIdx, err = strconv.Atoi(ctx.Query("startIdx"))
		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
			{Key: "total_count", Value: bson.D{
				{Key: "$sum", Value: 1}},
			},
			{Key: "data", Value: bson.D{
				{Key: "$push", Value: "$$ROOT"}},
			},
		}}}
		projectStage := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{
				{Key: "$slice", Value: []any{"$data", startIdx, recordPerPage}},
			}},
		}}}
		result, err := userCollection.Aggregate(ctxt, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while fetching the users.",
			})
			return
		}
		var allUsers []bson.M
		if err = result.All(ctxt, &allUsers); err != nil {
			log.Fatal(err)
		}
		ctx.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")
		if err := helper.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var ctxt, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		err := userCollection.FindOne(ctxt, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, user)
	}
}
