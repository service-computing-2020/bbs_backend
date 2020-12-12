package service

import "github.com/service-computing-2020/bbs_backend/models"

func GetAllPostsByForumID(forum_id int) ([]models.Post, error){
	return models.GetAllPostsByForumID(forum_id)
}