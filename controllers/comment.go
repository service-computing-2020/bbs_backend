package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"github.com/service-computing-2020/bbs_backend/models"
	"github.com/service-computing-2020/bbs_backend/service"
)

// 用户创建 COMMENT
// CreateComment godoc
// @Summary CreateComment
// @Description	CreateComment
// @Tags Comments
// @Accept	mpfd
// @Produce	json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Param content formData string true "Comment 的内容"
// @Success 200 {object} responses.StatusOKResponse "创建 comment 成功"
// @Failure 400 {object} responses.StatusBadRequestResponse "评论的内容不得为空"
// @Failure 500 {object} responses.StatusInternalServerError "插入用户创建的comment失败"
// @Router /forums/{forum_id}/posts/{post_id}/comments [post]
func CreateComment(c *gin.Context) {
	log.Info("user create comment")
	var data interface{}
	post_id, _ := strconv.Atoi(c.Param("post_id"))
	content := c.PostForm("content")
	user_id := service.GetUserFromContext(c).UserId
	if content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "评论的内容不得为空", "data": data})
		return
	}

	_, err := models.CreateComment(models.Comment{PostID: post_id, UserID: user_id, Content: content}) //CommentID与CreateAt由数据库生成
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "插入用户创建的comment失败", "data": data})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "创建 comment 成功", "data": data})
	return
}

// 获取某个 post 下的所有comment
// GetAllCommentsByPostID godoc
// @Summary GetAllCommentsByPostID
// @Description GetAllCommentsByPostID
// @Tags Comments
// @Accept json
// @Produce json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.Comment} "获取评论成功"
// @Failure 500 {object} responses.StatusInternalServerError "查询数据库出现异常"
// @Router /forums/{forum_id}/posts/{post_id}/comments [get]
func GetAllCommentsByPostID(c *gin.Context) {
	log.Info("get all comments by post_id controller")
	post_id, _ := strconv.Atoi(c.Param("post_id"))

	data, err := service.GetAllCommentsByPostID(post_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询数据库出现异常" + err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("获取帖子 %d 下的全部评论成功", post_id),
		"data": data,
	})
}

// 根据 id 获取某个comment的详情
// GetOneCommentDetailByCommentID godoc
// @Summary GetOneCommentDetailByCommentID
// @Description GetOneCommentDetailByCommentID
// @Tags Comments
// @Accept json
// @Produce json
// @Param token header string true "将token放在请求头部的‘Authorization‘字段中，并以‘Bearer ‘开头""
// @Success 200 {object} responses.StatusOKResponse{data=[]models.CommentDetail} "获取评论成功"
// @Failure 400 {object} responses.StatusInternalServerError "数据库查询异常，或者该comment不存在"
// @Router /forums/{forum_id}/posts/{post_id}/comments/{comment_id} [get]
func GetOneCommentDetailByCommentID(c *gin.Context) {
	log.Info("get one comment detail by comment_id")
	comment_id, _ := strconv.Atoi(c.Param("comment_id"))

	data, err := service.GetOneCommentDetailByCommentID(comment_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 403, "msg": "数据库查询异常，或者该comment不存在：" + err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("获取第 %d 号评论成功", comment_id),
		"data": data,
	})
}
