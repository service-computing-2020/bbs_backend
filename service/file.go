package service

import "github.com/service-computing-2020/bbs_backend/models"

func GetFilesByPostID(post_id int) ([]models.ExtendedFile, error) {
	files, err := models.GetFilesByPostID(post_id)
	if err != nil {
		return nil, err
	}
	return files, nil
}
