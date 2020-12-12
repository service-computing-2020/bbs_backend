package controllers

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"github.com/service-computing-2020/bbs_backend/models"
	"github.com/service-computing-2020/bbs_backend/service"
)

// 用户创建 HOLE

func CreateHole(c *gin.Context) {
	log.Info("user create hole")
	var data interface{}
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	title, content := c.PostForm("title"), c.PostForm("content")
	user_id := service.GetUserFromContext(c).UserId

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "标题或者内容不得为空", "data": data})
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
	names, err := service.MultipleFilesUpload(filesToBeUpload, bucketName, c.Request.URL.Path, ".png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "上传文件失败，服务器内部错误" + err.Error(), "data": data})
		return
	}
	// 首先插入 hole, 获取hole_id
	hole_id, err := models.CreateHole(models.Hole{ForumID: forum_id, UserID: user_id, Title: title, Content: content})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "插入用户创建的hole失败", "data": data})
		return
	}

	for _, name := range names {
		_, err := models.CreateFile(models.ExtendedFile{PostID: int(hole_id), Bucket: bucketName, FileName: name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库插入异常 " + err.Error(), "data": data})
			return
		}
		fmt.Println(name)
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建 hole 成功", "data": data})
	return
}

// 获取某个 forum 下的所有hole
func GetAllHolesByForumID(c *gin.Context) {
	log.Info("get all holes by forum_id controller")
	forum_id, _ := strconv.Atoi(c.Param("forum_id"))

	data, err := service.GetAllHolesByForumID(forum_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询数据库出现异常" + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("获取论坛 %d 下的全部树洞帖子成功", forum_id),
		"data": data,
	})
}

// 根据 id 获取某个hole的详情
func GetOneHoleDetailByHoleID(c *gin.Context) {
	log.Info("get one hole detail by hole_id")
	hole_id, _ := strconv.Atoi(c.Param("hole_id"))

	data, err := service.GetOneHoleDetailByHoleID(hole_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "数据库查询异常，或者该hole不存在：" + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("获取第 %d 号树洞帖子成功", hole_id),
		"data": data,
	})
}
