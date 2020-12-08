package main

import (
	"github.com/service-computing-2020/bbs_backend/controllers"
	"github.com/service-computing-2020/bbs_backend/middlewares"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	api := router.Group("/api")
	{
		userRouter := api.Group("/users")
		{
			userRouter.POST("/register", controllers.UserRegister)
			userRouter.POST("/login", controllers.UserLogin)

			userRouter.Use(middlewares.VerifyJWT())
			userRouter.GET("/", controllers.GetAllUsers)
		}

	}
	router.Run(":5000")
}