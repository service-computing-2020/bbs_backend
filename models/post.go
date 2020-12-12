package models

import "strconv"

type Post struct {
	PostID	int		`json:"post_id"`
	ForumID	int		`json:"forum_id"`
	UserID  int		`json:"user_id"`
	Title	string	`json:"title"`
	Content string	`json:"content"`
	CreateAt	string	`json:"content"`
	Like	int	`json:"like"`
}

// 将数据库查询结果转换为 POST
func convertMapToPost(post map[string]string) Post {
	post_id, _ := strconv.Atoi(post["post_id"])
	forum_id, _ := strconv.Atoi(post["forum_id"])
	user_id, _ := strconv.Atoi(post["user_id"])
	like, _ := strconv.Atoi(post["like"])
	return Post{
		PostID: post_id,
		ForumID: forum_id,
		UserID: user_id,
		Title: post["title"],
		Content: post["content"],
		CreateAt: post["create_at"],
		Like: like,
	}
}

// 创建一个帖子
func CreatePost(post Post) (int64, error) {
	sentence := "INSERT INTO post(forum_id, user_id, title, content) VALUES (?, ?, ?, ?)"
	return Execute(sentence, post.ForumID, post.UserID, post.Title, post.Content)
}
