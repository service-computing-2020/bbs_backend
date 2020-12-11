package controllers

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strconv"

	"github.com/service-computing-2020/bbs_backend/service"

	"github.com/gin-gonic/gin"
	"github.com/service-computing-2020/bbs_backend/models"
)

func GetAllPublicFroums(c *gin.Context) {
	data, err := models.GetAllPublicForums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error(), "data": data})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "获取全部公开论坛", "data": data})
}

// 论坛参数
type ForumParam struct {
	ForumName   string `json:"forum_name"`
	IsPublic    bool   `json:"is_public"`
	Description string `json:"description"`
}

func CreateForum(c *gin.Context) {
	var param ForumParam
	data := make(map[string]string)
	err := c.BindJSON(&param)
	fmt.Println(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数不合法: " + err.Error(), "data": data})
		return
	}
	// TODO:
	// 1.检查论坛同名问题
	// 2.写入forum_user的关系
	err = service.CreateForum(param.ForumName, param.Description, param.IsPublic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "论坛创建成功"})
}

// 上传论坛封面
func UploadCover(c *gin.Context) {
	data := make(map[string]string)
	// TODO:
	// 查表当前的用户是否是论坛的主人或者管理员
	file, header, err := c.Request.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求格式不正确: " + err.Error(), "data": data})
	} else {
		fmt.Println(c.Request.URL.String())
		fmt.Println("base", path.Base(c.Request.URL.Path))
		fmt.Println("not base", c.Request.URL.Path)
		_, err := service.FileUpload(file, header, path.Base(c.Request.URL.Path), c.Request.URL.Path, ".png")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "文件服务错误: " + err.Error(), "data": data})
		} else {
			forumID, err := strconv.Atoi(c.Param("forum_id"))
			if err != nil {
				// id错误?
			}
			// TODO:
			err = models.UpdateCover(c.Request.URL.Path+".png", forumID)
			if err != nil {
				// 更新数据库失败
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "上传封面成功", "data": data})
		}
	}
}

func GetCover(c *gin.Context) {
	var data interface{}
	rawImage, err := service.FileDownload(c.Request.URL.Path, "cover", ".png")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "获取封面失败" + err.Error(), "data": data})
	} else {
		// 图片最多2个M
		image := make([]byte, 2000000)
		len, err := rawImage.Read(image)
		if err != nil {
			if err != io.EOF && err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "读取图片失败 " + err.Error(), "data": data})
			} else {
				c.Writer.WriteHeader(http.StatusOK)
				c.Header("Content-Disposition", "attachment; filename=hello.txt")
				c.Header("Content-Type", "image/jpeg")
				c.Header("Accept-Length", fmt.Sprintf("%d", len))
				c.Writer.Write(image)
			}
		}
	}
}
