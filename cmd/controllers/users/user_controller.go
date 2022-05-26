package users

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ko1ke/know-sync-api/cmd/domain/users"
	"github.com/ko1ke/know-sync-api/cmd/services"
	"github.com/ko1ke/know-sync-api/cmd/utils/auth_utils"
	"github.com/ko1ke/know-sync-api/cmd/utils/pg_error_utils"
	"github.com/ko1ke/know-sync-api/cmd/utils/res_utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func SignUp(c *gin.Context) {
	var user users.User
	var err error

	if err = c.ShouldBindJSON(&user); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnprocessableEntity, &res_utils.ErrObj{Message: err.Error()})
		return
	}

	user.Password, err = auth_utils.GeneratehashPassword(user.Password)
	if err != nil {
		logrus.Fatalln("Error in password hashing.")
	}

	newUser, saveErr := services.CreateUser(user)
	if saveErr != nil {
		logrus.Error(saveErr)
		c.JSON(http.StatusUnprocessableEntity, &res_utils.ErrObj{Message: pg_error_utils.ParseError(saveErr).Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newUser.ID, "username": newUser.Username, "email": newUser.Email})
}

func SignIn(c *gin.Context) {
	var auth users.Authentication
	if err := c.ShouldBindJSON(&auth); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnprocessableEntity, &res_utils.ErrObj{Message: err.Error()})
		return
	}

	user, err := services.GetUserByEmail(auth.Email)

	if user == nil || err != nil {
		logrus.Error("email is not registered")
		c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: "Eメールが登録されていません"})
		return
	}

	check := auth_utils.CheckPasswordHash(auth.Password, user.Password)

	if !check {
		logrus.Error("invalid password")
		c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: "パスワードが不正です"})
		return
	}

	ts, err := auth_utils.CreateToken(user.ID)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	err = auth_utils.CreateAuth(user.ID, ts)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	tokens := gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"accessToken":  ts.AccessToken,
		"refreshToken": ts.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

func getUserID(userIDParam string) (uint, error) {
	userID, userErr := strconv.ParseUint(userIDParam, 10, 64)
	if userErr != nil {
		return 0, userErr
	}
	return uint(userID), nil
}

func GetUserFromToken(r *http.Request) (*users.User, error) {
	metadata, err := auth_utils.ExtractTokenMetadata(r)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	userId, err := auth_utils.FetchAuth(metadata)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	user, _ := services.GetUser(userId)
	return user, nil
}

func AuthUser(c *gin.Context) {
	user, err := GetUserFromToken(c.Request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username, "email": user.Email})
}

func SignOut(c *gin.Context) {
	au, err := auth_utils.ExtractTokenMetadata(c.Request)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: "トークンが不正です"})
		return
	}
	deleted, delErr := auth_utils.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		logrus.Error(delErr)
		c.JSON(http.StatusInternalServerError, &res_utils.ErrObj{Message: "エラーが発生しました"})
		return
	}
	c.JSON(http.StatusOK, "ログアウトしました")
}

func Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnprocessableEntity, &res_utils.ErrObj{Message: err.Error()})
		return
	}
	refreshToken := mapToken["refreshToken"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: "ログインしてください"})
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		logrus.Error(err)
		c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: "ログインしてください"})
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			logrus.Error(err)
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnprocessableEntity, &res_utils.ErrObj{Message: "エラーが発生しました"})
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := auth_utils.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			logrus.Error(delErr)
			c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: "エラーが発生しました"})
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := auth_utils.CreateToken(uint(userId))
		if createErr != nil {
			logrus.Error(createErr)
			c.JSON(http.StatusForbidden, &res_utils.ErrObj{Message: "エラーが発生しました"})
			return
		}
		//save the tokens metadata to redis
		saveErr := auth_utils.CreateAuth(uint(userId), ts)
		if saveErr != nil {
			logrus.Error(saveErr)
			c.JSON(http.StatusForbidden, &res_utils.ErrObj{Message: "エラーが発生しました"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"accessToken":  ts.AccessToken,
			"refreshToken": ts.RefreshToken,
		})
	} else {
		c.JSON(http.StatusUnauthorized, &res_utils.ErrObj{Message: "ログインしてください"})
	}
}

func DeleteUser(c *gin.Context) {
	userID, idErr := getUserID(c.Param("user_id"))
	if idErr != nil {
		logrus.Error(idErr)
		c.JSON(http.StatusNotFound, "ユーザーが存在しません")
		return
	}

	if err := services.DeleteUser(userID); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": userID})
}
