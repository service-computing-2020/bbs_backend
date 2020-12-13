package main

import (
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/service-computing-2020/bbs_backend/controllers"
	"github.com/service-computing-2020/bbs_backend/middlewares"
	//"github.com/service-computing-2020/bbs_backend/middlewares"
)

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.MaxMultipartMemory = 5 << 20 // 限制文件大小为5MB
	api := router.Group("/api")
	{
		userRouter := api.Group("/users")
		{

			userRouter.POST("", controllers.UserRegister)
			userRouter.PUT("", controllers.UserLogin)
			userRouter.Use(middlewares.VerifyJWT())
			userRouter.GET("", controllers.GetAllUsers)

			// 单个用户路由
			singleUserRouter := userRouter.Group("/:user_id")
			{
				singleUserRouter.POST("/avatar", controllers.UploadAvatar)
				singleUserRouter.GET("/avatar", controllers.GetAvatar)
			}
		}
		forumRouter := api.Group("/forums")
		{
			forumRouter.GET("", controllers.GetAllPublicFroums)
			forumRouter.POST("", middlewares.VerifyJWT(), controllers.CreateForum)
			// 单个论坛路由
			singleForumRouter := forumRouter.Group("/:forum_id")
			{
				singleForumRouter.POST("/cover", middlewares.VerifyJWT(), controllers.UploadCover)
				singleForumRouter.GET("/cover", controllers.GetCover)

				// post 路由
				postRouter := singleForumRouter.Group("/posts")
				postRouter.Use(middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum())
				{
					postRouter.POST("", controllers.CreatePost)
					postRouter.GET("", controllers.GetAllPostsByForumID)

					singlePostRouter := postRouter.Group("/:post_id")
					{
						singlePostRouter.GET("", controllers.GetOnePostDetailByPostID)

						fileRouter := singlePostRouter.Group("/files")
						{
							fileRouter.GET("", controllers.GetFilesByPostID)
						}
					}
				}

				// hole 路由
				holeRouter := singleForumRouter.Group("/holes")
				holeRouter.Use(middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum())
				{
					holeRouter.POST("", controllers.CreateHole)
					holeRouter.GET("", controllers.GetAllHolesByForumID)

					singleHoleRouter := holeRouter.Group("/:hole_id")
					{
						singleHoleRouter.GET("", controllers.GetOneHoleDetailByHoleID)
					}
				}
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
