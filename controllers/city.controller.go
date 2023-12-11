package controllers

import (
	"net/http"
	"strconv"

	"github.com/bryansamperura/ticket-booking/models"
	"github.com/labstack/echo/v4"
)

// FetchAllCities returns a list of all cities
// @Summary Get a list of all cities
// @Description Retrieve a list of all cities
// @Tags City
// @Security Bearer
// @Produce json
// @Success 200 {array} models.City
// @Failure 500 {object} models.HTTPError
// @Router /cities [get]
func FetchAllCities(c echo.Context) error {
	result, err := models.FindAllCity()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

// GetCityById return city by ID
// @Summary Get city by id
// @Description Returns the city with the given id
// @Tags City
// @Security Bearer
// @Param id path int true "City ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} models.City
// @Failure 404 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /city/{id} [get]
func GetCityById(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.FindCityById(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusOK, result)
}

// StoreCity stores city data
// @Summary Create a new city
// @Description Save a new city to the database
// @Tags City
// @Security Bearer
// @Accept json
// @Consumes json
// @Param city body models.CityRequest true "City Name"
// @Success 201 {object} models.City
// @Failure 500 {object} models.HTTPError
// @Router /city [post]
func StoreCity(c echo.Context) error {

	request := new(models.CityRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	result, err := models.StoreCity(request.CityName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

// UpdateCity stores city data
// @Summary Update City
// @Description Updates an existing city in the database
// @Tags City
// @Security Bearer
// @Accept json
// @Consumes json
// @Param id path int true "City ID"
// @Param city body models.CityRequest true "City Name"
// @Success 204 {object} string
// @Failure 500 {object} models.HTTPError
// @Router /city/{id} [put]
func UpdateCity(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	request := new(models.CityRequest)

	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	result, err := models.UpdateCity(id, request.CityName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

// DeleteCity delete city by id
// @Summary Delete City
// @Description Deletes an existing city from the database
// @Tags City
// @Security Bearer
// @Produce json
// @Param id path int true "City ID"
// @Success 204 {object} string
// @Failure 400 {object} models.HTTPError
// @Failure 500 {object} models.HTTPError
// @Router /city/{id} [delete]
func DeleteCity(c echo.Context) error {
	// Get the path parameter "id" as a string
	idStr := c.Param("id")

	// Convert the string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}

	result, err := models.DeleteCity(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if result.Status == 404 {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Not Found"})
	}

	return c.JSON(http.StatusNoContent, result)
}
