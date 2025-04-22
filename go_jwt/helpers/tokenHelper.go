package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Adnanoff029/go_jwt/database"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type SignatureDetails struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func GenerateAllTokens(email, firstName, lastName, userType, uid string) (token string, refreshToken string, err error) {
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("Error loading the env file")
	}
	secretKey := os.Getenv("SECRET_KEY")
	claims := &SignatureDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	refreshClaims := &SignatureDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 120)),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllTokens(token, refreshToken, userId string) {
	ctxt, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	filter := bson.M{"user_id": userId}
	update := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    Updated_at,
		},
	}

	_, err := userCollection.UpdateOne(
		ctxt,
		filter,
		update,
		options.UpdateOne().SetUpsert(true),
	)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func ValidateToken(clientToken string) (claims SignatureDetails, msg string) {
	err1 := godotenv.Load()
	if err1 != nil {
		log.Fatal("Error loading the env file")
	}
	secretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(
		clientToken,
		&SignatureDetails{},
		func(token *jwt.Token) (any, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			msg = "Invalid token signature"
			return
		}
		if err == jwt.ErrTokenExpired {
			msg = "Token has expired"
			return
		}
		if err == jwt.ErrTokenNotValidYet {
			msg = "Token is not valid yet"
			return
		}
		msg = "Invalid token"
		return
	}

	claims = *token.Claims.(*SignatureDetails)
	return claims, msg
}
