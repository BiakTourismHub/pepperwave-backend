package routes

import (
	"github.com/bryansamperura/ticket-booking/controllers"
	_ "github.com/bryansamperura/ticket-booking/docs"
	"github.com/bryansamperura/ticket-booking/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Init() *echo.Echo {
	e := echo.New()

	e.Static("/uploads", "uploads")

	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	e.Use(cors)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Server Is Successfully Running"})
	})

	Authorization := middlewares.AuthMiddleware

	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)
	e.GET("/account-info", controllers.GetAccountInfo, Authorization)

	e.POST("/test", controllers.Test)

	e.GET("/cities", controllers.FetchAllCities, Authorization)
	e.POST("/city", controllers.StoreCity, Authorization)
	e.GET("/city/:id", controllers.GetCityById, Authorization)
	e.PUT("/city/:id", controllers.UpdateCity, Authorization)
	e.DELETE("/city/:id", controllers.DeleteCity, Authorization)

	e.GET("/customers", controllers.FetchAllCustomers, Authorization)
	e.POST("/customer", controllers.StoreCustomer, Authorization)
	e.GET("/customer/:id", controllers.GetCustomerById, Authorization)
	e.PUT("/customer/:id", controllers.UpdateCustomer, Authorization)
	e.DELETE("/customer/:id", controllers.DeleteCustomer, Authorization)

	e.GET("/destination", controllers.FetchAllDestination)
	e.POST("/destination", controllers.StoreDestination)
	e.GET("/destination/:id", controllers.GetDestinationById)
	e.PUT("/destination/:id", controllers.UpdateDestination)
	e.DELETE("/destination/:id", controllers.DeleteDestination)

	e.GET("/booking", controllers.FetchAllBooking, Authorization)
	e.POST("/booking", controllers.StoreBooking, Authorization)
	e.GET("/booking/:customer_id", controllers.GetBookingById, Authorization)

	e.GET("/admin", controllers.FetchAllCustomers)
	e.POST("/admin", controllers.StoreAdmin)
	e.GET("/admin/:id", controllers.GetAdminById)
	e.PUT("/admin/:id", controllers.UpdateAdmin)
	e.DELETE("/admin/:id", controllers.DeleteAdmin)

	return e
}
