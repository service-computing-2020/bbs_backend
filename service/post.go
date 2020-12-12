package service

import (
	"github.com/pkg/errors"
	"github.com/service-computing-2020/bbs_backend/models"
)

func GetAllPostsByForumID(forum_id int) ([]models.Post, error){
	return models.GetAllPostsByForumID(forum_id)
}

// 根据 post_id 获取一个post的详情
func GetOnePostDetailByPostID(post_id int) ([]models.PostDetail, error) {
	var postDetails []models.PostDetail

	posts, err := models.GetOnePostByPostID(post_id)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return nil, errors.New("该post_id对应的post不存在")
	}
	post := posts[0]

	files, err := models.GetFilesByPostID(post_id)
	if err != nil {
		return nil, err
	}

	postDetail := models.PostDetail{Post: post, Files: files}
	postDetails = []models.PostDetail{postDetail}
	return postDetails, nil
}