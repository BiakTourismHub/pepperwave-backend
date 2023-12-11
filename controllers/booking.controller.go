package controllers

import (
	"net/http"
	"strconv"

	"github.com/bryansamperura/ticket-booking/models"
	"github.com/labstack/echo/v4"
)

// FetchAllBooking returns a list of all booking
// @Summary Get a list of all booking
// @Description Retrieve a list of all booking
// @Tags Booking
// @Security Bearer
// @Produce json
// @Success 200 {array} models.Booking
// @Failure 500 {object} models.HTTPError
// @Router /booking [get]
func FetchAllBooking(c echo.Context) error {
	result, err := models.FindAllBooking()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// GetBookingById return booking by ID
// @Summary Get booking by id
// @Description Returns the booking with the given id
// @Tags Booking
// @Security Bearer
// @Param id path int true "Customer ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Booking
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /booking/{id} [get]
func GetBookingById(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("customer_id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.FindBookingById(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusOK, result)
}

// StoreBooking stores booking data
// @Summary Create a new booking
// @Description Save a new booking to the database
// @Tags Booking
// @Security Bearer
// @Accept json
// @Consumes json
// @Param booking body models.BookingRequest true "Booking Name"
// @Success 201 {object} models.Booking
// @Failure 500 {object} models.HTTPError
// @Router /booking [post]
func StoreBooking(c echo.Context) error {

	request := new(models.BookingRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	result, err := models.StoreBooking(request.CustomerID, request.Qty, request.DestinationID, request.TanggalBooking)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}
