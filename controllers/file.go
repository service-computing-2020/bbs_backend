package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"github.com/service-computing-2020/bbs_backend/service"
)

func GetFilesByPostID(c *gin.Context) {
	log.Info("get files by post_id")
	post_id, _ := strconv.Atoi(c.Param("post_id"))

	files, err := service.GetFilesByPostID(post_id)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据库查询错误 " + err.Error(), "data": nil})
	    return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg": fmt.Sprintf("获取第 %d 帖子下的所有文件成功", post_id),
		"data": files,
	})
}
