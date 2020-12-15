package main

import (
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/service-computing-2020/bbs_backend/controllers"
	"github.com/service-computing-2020/bbs_backend/middlewares"

	//"github.com/service-computing-2020/bbs_backend/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/service-computing-2020/bbs_backend/docs" // docs is generated by Swag CLI, you have to import it.
)

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())
	// api文档自动生成
	url := ginSwagger.URL("http://localhost:5000/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
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
				singleUserRouter.GET("/subscribe", controllers.GetOneUserSubscribe)
			}
		}
		forumRouter := api.Group("/forums")
		{
			forumRouter.GET("", middlewares.VerifyJWT(), controllers.GetAllPublicFroums)
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
							fileRouter.GET("/:filename", controllers.GetOneFile)
						}

						// comment 路由
						commentRouter := singlePostRouter.Group("/comments")
						commentRouter.Use(middlewares.VerifyJWT(), middlewares.CanUserWatchTheForum())
						{
							commentRouter.POST("", controllers.CreateComment)
							commentRouter.GET("", controllers.GetAllCommentsByPostID)

							singleCommentRouter := commentRouter.Group("/:comment_id")
							{
								singleCommentRouter.GET("", controllers.GetOneCommentDetailByCommentID)
							}
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

				// role 路由
				roleRouter := singleForumRouter.Group("/role")
				{
					roleRouter.POST("", middlewares.VerifyJWT(), controllers.SubscribeForum)
					roleRouter.DELETE("", middlewares.VerifyJWT(), controllers.UnSubscribeForum)
					roleRouter.GET("/:user_id", controllers.GetRoleInForum)
					roleRouter.PATCH("/:user_id", middlewares.VerifyJWT(), controllers.UpdateRoleInForum)
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
