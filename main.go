package main

import (
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/service-computing-2020/bbs_backend/controllers"
	//"github.com/service-computing-2020/bbs_backend/middlewares"
)

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	api := router.Group("/api")
	{
		userRouter := api.Group("/users")
		{

			userRouter.POST("/", controllers.UserRegister)
			userRouter.PUT("/", controllers.UserLogin)
			//userRouter.Use(middlewares.VerifyJWT())
			userRouter.GET("/", controllers.GetAllUsers)

			// 单个用户路由
			singleUserRouter := userRouter.Group("/:user_id")
			{
				singleUserRouter.POST("/avatar", controllers.UploadAvatar)
				singleUserRouter.GET("/avatar", controllers.GetAvatar)
			}
		}

	}
	router.Run(":5000")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}