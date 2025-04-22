package helpers

import(
	"github.com/gin-gonic/gin"
	"errors"
)

func CheckUserType(ctx *gin.Context, role string)  (err error) {
	userType := ctx.GetString("user_type")
	err = nil
	if userType != role {
		err = errors.New("Unauthorized to access the resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(ctx *gin.Context, userId string) (err error) {
	userType := ctx.GetString("user_type") // ADMIN USER
	uid := ctx.GetString("uid")
	err = nil
	// If a USER access the resource but his ID is not equal to the entered ID
	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access the resource")
		return err
	}
	// USER and uid == userId // err would still be nil after the below function call
	// ADMIN
	// err = CheckUserType(ctx, userType) // Checks if the user type of the current user
	return err
}

