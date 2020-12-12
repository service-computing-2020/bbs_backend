package controllers

import (
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/service-computing-2020/bbs_backend/models"
	"github.com/service-computing-2020/bbs_backend/service"
)

// 用户创建 POST
func CreatePost(c *gin.Context) {
	var data interface{}
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	title, content := c.PostForm("title"), c.PostForm("content")
	user_id := service.GetUserFromContext(c).UserId



	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg":"标题或者内容不得为空" , "data": data})
		return
	}

	var filesToBeUpload []service.File
	for _, fileHeader := range files {
		f, err := fileHeader.Open()
		if err != nil {
		    c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "您所上传的文件无法打开", "data": data})
		    return
		}
		filesToBeUpload = append(filesToBeUpload, service.File{F: f, H: fileHeader})
	}
	bucketName := path.Base(c.Request.URL.Path)
	names, err := service.MultipleFilesUpload(filesToBeUpload, bucketName ,c.Request.URL.Path, ".png" )
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "上传文件失败，服务器内部错误" + err.Error(), "data": data})
	    return
	}
	// 首先插入 post, 获取post_id
	post_id, err := models.CreatePost(models.Post{ForumID: forum_id, UserID: user_id, Title: title, Content: content})
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "插入用户创建的post失败", "data": data})
	    return
	}

	for _, name := range names {
		_, err := models.CreateFile(models.ExtendedFile{PostID: int(post_id), Bucket: bucketName, FileName: name})
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建 post 成功", "data": data})
	return
}