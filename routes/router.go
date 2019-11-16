package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spm/api"
	"spm/middleware"
)

func APIRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	router := gin.Default()

	r := router.Group("/api")
	{
		// Ping test
		r.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		r.POST("/signup", api.SignUp)
		r.POST("/signin", api.SignIn)

		auth := r.Group("/", middleware.Auth())
		{
			auth.GET("/signout", api.SignOut)

			auth.GET("/profile", api.GetProfile)
			auth.PUT("/profile", api.EditProfile)

			auth.POST("/appointment", api.MakeAppointment)
			auth.PUT("/appointment", api.EditAppointment)
			auth.DELETE("/appointment", api.DeleteAppointment)
			auth.GET("/appointment", api.GetAppointments)

			auth.GET("/all", api.GetAllAppointments)

			auth.GET("/profile/:id", api.GetAnyProfile)

		}

	}

	// add router here

	return router
}
