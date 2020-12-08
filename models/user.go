package models

import (
	"strconv"
)

type User struct {

	UserId     	int		`json:"user_id"`
	Username    string	`json:"username"`
	Email		string  `json:"email"`
	Password	string	`json:"password"`
	IsAdmin		bool	`json:"is_admin"`
	Avatar		string	`json:"avatar"`
	CreateAt	string	`json:"create_at"`

}

// 将数据库查询的结果转换为 User
func convertMapToUser(user map[string]string) User {
	user_id, _ := strconv.Atoi(user["user_id"])
	is_admin := false
	if user["is_admin"] == "1" {
		is_admin = true
	}
	return User{UserId: user_id, Username: user["username"], Email: user["email"],Password: user["password"], IsAdmin: is_admin, Avatar: user["avatar"], CreateAt: user["create_at"]}
}

// 创建用户
func CreateUser(user User) error{
	sentence := "INSERT INTO user(username, password, email ,is_admin, avatar) VALUES(?, ?, ?, ?, ?)"
	_, err := Execute(sentence, user.Username, user.Password, user.Email, user.IsAdmin, user.Avatar)
	return err
}

// 根据用户id获取用户
func GetUserById(user_id int) ([]User, error) {
	var ret []User

	res, err := QueryRows("SELECT user_id, username, password, email ,is_admin, create_at, avatar FROM user WHERE user_id = ?", user_id)

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, nil
}

// 根据用户名获取用户
func GetUserByUsername(username string) ([]User, error) {
	var ret []User

	res, err := QueryRows("SELECT user_id, username, password, email, is_admin, create_at, avatar FROM user WHERE username = ?", username)

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, err
}

func GetUserByEmail(email string) ([]User, error) {
	var ret []User

	res, err := QueryRows("SELECT user_id, username, password, email,is_admin, create_at, avatar FROM user WHERE email = ?", email)

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, err
}

func GetAllUsers() ([]User, error) {
	var ret []User

	res, err := QueryRows("SELECT user_id, username, password, email,is_admin, create_at, avatar FROM user")

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToUser(r))
	}

	return ret, err
}


