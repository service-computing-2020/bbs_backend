package models

import (
	"errors"
	"strconv"
)

type Post struct {
	PostID   int    `json:"post_id"`
	ForumID  int    `json:"forum_id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt string `json:"create_at"`
	Like     int    `json:"like"`
	Username string `json:"username"`
}

type PostDetail struct {
	Post
	Files []ExtendedFile		`json:"files"`
}

// 将数据库查询结果转换为 POST
func convertMapToPost(post map[string]string) Post {
	post_id, _ := strconv.Atoi(post["post_id"])
	forum_id, _ := strconv.Atoi(post["forum_id"])
	user_id, _ := strconv.Atoi(post["user_id"])
	like, _ := strconv.Atoi(post["like"])
	return Post{
		PostID:   post_id,
		ForumID:  forum_id,
		UserID:   user_id,
		Title:    post["title"],
		Content:  post["content"],
		CreateAt: post["create_at"],
		Like:     like,
		Username: post["username"],
	}
}

// 创建一个帖子
func CreatePost(post Post) (int64, error) {
	sentence := "INSERT INTO post(forum_id, user_id, title, content) VALUES (?, ?, ?, ?)"
	return Execute(sentence, post.ForumID, post.UserID, post.Title, post.Content)
}

// 获取某个 forum 下的全部 posts
func GetAllPostsByForumID(forum_id int) ([]Post, error) {
	var ret []Post
	res, err := QueryRows("SELECT post.post_id, post.forum_id, post.user_id, post.title, post.content, post.create_at, post.like, user.username \n\tFROM post INNER JOIN user ON post.user_id = user.user_id WHERE post.forum_id=? ORDER BY post.create_at DESC", forum_id)
	if err != nil {
		return ret, err
	}

	for _, p := range res {

		ret = append(ret, convertMapToPost(p))
	}
	return ret, nil
}

// 根据id获取某个 Post
func GetOnePostByPostID(post_id int) ([]Post, error) {
	var ret []Post
	res, err := QueryRows("SELECT * FROM post WHERE post_id=?", post_id)
	if err != nil {
		return ret, err
	}

	for _, p := range res {
		ret = append(ret, convertMapToPost(p))
	}
	return ret, nil
}

func LikeOnePostByUserIDAndPostID(userID int, postID int) error {
	sql1 :=
		`
			SELECT * FROM post_like WHERE user_id = ? AND post_id = ?;
		`
	ret, err := QueryRows(sql1, userID, postID)
	if err != nil {
		return err
	}
	if len(ret) > 0 {
		return errors.New("您已经点赞过了")
	}

	sql2 :=
		`
			INSERT INTO post_like(user_id, post_id) VALUES (?, ?);
		`
	_, err = Execute(sql2, userID, postID)
	if err != nil {
		return err
	}

	sql3 :=
		`
			UPDATE post SET post.like = post.like+1 WHERE post.post_id = ?
		`
	_, err = Execute(sql3, postID)
	if err != nil {
		return err
	}
	return nil
}
