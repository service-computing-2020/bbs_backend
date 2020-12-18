package service

import (
	"github.com/pkg/errors"
	"github.com/service-computing-2020/bbs_backend/models"
	"log"
)

func IsUsernameExist(username string) (bool, error) {
	users, err := models.GetUserByUsername(username)
	if err != nil {
		return false, err
	}
	if len(users) > 0 {
		return true, nil
	}
	return false, nil
}

func IsEmailExist(email string) (bool, error) {
	users, err := models.GetUserByEmail(email)
	log.Println(len(users))
	if err != nil {
		return false, err
	}
	if len(users) > 0 {
		return true, nil
	}
	return false, nil
}

func VerifyByUsernameAndPassword(username string, password string) (bool, error) {
	users, err := models.GetUserByUsername(username)
	if err != nil {
		return false, err
	}

	if users[0].Password == password {
		return true, nil
	}
	return false, nil

}

func VerifyByEmailAndPassword(email string, password string) (bool, error) {
	users, err := models.GetUserByEmail(email)
	if err != nil {
		return false, err
	}

	if users[0].Password == password {
		return true, nil
	}
	return false, nil

}

func CreateUser(username string, password string, email string) error{
	user := models.User{Username: username, Email: email, Password: password, Avatar: "0.jpg"}
	return models.CreateUser(user)
}

func ProduceTokenByUsernameAndPasword(username string, password string) (string, error) {
	users, err := models.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	token, err := GenerateToken(users[0].UserId, username, password)
	if err != nil {
		return "", err
	}
	return token, nil
}
func ProduceTokenByEmailAndPassword(email string, password string) (string, error) {
	users, err := models.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	token, err := GenerateToken(users[0].UserId, email, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetOneUserSubscribe(userID int) (models.SubscribeList, error) {
	return models.GetOneUserSubscribe(userID)
}

func GetOneUserDetail(userID int) (models.UserDetail, error) {
	var userDetail models.UserDetail
	user, err := models.GetUserById(userID)
	if err != nil {
		return userDetail, err
	}

	if len(user) == 0 {
		return userDetail, errors.New("该用户不存在")
	}
	subscribe, err := models.GetOneUserSubscribe(userID)
	if err != nil {
		return userDetail, err
	}

	likeList, err := models.GetOneUserLikeListByUserID(userID)
	if err != nil {
		return userDetail, err
	}
	userDetail.User = user[0]
	userDetail.SubscribeList = subscribe
	userDetail.LikeList = likeList

	return userDetail, nil
}






