package models

import (
	"strconv"
)

type Forum struct {
	ForumID     int    `json:"forum_id"`
	ForumName   string `json:"forum_name"`
	IsPublic    bool   `json:"is_public"`
	Description string `json:"description"`
	CreateAt    string `json:"create_at"`
	Cover       string `json:"cover"`
}

// 包含 forum 信息的用户
type ForumUser struct {
	User
	Role    string `json:"role"`
	ForumID int    `json:"forum_id"`
}

// 将数据库查询结果转换为Forum
func convertMapToForum(forum map[string]string) Forum {
	forum_id, _ := strconv.Atoi(forum["forum_id"])
	is_public := true
	if forum["is_public"] == "0" {
		is_public = false
	}
	return Forum{
		ForumID:     forum_id,
		ForumName:   forum["forum_name"],
		IsPublic:    is_public,
		Description: forum["description"],
		CreateAt:    forum["create_at"],
		Cover:       forum["cover"],
	}
}

func convertMapToForumUser(forumUser map[string]string) ForumUser {
	forum_id, _ := strconv.Atoi(forumUser["forum_id"])
	return ForumUser{
		ForumID: forum_id,
		Role:    forumUser["role"],
		User:    convertMapToUser(forumUser),
	}
}

// 创建论坛
func CreateForum(forum Forum) error {
	sentence := "INSERT INTO forum(forum_name, is_public, description, cover) VALUES(?, ?, ?, ?)"
	_, err := Execute(sentence, forum.ForumName, forum.IsPublic, forum.Description, forum.Cover)
	return err
}

// 往论坛写入封面的url
func UpdateCover(path string, forumID int) error {
	sentence := "UPDATE forum SET cover=? WHERE forum_id=?"
	_, err := Execute(sentence, path, forumID)
	return err
}

// 获取论坛封面的url
func GetCoverURL(forumID int) error {
	sentence := "SELECT cover FROM forum WHERE user_id=?"
	_, err := Execute(sentence, forumID)
	return err
}

// 获取所有公开论坛，用于未登陆时获取公开论坛
func GetAllPublicForums() ([]Forum, error) {
	var ret []Forum

	res, err := QueryRows("SELECT forum_id, forum_name, is_public, description, create_at, cover FROM forum WHERE is_public = 1")

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToForum(r))
	}

	return ret, nil
}

// 根据论坛的 ID 获取论坛信息
func GetForumByID(forum_id int) ([]Forum, error) {
	var ret []Forum

	res, err := QueryRows("SELECT * FROM forum WHERE forum_id=?", forum_id)
	if err != nil {
		return nil, err
	}

	for _, r := range res {
		ret = append(ret, convertMapToForum(r))
	}

	return ret, nil
}

// 根据论坛id获取其成员信息/关注者信息
func GetForumUserByForumID(forum_id int) ([]ForumUser, error) {
	var users []ForumUser

	sentence :=
		`
SELECT user.user_id, user.username, user.email, user.is_admin, user.create_at, user.avatar, forum_user.role
			FROM user INNER JOIN forum_user
            ON user.user_id = forum_user.user_id
			WHERE forum_user.forum_id = ?;
		`
	res, err := QueryRows(sentence, forum_id)
	if err != nil {
		return nil, err
	}
	for _, u := range res {
		users = append(users, convertMapToForumUser(u))
	}

	return users, nil
}

// 根据用户id获取公开及其相关论坛
