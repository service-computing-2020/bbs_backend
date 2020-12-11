package service

import (
	"github.com/service-computing-2020/bbs_backend/models"
)

// 创建论坛
func CreateForum(forumName, description string, isPublic bool) error {
	forum := models.Forum{
		ForumName:   forumName,
		IsPublic:    isPublic,
		Description: description,
	}
	return models.CreateForum(forum)
}
