package models

import "strconv"

type Hole struct {
	HoleID   int    `json:"hole_id"`
	ForumID  int    `json:"forum_id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	CreateAt string `json:"content"`
}

// 将数据库查询结果转换为 HOLE
func convertMapToHole(hole map[string]string) Hole {
	hole_id, _ := strconv.Atoi(hole["hole_id"])
	forum_id, _ := strconv.Atoi(hole["forum_id"])
	user_id, _ := strconv.Atoi(hole["user_id"])
	return Hole{
		HoleID:   hole_id,
		ForumID:  forum_id,
		UserID:   user_id,
		Title:    hole["title"],
		Content:  hole["content"],
		CreateAt: hole["create_at"],
	}
}

// 创建一个匿名帖子(树洞帖子)
func CreateHole(hole Hole) (int64, error) {
	sentence := "INSERT INTO hole(forum_id, user_id, title, content) VALUES (?, ?, ?, ?)"
	return Execute(sentence, hole.ForumID, hole.UserID, hole.Title, hole.Content)
}
