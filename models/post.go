package models

import "strconv"

type Post struct {
	PostID   int    `json:"post_id"`
	ForumID  int    `json:"forum_id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt string `json:"create_at"`
	Like     int    `json:"like"`
}

type PostDetail struct {
	Post
	Files []ExtendedFile
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
	res, err := QueryRows("SELECT * FROM post WHERE forum_id=?", forum_id)
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
