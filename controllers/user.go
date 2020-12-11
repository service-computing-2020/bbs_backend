package controllers

import (
	"fmt"
	"github.com/service-computing-2020/bbs_backend/models"
	"github.com/service-computing-2020/bbs_backend/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"path"
)

// 用户注册需要提供的字段
type RegisterParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email	 string `json:"email"`
}
// 用户注册控制器
func UserRegister(c *gin.Context) {
	var param RegisterParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数不合法: "+err.Error()})
		return
	}

	if ok, err := service.IsUsernameExist(param.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error()})
		return
	} else if ok {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "该用户名已经被使用"})
		return
	}

	if ok, err := service.IsEmailExist(param.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error()})
		return
	} else if ok {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "该邮箱已经被使用"})
		return
	}

	err = service.CreateUser(param.Username, param.Password, param.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500 ,"msg": "服务器错误: "+err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功"})
}


// 登陆参数
type LoginParam struct {
	Input 	 string  `json:"input"`
	Password string	 `json:"password"`
}
// 用户登录控制器
func UserLogin(c *gin.Context) {
	var param LoginParam
	data := make(map[string]string)
	err := c.BindJSON(&param)
	fmt.Println(param)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数不合法: "+err.Error(), "data":data})
		return
	}
	if ok, err := service.IsUsernameExist(param.Input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error(), "data":data})
		return
	} else if ok {
		pass, err := service.VerifyByUsernameAndPassword(param.Input, param.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error(), "data":data})
			return
		}
		if pass {
			data["token"], err = service.ProduceTokenByUsernameAndPasword(param.Input, param.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error(), "data":data})
				return
			}

			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data":data})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "密码错误", "data":data})
			return
		}
	}

	if ok, err := service.IsEmailExist(param.Input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error(), "data":data})
		return
	} else if ok {
		pass, err := service.VerifyByEmailAndPassword(param.Input, param.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error(), "data":data})
			return
		}
		if pass {
			data["token"], err = service.ProduceTokenByEmailAndPassword(param.Input, param.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error(), "data":data})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登录成功", "data":data})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "msg": "密码错误", "data":data})
			return
		}
	}

	c.JSON(http.StatusForbidden, gin.H{"code": 403,"msg":"该用户名或邮箱不存在", "data":data})
}

func GetAllUsers(c *gin.Context) {

	data, err := models.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "服务器错误: "+err.Error(), "data":data})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code":200, "msg": "获取全部用户", "data":data})
}


// 上传用户头像图像
func UploadAvatar(c * gin.Context) {
	data := make(map[string]string)

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求格式不正确: " + err.Error(), "data": data})
	} else {
		fmt.Println(c.Request.URL.String())
		// 图片统一改成png上传
		_, err := service.FileUpload(file, header, path.Base(c.Request.URL.Path), c.Request.URL.Path, ".png")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code":500, "msg": "文件服务错误: " + err.Error(), "data": data})
		} else {

			c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "上传头像成功", "data": data})
		}
	}
}

// 获取用户图片
func GetAvatar(c *gin.Context) {
	var data interface{}
	rawImage, err := service.FileDownload(c.Request.URL.Path, "avatar", ".png")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "获取头像失败" + err.Error(), "data": data})

	} else {
		// 图片最多2个M
		image := make([]byte, 2000000 )
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