package routes

import (
	"github.com/labstack/echo/v4"
	"stan.com/stantest/controllers"
)

func SetupRoutes(e *echo.Echo) {
	// episode processing api version 1
	v1 := e.Group("/api/v1")
	{
		// controllers mapping
		users := v1.Group("/episodes")
		{
			users.POST("", controllers.DealwithEpisodes)
		}

		// checking healthy maybe needed by third party
		v1.GET("/health", func(c echo.Context) error {
			return c.JSON(200, map[string]string{
				"status": "ok",
			})
		})
	}
}
