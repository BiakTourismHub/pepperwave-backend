package main

import (
	"github.com/bryansamperura/ticket-booking/db"
	"github.com/bryansamperura/ticket-booking/routes"
)

// @title API Documentation - Ticket Wisata Booking API
// @version 1.0
// @description This is the documentation for the ticket booking API. It provides information about endpoints and their functionality.
// @contact.name Bryan Samperura

// @contact.url https://www.linkedin.com/in/bryansamperura/
// @license.name MIT License
// @license.url http://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @SecurityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// Use the Authorization header with a Bearer token for authentication.

func main() {
	db.Init()

	e := routes.Init()

	e.Logger.Fatal(e.Start(":3000"))
}
